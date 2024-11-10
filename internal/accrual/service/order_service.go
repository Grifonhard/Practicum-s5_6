package service

import (
	"context"
	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/model"
)

type orderRepository interface {
	Create(ctx context.Context, o *model.Order) error
}

type OrderService struct {
	repo orderRepository
}

func NewOrderService(repo orderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (u *OrderService) CreateOrder(ctx context.Context, o *model.Order) error {
	return u.repo.Create(ctx, o)
}
