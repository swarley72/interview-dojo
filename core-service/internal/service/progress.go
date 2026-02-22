package service

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/swarley72/interview-dojo/core-service/internal/repository"
)

var (
	ErrProgressNotFound     = errors.New("user progress not found")
	ErrInvalidAnswerQuality = errors.New("invalid answer quality, must be 0, 3, 4 or 5")
	ErrNoQuestionsAvailable = errors.New("no questions available for review")
)

var validAnswerQualities = map[int]bool{
	0: true,
	3: true,
	4: true,
	5: true,
}

type UserProgressService interface {
	GetProgress(ctx context.Context, userID string, questionID string) (repository.UserProgress, error)
	RecordAnswer(ctx context.Context, userID string, questionID string, answerQuality int) (repository.UserProgress, error)
	GetNextQuestion(ctx context.Context, userID string, filters repository.NextQuestionFilters) (repository.Question, error)
}

type userProgressService struct {
	questionRepo     repository.QuestionRepository
	userProgressRepo repository.UserProgressRepository
}

func (s *userProgressService) GetProgress(ctx context.Context, userID string, questionID string) (repository.UserProgress, error) {
	userProgress, err := s.userProgressRepo.GetProgress(ctx, userID, questionID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repository.UserProgress{}, ErrProgressNotFound
		}
		return repository.UserProgress{}, err
	}
	return userProgress, nil
}

func (s *userProgressService) RecordAnswer(ctx context.Context, userID string, questionID string, answerQuality int) (repository.UserProgress, error) {
	if !validAnswerQualities[answerQuality] {
		return repository.UserProgress{}, ErrInvalidAnswerQuality
	}

	progress, err := s.userProgressRepo.GetProgress(ctx, userID, questionID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return repository.UserProgress{}, err
	}

	result := calculateSM2(progress.Repetitions, progress.EaseFactor, progress.IntervalDays, answerQuality)

	return s.userProgressRepo.UpsertProgress(ctx, repository.UpsertProgressParams{
		UserID:       userID,
		QuestionID:   questionID,
		Repetitions:  result.Repetitions,
		EaseFactor:   result.EaseFactor,
		IntervalDays: result.IntervalDays,
		NextReviewAt: time.Now().AddDate(0, 0, int(result.IntervalDays)),
	})
}

func (s *userProgressService) GetNextQuestion(ctx context.Context, userID string, filters repository.NextQuestionFilters) (repository.Question, error) {
	questionID, err := s.userProgressRepo.GetDueQuestionID(ctx, userID, filters)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return repository.Question{}, err
	}

	if errors.Is(err, pgx.ErrNoRows) {
		questionID, err = s.questionRepo.GetNewQuestionID(ctx, userID, filters)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return repository.Question{}, ErrNoQuestionsAvailable
			}
			return repository.Question{}, err
		}
	}

	return s.questionRepo.GetQuestionByID(ctx, questionID)
}

func NewUserProgressService(
	questionRepo repository.QuestionRepository,
	userProgressRepo repository.UserProgressRepository,
) UserProgressService {
	return &userProgressService{questionRepo: questionRepo, userProgressRepo: userProgressRepo}
}
