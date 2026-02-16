package repository

import (
	"context"
)

type TagRepository interface {
	CreateTag(ctx context.Context, name string) (Tag, error)
	ListTags(ctx context.Context) ([]Tag, error)
	GetTagByID(ctx context.Context, id int) (Tag, error)
	DeleteTag(ctx context.Context, id int) error
}

type Tag struct {
	ID   int
	Name string
}
