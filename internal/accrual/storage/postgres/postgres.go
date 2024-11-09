package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/errors"
	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/model"
)

type Config struct {
	DatabaseDSN    string
	ConnectTimeout time.Duration
}

type Storager struct {
	db *pgx.Conn
}

func NewPostgresStorage(cfg *Config) (*Storager, error) {
	connCtx, cancel := context.WithTimeoutCause(context.Background(), cfg.ConnectTimeout, errors.ErrConnectTimeout)
	defer cancel()
	conn, err := pgx.Connect(connCtx, cfg.DatabaseDSN)
	if err != nil {
		return nil, fmt.Errorf("postgres connect: %w", err)
	}

	return &Storager{
		db: conn,
	}, nil
}

func (s *Storager) CreateOrder(ctx context.Context, o model.Order) error {
	query := "INSERT INTO orders (number, status, accrual) VALUES ($1, $2, $3)"
	_, err := s.db.Exec(ctx, query, o.Number, o.Status, o.Accrual)
	if err != nil {
		return errors.ErrQueryExecution
	}

	return nil
}

func (s *Storager) Close() error {
	err := s.db.Close(context.Background())

	if err != nil {
		return fmt.Errorf("postgres close: %w", err)
	}

	return nil
}
