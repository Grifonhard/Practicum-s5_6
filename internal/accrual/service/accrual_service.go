package service

import (
	"context"
	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/model"
	"github.com/Grifonhard/Practicum-s5_6/internal/lib/math"
	"log/slog"
	"strings"
)

type accrualRepository interface {
	RegisterAccrual(context.Context, model.AccrualProgram) error
	GetAllAccrualPrograms(context.Context) ([]model.AccrualProgram, error)
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

func (s *AccrualService) CalculateAccruals(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			orders, err := s.orderRepo.GetRegisteredOrdersWithGoods(ctx)
			if err != nil {
				slog.ErrorContext(ctx, "get registered orders", "err", err)
			}

			accruals, err := s.accrualRepo.GetAllAccrualPrograms(ctx)
			if err != nil {
				slog.ErrorContext(ctx, "get all accrual programs", "err", err)
			}

			for _, order := range orders {
				err = s.orderRepo.UpdateOrderStatus(ctx, order.Number, model.OrderStatusProcessing)
				if err == nil {
					slog.ErrorContext(ctx, "update order status", "err", err)
				}

				var orderAccrual uint64

				for _, accrual := range accruals {
					match := strings.ToLower(accrual.Match)

					matchedGoods := selectMatchedGoods(order.Goods, match)

					for _, good := range matchedGoods {
						switch accrual.RewardType {
						case model.RewardTypePoints:
							orderAccrual += uint64(accrual.Reward)
						case model.RewardTypePercent:
							reward := math.Percent(int(accrual.Reward), int(good.Price))
							orderAccrual += uint64(reward)
						}
					}
				}

				err = s.orderRepo.UpdateOrderAccrual(ctx, order.Number, orderAccrual)
				if err != nil {
					updateErr := s.orderRepo.UpdateOrderStatus(ctx, order.Number, model.OrderStatusInvalid)
					if updateErr != nil {
						slog.ErrorContext(ctx, "update order status", "err", updateErr)
					}
				}
			}
		}
	}()
}

func selectMatchedGoods(goods []model.Good, match string) []model.Good {
	result := make([]model.Good, 0, len(goods))

	for _, good := range goods {
		desc := strings.ToLower(good.Description)
		if strings.Contains(desc, match) {
			result = append(result, good)
		}
	}

	return result
}
