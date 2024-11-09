package service

import (
	"context"
	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/model"
)

type OrderService struct {
	OrderRepository model.OrderRepository
}

type OrderServiceConfig struct {
	OrderRepository model.OrderRepository
}

func NewOrderService(c *OrderServiceConfig) model.OrderService {
	return &OrderService{
		OrderRepository: c.OrderRepository,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, _ *model.Order) error {

	return nil
}
