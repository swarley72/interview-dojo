package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresUserRepository struct {
	pool *pgxpool.Pool
}

func (r *postgresUserRepository) CreateUser(ctx context.Context, login, passwordHash string) (User, error) {
	var user User
	query := `INSERT INTO users (login, password_hash)
			  VALUES ($1, $2)
			  RETURNING id, login, password_hash, created_at, updated_at, is_super_user;`
	err := r.pool.QueryRow(ctx, query, login, passwordHash).Scan(
		&user.ID,
		&user.Login,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.IsSuperUser,
	)
	return user, err
}

func (r *postgresUserRepository) GetUserByLogin(ctx context.Context, login string) (User, error) {
	var user User
	query := `SELECT id, login, password_hash, created_at, updated_at, is_super_user FROM users WHERE login = $1`
	err := r.pool.QueryRow(ctx, query, login).Scan(
		&user.ID,
		&user.Login,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.IsSuperUser,
	)
	return user, err
}

func (r *postgresUserRepository) GetUserByID(ctx context.Context, id string) (User, error) {
	var user User
	query := `SELECT id, login, password_hash, created_at, updated_at, is_super_user FROM users WHERE id = $1`
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Login,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.IsSuperUser,
	)
	return user, err
}

func NewPostgresUserRepository(pool *pgxpool.Pool) UserRepository {
	return &postgresUserRepository{pool}
}
