package service

import (
	"context"
	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/model"
	"github.com/Grifonhard/Practicum-s5_6/internal/lib/helpers"
	"log/slog"
	"strings"
	"time"
)

const (
	accrualsWorkerDelay = 500 * time.Millisecond
)

type accrualRepository interface {
	RegisterAccrual(context.Context, model.AccrualProgram) error
	GetAllAccrualPrograms(context.Context) ([]model.AccrualProgram, error)
}

type goodRepository interface {
	GetGoodsOfOrdersByStatus(context.Context, string) ([]model.Good, error)
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

func (s *AccrualService) RunAccrualsWorker(ctx context.Context) {
	go s.runAccrualsWorker(ctx)
}

func (s *AccrualService) runAccrualsWorker(ctx context.Context) {
	ticker := time.NewTicker(accrualsWorkerDelay)
	defer ticker.Stop()

	for range ticker.C {
		select {
		case <-ctx.Done():
			return
		default:
		}

		goods, err := s.goodRepo.GetGoodsOfOrdersByStatus(ctx, model.OrderStatusRegistered)
		if err != nil {
			slog.ErrorContext(ctx, "get registered orders", "err", err)
		}

		groupedGoods := helpers.GroupBy(goods, func(g model.Good) uint64 { return g.OrderNumber })

		accruals, err := s.accrualRepo.GetAllAccrualPrograms(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "get all accrual programs", "err", err)
		}

		for number, goodsList := range groupedGoods {
			err = s.orderRepo.UpdateOrderStatus(ctx, number, model.OrderStatusProcessing)
			if err != nil {
				slog.ErrorContext(ctx, "update order status", "err", err)
			}

			var orderAccrual float64
			for _, accrual := range accruals {
				match := strings.ToLower(accrual.Match)

				filteredGoods := helpers.Filter(goodsList, func(good model.Good, _ int) bool {
					desc := strings.ToLower(good.Description)
					return strings.Contains(desc, match)
				})

				for _, good := range filteredGoods {
					reward := CalculateReward(good, accrual)
					orderAccrual += reward
				}
			}

			err = s.orderRepo.UpdateOrderAccrual(ctx, number, orderAccrual)
			if err != nil {
				updateErr := s.orderRepo.UpdateOrderStatus(ctx, number, model.OrderStatusInvalid)
				if updateErr != nil {
					slog.ErrorContext(ctx, "update order status", "err", updateErr)
				}
			} else {
				slog.InfoContext(ctx, "update order accrual", "order", number, "accrual", orderAccrual)
			}
		}
	}
}
