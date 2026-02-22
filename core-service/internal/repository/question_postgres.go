package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresQuestionRepository struct {
	pool *pgxpool.Pool
}

func (q *postgresQuestionRepository) CreateQuestion(ctx context.Context, params CreateQuestionParams) (Question, error) {
	tx, err := q.pool.Begin(ctx)
	if err != nil {
		return Question{}, err
	}
	defer tx.Rollback(ctx)

	questionQuery := `
		INSERT INTO questions (type, title, content_md, answer_md, difficulty, excalidraw_json)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, type, title, content_md, answer_md, difficulty, excalidraw_json, created_at, updated_at
	`
	questionTagsQuery := `INSERT INTO question_tags (question_id, tag_id) VALUES ($1, $2)`

	var question Question

	err = tx.QueryRow(
		ctx,
		questionQuery,
		params.Type,
		params.Title,
		params.ContentMD,
		params.AnswerMD,
		params.Difficulty,
		params.ExcalidrawJSON,
	).Scan(
		&question.ID,
		&question.Type,
		&question.Title,
		&question.ContentMD,
		&question.AnswerMD,
		&question.Difficulty,
		&question.ExcalidrawJSON,
		&question.CreatedAt,
		&question.UpdatedAt,
	)
	if err != nil {
		return Question{}, err
	}

	if len(params.TagIDs) > 0 {
		batch := &pgx.Batch{}
		for _, tagID := range params.TagIDs {
			batch.Queue(questionTagsQuery, question.ID, tagID)
		}

		br := tx.SendBatch(ctx, batch)

		for range len(params.TagIDs) {
			_, err := br.Exec()
			if err != nil {
				br.Close()
				return Question{}, err
			}
		}
		br.Close()
	}

	question.TagIDs = params.TagIDs

	err = tx.Commit(ctx)
	if err != nil {
		return Question{}, err
	}

	return question, err
}

func (q *postgresQuestionRepository) GetQuestionByID(ctx context.Context, id string) (Question, error) {
	var question Question
	questionQuery := `
		SELECT id, type, title, content_md, answer_md, difficulty, excalidraw_json, created_at, updated_at
		FROM questions
		WHERE id = $1
	`
	err := q.pool.QueryRow(ctx, questionQuery, id).Scan(
		&question.ID,
		&question.Type,
		&question.Title,
		&question.ContentMD,
		&question.AnswerMD,
		&question.Difficulty,
		&question.ExcalidrawJSON,
		&question.CreatedAt,
		&question.UpdatedAt,
	)
	if err != nil {
		return Question{}, err
	}

	tags, err := q.loadTagsForQuestions(ctx, []string{question.ID})
	if err != nil {
		return Question{}, err
	}

	question.TagIDs = tags[question.ID]

	return question, nil
}

func (q *postgresQuestionRepository) UpdateQuestion(ctx context.Context, id string, params UpdateQuestionParams) (Question, error) {
	tx, err := q.pool.Begin(ctx)
	if err != nil {
		return Question{}, err
	}
	defer tx.Rollback(ctx)

	var setClauses []string
	var args []any
	argIndex := 1

	if params.Type != nil {
		setClauses = append(setClauses, fmt.Sprintf("type = $%d", argIndex))
		args = append(args, *params.Type)
		argIndex++
	}

	if params.Title != nil {
		setClauses = append(setClauses, fmt.Sprintf("title = $%d", argIndex))
		args = append(args, *params.Title)
		argIndex++
	}

	if params.Difficulty != nil {
		setClauses = append(setClauses, fmt.Sprintf("difficulty = $%d", argIndex))
		args = append(args, *params.Difficulty)
		argIndex++
	}

	if params.ContentMD != nil {
		setClauses = append(setClauses, fmt.Sprintf("content_md = $%d", argIndex))
		args = append(args, *params.ContentMD)
		argIndex++
	}

	if params.AnswerMD != nil {
		setClauses = append(setClauses, fmt.Sprintf("answer_md = $%d", argIndex))
		args = append(args, *params.AnswerMD)
		argIndex++
	}

	if params.ExcalidrawJSON != nil {
		setClauses = append(setClauses, fmt.Sprintf("excalidraw_json = $%d", argIndex))
		args = append(args, *params.ExcalidrawJSON)
		argIndex++
	}

	if len(setClauses) == 0 && params.TagIDs == nil {
		return q.GetQuestionByID(ctx, id)
	}

	setClauses = append(setClauses, "updated_at = now()")

	query := fmt.Sprintf(
		"UPDATE questions SET %s WHERE id = $%d RETURNING id, type, title, content_md, answer_md, difficulty, excalidraw_json, created_at, updated_at",
		strings.Join(setClauses, ", "),
		argIndex,
	)

	args = append(args, id)
	var question Question
	err = tx.QueryRow(ctx, query, args...).Scan(
		&question.ID,
		&question.Type,
		&question.Title,
		&question.ContentMD,
		&question.AnswerMD,
		&question.Difficulty,
		&question.ExcalidrawJSON,
		&question.CreatedAt,
		&question.UpdatedAt,
	)
	if err != nil {
		return Question{}, err
	}

	if params.TagIDs != nil {
		_, err = tx.Exec(ctx, "DELETE FROM question_tags WHERE question_id = $1", id)
		if err != nil {
			return Question{}, err
		}

		if len(params.TagIDs) > 0 {
			placeholders := make([]string, 0, len(params.TagIDs))
			args := make([]any, 0, len(params.TagIDs)+1)
			args = append(args, id)
			for i, tagID := range params.TagIDs {
				placeholders = append(placeholders, fmt.Sprintf("($1, $%d)", i+2))
				args = append(args, tagID)
			}

			insertTagsQuery := fmt.Sprintf(
				"INSERT INTO question_tags (question_id, tag_id) VALUES %s",
				strings.Join(placeholders, ", "),
			)
			_, err = tx.Exec(ctx, insertTagsQuery, args...)
			if err != nil {
				return Question{}, err
			}
			question.TagIDs = params.TagIDs
		}

	} else {
		tagsQuery := "SELECT tag_id FROM question_tags WHERE question_id = $1 ORDER BY tag_id"

		rows, err := tx.Query(ctx, tagsQuery, id)
		if err != nil {
			return Question{}, err
		}
		defer rows.Close()

		for rows.Next() {
			var tagID int32
			err := rows.Scan(&tagID)
			if err != nil {
				return Question{}, err
			}
			question.TagIDs = append(question.TagIDs, tagID)
		}

	}

	err = tx.Commit(ctx)
	if err != nil {
		return Question{}, err
	}

	return question, nil
}

func (q *postgresQuestionRepository) DeleteQuestion(ctx context.Context, id string) error {
	query := "DELETE FROM questions WHERE id = $1"
	_, err := q.pool.Exec(ctx, query, id)
	return err
}

func (q *postgresQuestionRepository) ListQuestions(ctx context.Context, filters ListQuestionsFilters) (ListQuestionsResult, error) {
	query := "SELECT id, type, title, difficulty, created_at, updated_at FROM questions"

	var conditions []string
	var filterArgs []any
	argIndex := 1
	whereClause := ""

	if filters.Type != nil {
		conditions = append(conditions, fmt.Sprintf("type = $%d", argIndex))
		filterArgs = append(filterArgs, *filters.Type)
		argIndex++
	}

	if filters.Difficulty != nil {
		conditions = append(conditions, fmt.Sprintf("difficulty = $%d", argIndex))
		filterArgs = append(filterArgs, *filters.Difficulty)
		argIndex++
	}

	if len(filters.TagIDs) > 0 {
		conditions = append(conditions, fmt.Sprintf("id IN (SELECT question_id FROM question_tags WHERE tag_id = ANY($%d))",
			argIndex))
		filterArgs = append(filterArgs, filters.TagIDs)
		argIndex++
	}

	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	query += whereClause
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args := make([]any, len(filterArgs))
	copy(args, filterArgs)
	args = append(args, filters.Limit, filters.Offset)

	questionRows, err := q.pool.Query(ctx, query, args...)
	if err != nil {
		return ListQuestionsResult{}, err
	}
	defer questionRows.Close()

	var questions []Question
	for questionRows.Next() {
		var question Question
		err := questionRows.Scan(
			&question.ID,
			&question.Type,
			&question.Title,
			&question.Difficulty,
			&question.CreatedAt,
			&question.UpdatedAt,
		)
		if err != nil {
			return ListQuestionsResult{}, err
		}

		questions = append(questions, question)
	}

	if err := questionRows.Err(); err != nil {
		return ListQuestionsResult{}, err
	}
	var totalCount int32
	err = q.pool.QueryRow(ctx, "SELECT COUNT(*) FROM questions"+whereClause, filterArgs...).Scan(&totalCount)
	if err != nil {
		return ListQuestionsResult{}, err
	}

	ids := make([]string, len(questions))
	for i, q := range questions {
		ids[i] = q.ID
	}
	tags, err := q.loadTagsForQuestions(ctx, ids)
	if err != nil {
		return ListQuestionsResult{}, err
	}

	for i := range questions {
		questions[i].TagIDs = tags[questions[i].ID]
	}

	return ListQuestionsResult{Questions: questions, TotalCount: totalCount}, nil
}

func (q *postgresQuestionRepository) GetNewQuestionID(ctx context.Context, userID string, filters NextQuestionFilters) (string, error) {
	query := `
	SELECT id FROM questions q
	WHERE q.id NOT IN (SELECT question_id FROM user_progress WHERE user_id = $1)
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

	query += " LIMIT 1"

	var questionID string
	err := q.pool.QueryRow(ctx, query, args...).Scan(&questionID)
	if err != nil {
		return "", err
	}

	return questionID, nil
}

func (r *postgresQuestionRepository) loadTagsForQuestions(ctx context.Context, questionIDs []string) (map[string][]int32, error) {
	query := "SELECT question_id, tag_id FROM question_tags WHERE question_id = ANY($1)"
	rows, err := r.pool.Query(ctx, query, questionIDs)
	if err != nil {
		return map[string][]int32{}, err
	}
	defer rows.Close()

	res := make(map[string][]int32, len(questionIDs))
	for rows.Next() {
		var questionID string
		var tagID int32

		err := rows.Scan(
			&questionID,
			&tagID,
		)
		if err != nil {
			return map[string][]int32{}, err
		}
		res[questionID] = append(res[questionID], tagID)
	}

	return res, rows.Err()
}

func NewPostgresQuestionRepository(pool *pgxpool.Pool) QuestionRepository {
	return &postgresQuestionRepository{pool}
}
