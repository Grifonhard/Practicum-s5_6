package service

import (
	"context"
	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/model"
)

type accrualRepository interface {
	RegisterAccrual(context.Context, model.AccrualProgram) error
	MatchAccrualsByGoods(ctx context.Context) error
}

type goodRepository interface {
	GetGoodsByOrderNumbers(context.Context, []uint64) ([]model.Good, error)
}

type AccrualService struct {
	accrualRepo accrualRepository
	orderRepo   orderRepository
	goodRepo    goodRepository
}

func NewAccrualService(
	accrualRepo accrualRepository,
	orderRepo orderRepository,
	goodRepo goodRepository,
) *AccrualService {
	svc := &AccrualService{
		accrualRepo: accrualRepo,
		orderRepo:   orderRepo,
		goodRepo:    goodRepo,
	}

	return svc
}

func (s *AccrualService) RegisterAccrual(ctx context.Context, accrual model.AccrualProgram) error {
	return s.accrualRepo.RegisterAccrual(ctx, accrual)
}

func (s *AccrualService) calculateAccruals(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		orders, _ := s.orderRepo.ListRegisteredOrders(ctx)
		var orderNumbers []uint64

		for _, order := range orders {
			orderNumbers = append(orderNumbers, order.Number)
		}

		_ = s.accrualRepo.MatchAccrualsByGoods(ctx)
	}
}
