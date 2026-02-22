package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresUserProgressRepository struct {
	pool *pgxpool.Pool
}

func (up *postgresUserProgressRepository) GetProgress(ctx context.Context, userID string, questionID string) (UserProgress, error) {
	query := `
	SELECT id, user_id, question_id, repetitions, ease_factor, interval_days, last_reviewed_at, next_review_at
	FROM user_progress
	WHERE user_id = $1 AND question_id = $2
	`
	var progress UserProgress
	err := up.pool.QueryRow(ctx, query, userID, questionID).Scan(
		&progress.ID,
		&progress.UserID,
		&progress.QuestionID,
		&progress.Repetitions,
		&progress.EaseFactor,
		&progress.IntervalDays,
		&progress.LastReviewedAt,
		&progress.NextReviewAt,
	)
	if err != nil {
		return UserProgress{}, err
	}

	return progress, nil
}

func (up *postgresUserProgressRepository) UpsertProgress(ctx context.Context, params UpsertProgressParams) (UserProgress, error) {
	query := `
	INSERT INTO user_progress (user_id, question_id, repetitions, ease_factor, interval_days, last_reviewed_at, next_review_at)
	VALUES ($1, $2, $3, $4, $5, now(), $6)
	ON CONFLICT (user_id, question_id) DO UPDATE SET
		repetitions = EXCLUDED.repetitions,
		ease_factor = EXCLUDED.ease_factor,
		interval_days = EXCLUDED.interval_days,
		last_reviewed_at = now(),
		next_review_at = EXCLUDED.next_review_at
	RETURNING id, user_id, question_id, repetitions, ease_factor, interval_days, last_reviewed_at, next_review_at 
	`
	var progress UserProgress
	err := up.pool.QueryRow(
		ctx,
		query,
		params.UserID,
		params.QuestionID,
		params.Repetitions,
		params.EaseFactor,
		params.IntervalDays,
		params.NextReviewAt,
	).Scan(
		&progress.ID,
		&progress.UserID,
		&progress.QuestionID,
		&progress.Repetitions,
		&progress.EaseFactor,
		&progress.IntervalDays,
		&progress.LastReviewedAt,
		&progress.NextReviewAt,
	)
	if err != nil {
		return UserProgress{}, err
	}

	return progress, nil
}

func (up *postgresUserProgressRepository) GetDueQuestionID(ctx context.Context, userID string, filters NextQuestionFilters) (string, error) {
	query := `
	SELECT question_id FROM user_progress up
	JOIN questions q ON q.id = up.question_id
	WHERE up.user_id = $1 AND up.next_review_at <= now()
	`
	args := []any{userID}

	if filters.QuestionType != nil {
		args = append(args, *filters.QuestionType)
		query += fmt.Sprintf(" AND q.type = $%d", len(args))
	}

	if len(filters.TagIDs) > 0 {
		args = append(args, filters.TagIDs)
		query += fmt.Sprintf(" AND q.id IN (SELECT DISTINCT question_id FROM question_tags WHERE tag_id = ANY($%d))", len(args))
	}
	query += " ORDER BY next_review_at ASC LIMIT 1"

	var questionID string
	err := up.pool.QueryRow(ctx, query, args...).Scan(&questionID)
	if err != nil {
		return "", err
	}

	return questionID, nil
}

func (up *postgresUserProgressRepository) ResetProgress(ctx context.Context, userID string) error {
	_, err := up.pool.Exec(ctx, "DELETE FROM user_progress WHERE user_id = $1", userID)
	if err != nil {
		return err
	}

	return nil
}

func NewPostgresUserProgressRepository(pool *pgxpool.Pool) UserProgressRepository {
	return &postgresUserProgressRepository{pool}
}
