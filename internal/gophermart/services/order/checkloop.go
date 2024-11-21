package order

import (
	"context"
	"errors"
	"time"

	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/logger"
	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/model"
	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/repository"
)

const (
	TIMESLEEPLOOP = 300 * time.Millisecond
)

func (m *Manager) updateOrdersInfoLoop() {
	logger.Info("update order loop up")
	defer logger.Info("update orders loop down")

	ticker := time.NewTicker(TIMESLEEPLOOP)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			orders, err := m.repository.GetNotComplitedOrders(context.Background())
			if err != nil {
				if errors.Is(err, repository.ErrOrdersNotFound) {
					continue // Нет заказов, ждем следующий тик.
				}
				logger.Error("fail while get not complited orders, error: %v", err)
				continue
			}

			for _, o := range orders {
				if err := m.updateOrderInfo(&o); err != nil {
					// TODO может поломанным менять статус?
					logger.Error("fail while check and update processing order: %v error: %v", o, err)
				}
			}
		}
	}
}

func (m *Manager) updateOrderInfo(o *model.Order) error {
	var newOrder model.Order
	var accrual float64
	var status int
	var isUpdate bool

	info, err := m.accrual.AccrualReq(o.ID)
	if err != nil {
		return err
	}

	if info.Status != o.Status {
		accrual, status, err = newOrder.ConvertAccrual(info)
		if err != nil {
			return err
		}
		err = m.repository.UpdateOrderStatus(context.Background(), o.ID, status)
		if err != nil {
			return err
		}
		isUpdate = true
	}

	if info.Status == model.PROCESSED && isUpdate {
		err = m.repository.InsertBalanceTransaction(context.Background(), o.UserID, o.ID, accrual)
		if err != nil {
			return err
		}
	}

	return nil
}
