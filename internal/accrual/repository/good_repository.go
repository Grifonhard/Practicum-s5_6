package repository

import (
	"context"
	"fmt"
	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GoodRepository struct {
	db *pgxpool.Pool
}

func NewGoodRepository(db *pgxpool.Pool) *GoodRepository {
	return &GoodRepository{
		db: db,
	}
}

func (r *GoodRepository) GetGoodsOfOrdersByStatus(ctx context.Context, orderStatus string) ([]model.Good, error) {
	query := `
        SELECT goods.id, goods.description, goods.price, goods.order_number, goods.created_at
        FROM accrual.orders AS orders
        JOIN accrual.goods AS goods ON goods.order_number = orders.number
        WHERE orders.status = $1
    `

	rows, err := r.db.Query(ctx, query, orderStatus)
	if err != nil {
		return nil, fmt.Errorf("select goods: %w", err)
	}
	defer rows.Close()

	var goods []model.Good
	for rows.Next() {
		good := model.Good{}
		err = rows.Scan(&good.ID, &good.Description, &good.Price, &good.OrderNumber, &good.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("unable to scan row: %w", err)
		}
		goods = append(goods, good)
	}

	return goods, nil
}
