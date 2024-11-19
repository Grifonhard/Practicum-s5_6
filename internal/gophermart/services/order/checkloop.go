package order

import (
	"time"

	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/logger"
	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/model"
)

const (
	TIMESLEEPLOOP = 300 * time.Millisecond
)

func (m *Manager) updateOrdersInfoLoop() {
	logger.Info("update order loop up")
	defer logger.Info("update orders loop down")
	for {
		orders, err := m.repository.GetNotComplitedOrders()
		if err != nil {
			logger.Error("fail while get not complited orders, error: %v", err)
		}
	
		for _, o := range orders {
			if err = m.updateOrderInfo(&o); err != nil {
				// TODO может поломанным менять статус?
				logger.Error("fail while check and update processing order: %v error: %v", o, err)
			}
		}
	
		time.Sleep(TIMESLEEPLOOP)
	}
}

func (m *Manager) updateOrderInfo(o *model.Order) error {
	var newOrder model.Order
	var accrual, status int
	var isUpdate bool

	info, err := m.accrual.AccrualReq(o.Id)
	if err != nil {
		return err
	}

	if info.Status != o.Status {
		accrual, status, err = newOrder.ConvertAccrual(info)
		if err != nil {
			return err
		}
		err = m.repository.UpdateOrderStatus(o.Id, status)
		if err != nil {
			return err
		}
		isUpdate = true
	}

	if info.Status == model.PROCESSED && isUpdate {
		err = m.repository.InsertBalanceTransaction(o.UserId, o.Id, accrual)
		if err != nil {
			return err
		}
	}

	return nil
}
