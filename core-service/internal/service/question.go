package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/swarley72/interview-dojo/core-service/internal/repository"
)

var (
	ErrQuestionNotFound = errors.New("question not found")
)

type QuestionService interface {
	CreateQuestion(ctx context.Context, params repository.CreateQuestionParams) (repository.Question, error)
	GetQuestionByID(ctx context.Context, id string) (repository.Question, error)
	UpdateQuestion(ctx context.Context, id string, params repository.UpdateQuestionParams) (repository.Question, error)
	DeleteQuestion(ctx context.Context, id string) error
	ListQuestions(ctx context.Context, filters repository.ListQuestionsFilters) (repository.ListQuestionsResult, error)
}

type questionService struct {
	questionRepo repository.QuestionRepository
}

func (q *questionService) CreateQuestion(ctx context.Context, params repository.CreateQuestionParams) (repository.Question, error) {
	return q.questionRepo.CreateQuestion(ctx, params)
}

func (q *questionService) UpdateQuestion(
	ctx context.Context,
	id string,
	params repository.UpdateQuestionParams,
) (repository.Question, error) {
	question, err := q.questionRepo.UpdateQuestion(ctx, id, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repository.Question{}, ErrQuestionNotFound
		}
		return repository.Question{}, err
	}

	return question, nil
}

func (q *questionService) GetQuestionByID(ctx context.Context, id string) (repository.Question, error) {
	question, err := q.questionRepo.GetQuestionByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repository.Question{}, ErrQuestionNotFound
		}
		return repository.Question{}, err
	}

	return question, nil
}

func (q *questionService) ListQuestions(ctx context.Context, filters repository.ListQuestionsFilters) (repository.ListQuestionsResult, error) {
	return q.questionRepo.ListQuestions(ctx, filters)
}

func (q *questionService) DeleteQuestion(ctx context.Context, id string) error {
	return q.questionRepo.DeleteQuestion(ctx, id)
}

func NewQuestionService(questionRepo repository.QuestionRepository) QuestionService {
	return &questionService{questionRepo}
}
