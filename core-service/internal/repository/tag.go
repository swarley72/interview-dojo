package repository

import (
	"context"
)

type TagRepository interface {
	CreateTag(ctx context.Context, name string) (Tag, error)
	ListTags(ctx context.Context) ([]Tag, error)
	GetTagByID(ctx context.Context, id int32) (Tag, error)
	DeleteTag(ctx context.Context, id int32) error
}

type Tag struct {
	ID   int32
	Name string
}
