package repository

import (
	"context"
	"time"
)

type QuestionRepository interface {
	CreateQuestion(ctx context.Context, question CreateQuestionParams) (Question, error)
	GetQuestionByID(ctx context.Context, id string) (Question, error)
	UpdateQuestion(ctx context.Context, id string, question UpdateQuestionParams) (Question, error)
	DeleteQuestion(ctx context.Context, id string) error
	ListQuestions(ctx context.Context, filters ListQuestionsFilters) ([]Question, error)
	GetNewQuestionID(ctx context.Context, userID string) (string, error)
}

type Question struct {
	ID         string
	Type       string
	Title      string
	Difficulty string
	TagIDs     []int
	ContentMD  *string
	AnswerMD   *string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type CreateQuestionParams struct {
	TagIDs     []int
	Type       string
	Title      string
	Difficulty string
	ContentMD  *string
	AnswerMD   *string
}

type UpdateQuestionParams struct {
	Type       *string
	Title      *string
	Difficulty *string
	ContentMD  *string
	AnswerMD   *string
}

type ListQuestionsFilters struct {
	Limit      int
	Offset     int
	Type       *string
	Difficulty *string
}
