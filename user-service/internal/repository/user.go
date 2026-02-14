package repository

import (
	"context"
	"time"
)

type UserRepository interface {
	CreateUser(ctx context.Context, login, passwordHash string) (User, error)
	GetUserByLogin(ctx context.Context, login string) (User, error)
	GetUserByID(ctx context.Context, id string) (User, error)
}

type User struct {
	ID           string
	Login        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	IsSuperUser  bool
}
