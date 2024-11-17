package service

import (
	"context"
	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/model"
)

type orderRepository interface {
	RegisterOrder(context.Context, uint64, []model.Good) error
	GetOrderByNumber(context.Context, uint64) (*model.Order, error)
	GetRegisteredOrdersWithGoods(context.Context) ([]model.OrderWithGoods, error)
	UpdateOrderAccrual(context.Context, uint64, uint64) error
	UpdateOrderStatus(context.Context, uint64, string) error
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

func (s *OrderService) GetRegisteredOrdersWithGoods(ctx context.Context) ([]model.OrderWithGoods, error) {
	return s.repo.GetRegisteredOrdersWithGoods(ctx)
}
