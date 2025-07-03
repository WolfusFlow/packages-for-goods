package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) GetPackSizes(ctx context.Context) ([]int, error) {
	rows, err := r.pool.Query(ctx, `SELECT size FROM pack_sizes ORDER BY size ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sizes := []int{}
	for rows.Next() {
		var size int
		if err := rows.Scan(&size); err != nil {
			return nil, err
		}
		sizes = append(sizes, size)
	}
	return sizes, nil
}

func (r *Repository) InsertPackSize(ctx context.Context, size int) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO pack_sizes (size) VALUES ($1) ON CONFLICT DO NOTHING`, size)
	return err
}

func (r *Repository) DeletePackSize(ctx context.Context, size int) error {
	cmd, err := r.pool.Exec(ctx, `DELETE FROM pack_sizes WHERE size = $1`, size)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("pack size not found")
	}
	return nil
}
