package repository

import (
	"context"
	"time"
)

type UserProgressRepository interface {
	GetProgress(ctx context.Context, userID string, questionID string) (UserProgress, error)
	UpsertProgress(ctx context.Context, params UpsertProgressParams) (UserProgress, error)
	GetDueQuestionID(ctx context.Context, userID string, filters NextQuestionFilters) (string, error)
}

type UserProgress struct {
	UserID         string
	QuestionID     string
	ID             int32
	Repetitions    int32
	IntervalDays   int32
	EaseFactor     float64
	LastReviewedAt time.Time
	NextReviewAt   time.Time
}

type UpsertProgressParams struct {
	UserID       string
	QuestionID   string
	Repetitions  int32
	IntervalDays int32
	EaseFactor   float64
	NextReviewAt time.Time
}
