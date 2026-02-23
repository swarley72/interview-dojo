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
	ListQuestions(ctx context.Context, filters ListQuestionsFilters) (ListQuestionsResult, error)
	GetNewQuestionID(ctx context.Context, userID string, filters NextQuestionFilters) (string, error)
}

type Question struct {
	ID             string
	Type           string
	Title          string
	Difficulty     string
	TagIDs         []int32
	ContentMD      *string
	AnswerMD       *string
	ExcalidrawJSON *string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Verified       bool
}

type CreateQuestionParams struct {
	TagIDs         []int32
	Type           string
	Title          string
	Difficulty     string
	ContentMD      *string
	AnswerMD       *string
	ExcalidrawJSON *string
	Verified       bool
}

type UpdateQuestionParams struct {
	Type           *string
	Title          *string
	Difficulty     *string
	ContentMD      *string
	AnswerMD       *string
	ExcalidrawJSON *string
	TagIDs         []int32
	Verified       *bool
}

type ListQuestionsFilters struct {
	Limit      int
	Offset     int
	Type       *string
	Difficulty *string
	Query      *string
	TagIDs     []int32
	Verified   *bool
}

type ListQuestionsResult struct {
	Questions  []Question
	TotalCount int32
}

type NextQuestionFilters struct {
	TagIDs       []int32
	QuestionType *string
}
