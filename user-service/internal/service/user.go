package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/swarley72/interview-dojo/user-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(ctx context.Context, login, password string) (AuthResponse, error)
	Login(ctx context.Context, login, password string) (AuthResponse, error)
	GetUser(ctx context.Context, id string) (repository.User, error)
	ValidateToken(accessToken string) (Claims, error)
	RefreshToken(ctx context.Context, refreshToken string) (TokenPair, error)
}

type Claims struct {
	UserID      string
	IsSuperUser bool
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type AuthResponse struct {
	User  repository.User
	Token TokenPair
}

type userService struct {
	accessTokenExp  time.Duration
	refreshTokenExp time.Duration
	jwtSecret       []byte
	repo            repository.UserRepository
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
)

const (
	claimUserID      = "user_id"
	claimIsSuperUser = "is_super_user"
	claimExp         = "exp"
)

func (us *userService) generateTokens(userID string, isSuperUser bool) (TokenPair, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		claimUserID:      userID,
		claimIsSuperUser: isSuperUser,
		claimExp:         time.Now().UTC().Add(us.accessTokenExp).Unix(),
	})
	signedAccessToken, err := accessToken.SignedString(us.jwtSecret)

	if err != nil {
		return TokenPair{}, err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		claimUserID: userID,
		claimExp:    time.Now().UTC().Add(us.refreshTokenExp).Unix(),
	})
	signedRefreshToken, err := refreshToken.SignedString(us.jwtSecret)

	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{AccessToken: signedAccessToken, RefreshToken: signedRefreshToken}, nil
}

func (us *userService) generatePassword(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(passwordHash), nil
}

func (us *userService) comparePasswords(passwordHash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))

	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return ErrInvalidCredentials
	}

	return err
}

func (us *userService) Register(ctx context.Context, login, password string) (AuthResponse, error) {
	passwordHash, err := us.generatePassword(password)

	if err != nil {
		return AuthResponse{}, err
	}

	user, err := us.repo.CreateUser(ctx, login, string(passwordHash))

	if err != nil {
		var pgErr *pgconn.PgError
		// https://www.postgresql.org/docs/10/errcodes-appendix.html
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return AuthResponse{}, ErrUserAlreadyExists
		}
		return AuthResponse{}, err
	}

	tokens, err := us.generateTokens(user.ID, user.IsSuperUser)

	if err != nil {
		return AuthResponse{}, err
	}

	return AuthResponse{User: user, Token: tokens}, nil
}

func (us *userService) Login(ctx context.Context, login, password string) (AuthResponse, error) {
	user, err := us.repo.GetUserByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return AuthResponse{}, ErrInvalidCredentials
		}
		return AuthResponse{}, err
	}

	err = us.comparePasswords(user.PasswordHash, password)
	if err != nil {
		return AuthResponse{}, err
	}

	tokens, err := us.generateTokens(user.ID, user.IsSuperUser)
	if err != nil {
		return AuthResponse{}, err
	}

	return AuthResponse{User: user, Token: tokens}, nil
}

func (us *userService) GetUser(ctx context.Context, id string) (repository.User, error) {
	user, err := us.repo.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repository.User{}, ErrUserNotFound
		}

		return repository.User{}, err

	}
	return user, nil
}

func (us *userService) parseToken(token string) (jwt.MapClaims, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return us.jwtSecret, nil
	})

	if err != nil {
		return jwt.MapClaims{}, err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return jwt.MapClaims{}, ErrInvalidToken
	}

	return claims, nil
}

func (us *userService) ValidateToken(token string) (Claims, error) {
	claims, err := us.parseToken(token)
	if err != nil {
		return Claims{}, err
	}

	userID, ok := claims[claimUserID].(string)
	if !ok {
		return Claims{}, ErrInvalidToken
	}

	isSuperUser, ok := claims[claimIsSuperUser].(bool)
	if !ok {
		return Claims{}, ErrInvalidToken
	}

	return Claims{UserID: userID, IsSuperUser: isSuperUser}, nil
}

func (us *userService) RefreshToken(ctx context.Context, token string) (TokenPair, error) {
	claims, err := us.parseToken(token)
	if err != nil {
		return TokenPair{}, err
	}

	userID, ok := claims[claimUserID].(string)
	if !ok {
		return TokenPair{}, ErrInvalidToken
	}

	user, err := us.GetUser(ctx, userID)
	if err != nil {
		return TokenPair{}, err
	}

	return us.generateTokens(user.ID, user.IsSuperUser)
}

func NewUserService(
	repo repository.UserRepository,
	jwtSecret []byte,
	accessTokenExp time.Duration,
	refreshTokenExp time.Duration,
) UserService {
	return &userService{
		repo:            repo,
		jwtSecret:       jwtSecret,
		accessTokenExp:  accessTokenExp,
		refreshTokenExp: refreshTokenExp}
}
