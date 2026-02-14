package server

import (
	"context"
	"errors"

	userpb "github.com/swarley72/interview-dojo/proto/user"
	"github.com/swarley72/interview-dojo/user-service/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	userpb.UnimplementedUserServiceServer
	service service.UserService
}

func (g *GRPCServer) Register(
	ctx context.Context,
	req *userpb.RegisterRequest,
) (*userpb.AuthResponse, error) {
	res, err := g.service.Register(ctx, req.Login, req.Password)
	if err != nil {
		return nil, mapError(err)
	}

	return &userpb.AuthResponse{
		User: &userpb.User{
			Id:          res.User.ID,
			Login:       res.User.Login,
			IsSuperUser: res.User.IsSuperUser,
		},
		Token: &userpb.Token{
			AccessToken:  res.Token.AccessToken,
			RefreshToken: res.Token.RefreshToken,
		},
	}, nil
}

func (g *GRPCServer) Login(
	ctx context.Context,
	req *userpb.LoginRequest,
) (*userpb.AuthResponse, error) {
	res, err := g.service.Login(ctx, req.Login, req.Password)
	if err != nil {
		return nil, mapError(err)
	}

	return &userpb.AuthResponse{
		User: &userpb.User{
			Id:          res.User.ID,
			Login:       res.User.Login,
			IsSuperUser: res.User.IsSuperUser,
		},
		Token: &userpb.Token{
			AccessToken:  res.Token.AccessToken,
			RefreshToken: res.Token.RefreshToken,
		},
	}, nil
}

func (g *GRPCServer) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.User, error) {
	user, err := g.service.GetUser(ctx, req.UserId)
	if err != nil {
		return nil, mapError(err)
	}

	return &userpb.User{
		Id:          user.ID,
		Login:       user.Login,
		IsSuperUser: user.IsSuperUser,
	}, nil
}

func (g *GRPCServer) ValidateToken(_ context.Context, req *userpb.ValidateTokenRequest) (*userpb.ValidateTokenResponse, error) {
	res, err := g.service.ValidateToken(req.AccessToken)
	if err != nil {
		return nil, mapError(err)
	}

	return &userpb.ValidateTokenResponse{
		UserId:      res.UserID,
		IsSuperUser: res.IsSuperUser,
	}, nil
}

func (g *GRPCServer) RefreshToken(ctx context.Context, req *userpb.RefreshTokenRequest) (*userpb.Token, error) {
	token, err := g.service.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, mapError(err)
	}
	return &userpb.Token{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, nil
}

func NewGRPCServer(userService service.UserService) *GRPCServer {
	return &GRPCServer{service: userService}
}

func mapError(err error) error {
	switch {
	case errors.Is(err, service.ErrUserAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, service.ErrInvalidCredentials):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, service.ErrUserNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, service.ErrInvalidToken):
		return status.Error(codes.Unauthenticated, err.Error())
	default:
		return status.Error(codes.Internal, "internal error")
	}
}
