package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"time"

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

func (r *OrderRepository) RegisterOrder(ctx context.Context, number uint64, goods []model.Good) error {
	createOrderQuery := "INSERT INTO accrual.orders (number, status, created_at) VALUES ($1, $2, $3)"
	createGoodQuery := "INSERT INTO accrual.goods (description, price, order_number, created_at) VALUES ($1, $2, $3, $4)"

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	defer tx.Rollback(ctx)

	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	_, err = tx.Exec(ctx, createOrderQuery, number, model.OrderStatusRegistered, time.Now())
	if err != nil {
		return fmt.Errorf("insert order: %w", err)
	}

	for _, g := range goods {
		_, err = tx.Exec(ctx, createGoodQuery, g.Description, g.Price, number, time.Now())
		if err != nil {
			return fmt.Errorf("insert good: %w", err)
		}
	}

	return tx.Commit(ctx)
}
