package repository

import (
	"context"
	"fmt"
	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GoodRepository struct {
	connPool *pgxpool.Pool
}

func NewGoodRepository(connPool *pgxpool.Pool) *GoodRepository {
	return &GoodRepository{
		connPool: connPool,
	}
}

func (r *GoodRepository) GetGoodsOfOrdersByStatus(ctx context.Context, orderStatus string) ([]model.Good, error) {
	query := `
        SELECT goods.id, goods.description, goods.price, goods.order_number, goods.created_at
        FROM accrual.orders AS orders
        JOIN accrual.goods AS goods ON goods.order_number = orders.number
        WHERE orders.status = $1
    `
	conn, err := r.connPool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("acquire connection: %w", err)
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, query, orderStatus)
	if err != nil {
		return nil, fmt.Errorf("select goods: %w", err)
	}
	defer rows.Close()

	goods, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Good])
	if err != nil {
		return nil, fmt.Errorf("collect rows: %w", err)
	}

	return goods, nil
}
