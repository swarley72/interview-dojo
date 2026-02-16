package service

import (
	"context"
	"strings"

	"github.com/swarley72/interview-dojo/core-service/internal/repository"
)

type TagService interface {
	CreateTag(ctx context.Context, name string) (repository.Tag, error)
	DeleteTag(ctx context.Context, id int) error
	ListTags(ctx context.Context) ([]repository.Tag, error)
}

type tagService struct {
	tagRepo repository.TagRepository
}

func (t *tagService) CreateTag(ctx context.Context, name string) (repository.Tag, error) {
	name = strings.TrimSpace(name)
	return t.tagRepo.CreateTag(ctx, name)
}

func (t *tagService) DeleteTag(ctx context.Context, id int) error {
	return t.tagRepo.DeleteTag(ctx, id)
}

func (t *tagService) ListTags(ctx context.Context) ([]repository.Tag, error) {
	return t.tagRepo.ListTags(ctx)
}

func NewTagService(tagRepo repository.TagRepository) TagService {
	return &tagService{tagRepo}
}
