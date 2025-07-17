package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Conn abstracts DB lifecycle and access
type Conn interface {
	Close() error
	Pool() *pgxpool.Pool
}

// pgConn implements Conn
type pgConn struct {
	pool *pgxpool.Pool
}

func (p *pgConn) Close() error {
	p.pool.Close()
	return nil
}

func (p *pgConn) Pool() *pgxpool.Pool {
	return p.pool
}

func Connect(url string) (Conn, error) {
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DB URL: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	return &pgConn{pool: pool}, nil
}
