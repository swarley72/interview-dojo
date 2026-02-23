package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresTagRepository struct {
	pool *pgxpool.Pool
}

func (t *postgresTagRepository) CreateTag(ctx context.Context, name string) (Tag, error) {
	var tag Tag
	query := `INSERT INTO tags (name) VALUES ($1) RETURNING id, name`
	err := t.pool.QueryRow(ctx, query, name).Scan(
		&tag.ID,
		&tag.Name,
	)
	return tag, err
}

func (t *postgresTagRepository) GetTagByID(ctx context.Context, id int32) (Tag, error) {
	var tag Tag
	query := `SELECT id, name FROM tags WHERE id = $1`
	err := t.pool.QueryRow(ctx, query, id).Scan(
		&tag.ID,
		&tag.Name,
	)
	return tag, err
}

func (t *postgresTagRepository) ListTags(ctx context.Context) ([]Tag, error) {
	query := `SELECT id, name FROM tags ORDER BY name ASC`
	rows, err := t.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []Tag
	for rows.Next() {
		var tag Tag
		err := rows.Scan(&tag.ID, &tag.Name)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, rows.Err()
}

func (t *postgresTagRepository) DeleteTag(ctx context.Context, id int32) error {
	query := `DELETE FROM tags WHERE id = $1`
	_, err := t.pool.Exec(ctx, query, id)
	return err
}

func NewPostgresTagRepository(pool *pgxpool.Pool) TagRepository {
	return &postgresTagRepository{pool}
}
