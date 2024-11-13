package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"time"

	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/model"
)

type AccrualRepository struct {
	db *pgx.Conn
}

func NewAccrualRepository(db *pgx.Conn) *AccrualRepository {
	return &AccrualRepository{
		db: db,
	}
}

func (r *AccrualRepository) RegisterAccrual(ctx context.Context, accrual model.AccrualProgram) error {
	query := "INSERT INTO accrual.accrual_programs (match, reward, reward_type, created_at) VALUES ($1, $2, $3, $4)"

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
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

func (r *AccrualRepository) MatchAccrualsByGoods(ctx context.Context, accrual model.AccrualProgram) error {
	return nil
}
