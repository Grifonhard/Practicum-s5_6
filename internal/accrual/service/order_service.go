package service

import (
	"context"
	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/model"
)

type orderRepository interface {
	RegisterOrder(context.Context, uint64, []model.Good) error
	GetOrderByNumber(context.Context, uint64) (*model.Order, error)
	UpdateOrderAccrual(context.Context, uint64, float64) error
	UpdateOrderStatus(context.Context, uint64, string) error
	GetAllOrders(context.Context) ([]model.Order, error)
}

type OrderService struct {
	repo orderRepository
}

func NewOrderService(repo orderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) RegisterOrder(ctx context.Context, number uint64, goods []model.Good) error {
	return s.repo.RegisterOrder(ctx, number, goods)
}

func (s *OrderService) GetOrderByNumber(ctx context.Context, number uint64) (*model.Order, error) {
	return s.repo.GetOrderByNumber(ctx, number)
}

func (s *OrderService) GetAllOrders(ctx context.Context) ([]model.Order, error) {
	return s.repo.GetAllOrders(ctx)
}
