package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/swarley72/interview-dojo/core-service/internal/repository"
)

type TagService interface {
	CreateTag(ctx context.Context, name string) (repository.Tag, error)
	DeleteTag(ctx context.Context, id int32) error
	ListTags(ctx context.Context) ([]repository.Tag, error)
}

type tagService struct {
	tagRepo repository.TagRepository
}

var ErrTagAlreadyExists = errors.New("tag already exists")

func (t *tagService) CreateTag(ctx context.Context, name string) (repository.Tag, error) {
	name = strings.TrimSpace(name)

	tag, err := t.tagRepo.CreateTag(ctx, name)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return repository.Tag{}, fmt.Errorf("tag %q: %w", name, ErrTagAlreadyExists)
		}
		return repository.Tag{}, err
	}

	return tag, nil
}

func (t *tagService) DeleteTag(ctx context.Context, id int32) error {
	return t.tagRepo.DeleteTag(ctx, id)
}

func (t *tagService) ListTags(ctx context.Context) ([]repository.Tag, error) {
	return t.tagRepo.ListTags(ctx)
}

func NewTagService(tagRepo repository.TagRepository) TagService {
	return &tagService{tagRepo}
}
