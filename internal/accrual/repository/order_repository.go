package repository

import (
	"context"
	"github.com/jackc/pgx/v5"

	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/errors"
	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/model"
)

type OrderRepository struct {
	db *pgx.Conn
}

func NewOrderRepository(db *pgx.Conn) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) Create(ctx context.Context, o *model.Order) error {
	query := "INSERT INTO orders (number, status, accrual) VALUES ($1, $2, $3)"
	_, err := r.db.Exec(ctx, query, o.Number, o.Status, o.Accrual)
	if err != nil {
		return errors.ErrQueryExecution
	}

	return nil
}
