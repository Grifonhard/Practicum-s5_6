package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"

	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/model"
)

type AccrualRepository struct {
	connPool *pgxpool.Pool
}

func NewAccrualRepository(connPool *pgxpool.Pool) *AccrualRepository {
	return &AccrualRepository{
		connPool: connPool,
	}
}

func (r *AccrualRepository) RegisterAccrual(ctx context.Context, accrual model.AccrualProgram) error {
	query := "INSERT INTO accrual.accrual_programs (match, reward, reward_type, created_at) VALUES ($1, $2, $3, $4)"

	conn, err := r.connPool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("acquire connection: %w", err)
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	defer tx.Rollback(ctx)

	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	_, err = tx.Exec(ctx, query, accrual.Match, accrual.Reward, accrual.RewardType, time.Now())
	if err != nil {
		return fmt.Errorf("insert accrual: %w", err)
	}

	return tx.Commit(ctx)
}

func (r *AccrualRepository) GetAllAccrualPrograms(ctx context.Context) ([]model.AccrualProgram, error) {
	query := "SELECT * FROM accrual.accrual_programs"
	conn, err := r.connPool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("acquire connection: %w", err)
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("select accruals: %w", err)
	}
	defer rows.Close()

	accruals, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.AccrualProgram])
	if err != nil {
		return nil, fmt.Errorf("collect rows: %w", err)
	}

	return accruals, nil
}
