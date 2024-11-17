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
	db *pgxpool.Pool
}

func NewAccrualRepository(db *pgxpool.Pool) *AccrualRepository {
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

func (r *AccrualRepository) GetAllAccrualPrograms(ctx context.Context) ([]model.AccrualProgram, error) {
	query := "SELECT * FROM accrual.accrual_programs"

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("select accruals: %w", err)
	}

	var accruals []model.AccrualProgram
	for rows.Next() {
		accrual := model.AccrualProgram{}
		err = rows.Scan(&accrual.ID, &accrual.Match, &accrual.Reward, &accrual.RewardType, &accrual.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("unable to scan row: %w", err)
		}
		accruals = append(accruals, accrual)
	}

	rows.Close()

	return accruals, nil
}
