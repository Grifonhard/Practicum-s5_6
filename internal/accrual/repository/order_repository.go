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

func (r *OrderRepository) GetOrderByNumber(ctx context.Context, number uint64) (*model.Order, error) {
	query := "SELECT * FROM accrual.orders WHERE number = ($1) LIMIT 1"

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	defer tx.Rollback(ctx)

	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}

	_, err = tx.Exec(ctx, query, number, number)
	if err != nil {
		return nil, fmt.Errorf("select order: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return &model.Order{}, nil
}

func (r *OrderRepository) ListRegisteredOrders(ctx context.Context) ([]model.Order, error) {
	query := "SELECT * FROM accrual.orders WHERE status = ($1)"

	rows, err := r.db.Query(ctx, query, model.OrderStatusRegistered)
	if err != nil {
		return nil, fmt.Errorf("select orders: %w", err)
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		order := model.Order{}
		err = rows.Scan(&order.Number, &order.Status, &order.Accrual, &order.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("unable to scan row: %w", err)
		}
		orders = append(orders, order)
	}

	return orders, nil
}
