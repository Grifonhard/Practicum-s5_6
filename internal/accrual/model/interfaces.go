package model

import (
	"context"
)

type OrderService interface {
	CreateOrder(ctx context.Context, o *Order) error
}

type OrderRepository interface {
	Create(ctx context.Context, u *Order) error
}
