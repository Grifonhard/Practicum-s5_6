package service

import (
	"context"
	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/model"
)

type orderRepository interface {
	RegisterOrder(context.Context, uint64, []model.Good) error
	GetOrderByNumber(context.Context, uint64) (*model.Order, error)
}

type OrderService struct {
	repo orderRepository
}

func NewOrderService(repo orderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (u *OrderService) RegisterOrder(ctx context.Context, number uint64, goods []model.Good) error {
	return u.repo.RegisterOrder(ctx, number, goods)
}

func (u *OrderService) GetOrderByNumber(ctx context.Context, number uint64) (*model.Order, error) {
	return u.repo.GetOrderByNumber(ctx, number)
}
