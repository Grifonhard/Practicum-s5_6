package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"

	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/model"
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
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

func (r *OrderRepository) GetOrderByNumber(ctx context.Context, number uint64) (*model.Order, error) {
	query := "SELECT * FROM accrual.orders WHERE number = ($1) LIMIT 1"

	rows, err := r.db.Query(ctx, query, number)
	if err != nil {
		return nil, fmt.Errorf("get order by number: %w", err)
	}
	defer rows.Close()

	order := model.Order{}

	for rows.Next() {
		err = rows.Scan(
			&order.Number,
			&order.Status,
			&order.Accrual,
			&order.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("unable to scan row: %w", err)
		}
	}

	return &order, nil
}

func (r *OrderRepository) UpdateOrderAccrual(ctx context.Context, number uint64, accrual float64) error {
	query := `UPDATE accrual.orders
              SET accrual = $1,
                  status = $2
              WHERE number = $3;`

	_, err := r.db.Query(ctx, query, accrual, model.OrderStatusProcessed, number)
	if err != nil {
		return fmt.Errorf("update order accrual: %w", err)
	}

	return nil
}

func (r *OrderRepository) UpdateOrderStatus(ctx context.Context, number uint64, status string) error {
	query := `UPDATE accrual.orders
              SET status = $1
              WHERE number = $2;`

	_, err := r.db.Query(ctx, query, status, number)
	if err != nil {
		return fmt.Errorf("update order status: %w", err)
	}

	return nil
}

func (r *OrderRepository) GetAllOrders(ctx context.Context) ([]model.Order, error) {
	query := "SELECT * FROM accrual.orders"

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get all orders: %w", err)
	}
	defer rows.Close()

	var orders []model.Order

	for rows.Next() {
		var order model.Order
		err = rows.Scan(
			&order.Number,
			&order.Status,
			&order.Accrual,
			&order.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("unable to scan row: %w", err)
		}
		orders = append(orders, order)
	}

	return orders, nil
}
