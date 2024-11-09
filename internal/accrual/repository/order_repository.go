package repository

import (
    "context"
    "github.com/jackc/pgx/v5"

    "github.com/Grifonhard/Practicum-s5_6/internal/accrual/errors"
    "github.com/Grifonhard/Practicum-s5_6/internal/accrual/model"
)

type OrderRepository struct {
    DB *pgx.Conn
}

func NewUserRepository(db *pgx.Conn) model.OrderRepository {
    return &OrderRepository{
        DB: db,
    }
}

func (r *OrderRepository) Create(ctx context.Context, o *model.Order) error {
    query := "INSERT INTO orders (number, status, accrual) VALUES ($1, $2, $3)"
    _, err := r.DB.Exec(ctx, query, o.Number, o.Status, o.Accrual)
    if err != nil {
        return errors.ErrQueryExecution
    }

    return nil

    return nil
}
