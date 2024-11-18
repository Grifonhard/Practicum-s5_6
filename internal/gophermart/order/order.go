package order

import (
	"errors"
	"sort"

	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/accrual"
	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/logger"
	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/order/storage"
	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/order/transactions"
	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/repository"
	"github.com/Grifonhard/Practicum-s5_6/internal/model"
)

// TODO запись в логи при возникновении ошибок

type Manager struct {
	s   *storage.Storage
	a   *accrual.Manager
	r   *repository.DB
	muT *transactions.Mutex
}

func New(r *repository.DB, acm *accrual.Manager) (*Manager, error) {
	var m Manager

	stor, err := storage.New(r)
	if err != nil {
		return nil, err
	}
	m.s = stor

	mu, err := transactions.New()
	if err != nil {
		return nil, err
	}
	m.muT = mu

	m.a = acm
	m.r = r

	return &m, nil
}

func (m *Manager) AddOrder(username string, orderID int) error {
	err := checkLuhn(orderID)
	if err != nil {
		return err
	}
	err = m.s.NewOrder(username, orderID)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) ListOrders(username string) ([]model.OrderDto, error) {
	err := m.updateOrdersInfo(username)
	if err != nil {
		logger.Error("fail update orders info: %v", err)
	}

	orders, err := m.s.GetOrders(username)
	if err != nil {
		return nil, err
	}
	sort.Slice(orders, func(i, j int) bool {
		return orders[i].Created.Before(orders[j].Created)
	})

	// часть получения инфы о заказах, по которым ещё нет данных
	// собираем недостающую инфу
	var ordersFront []model.OrderDto
	for _, o := range orders {
		order, err := m.convertToFrontOrder(&o)
		if errors.Is(err, ErrOrderNotReady) {
			logger.Debug("order is still processing: %v", o)
			continue
		}
		if errors.Is(err, ErrOrderInvalid) {
			logger.Debug("order is invalid: %v", o)
			continue
		}
		if err != nil {
			logger.Error("order %+v convert error: %v", o, err)
			return nil, err
		}
		ordersFront = append(ordersFront, *order)
	}

	return ordersFront, err
}

func (m *Manager) Balance(username string) (*model.BalanceDto, error) {
	err := m.updateOrdersInfo(username)
	if err != nil {
		logger.Error("fail update orders info: %v", err)
	}

	ts, err := m.s.GetTransactions(username)
	if err != nil {
		return nil, err
	}

	var sum int
	var withdrawn int

	for _, t := range ts {
		sum += t.Sum
		if sum < 0 {
			withdrawn += t.Sum
		}
	}

	return &model.BalanceDto{
		Current:   sum,
		Withdrawn: withdrawn,
	}, nil
}

func (m *Manager) Withdraw(username, order string, sum int) error {
	m.muT.Lock(username)
	defer m.muT.Unlock(username)

	balance, err := m.Balance(username)
	if err != nil {
		return err
	}
	if balance.Current < sum {
		return ErrNotEnoughBalance
	}

	ts, err := m.s.GetTransactionsByOrder(order)
	if err != nil {
		return err
	}

	switch len(ts) {
	case 1:
		// 1 заказ - 1 списание
	case 2:
		return ErrAlreadyDebited
	default:
		return ErrTooMuchTransact
	}

	return m.s.Withdraw(username, order, sum)
}

func (m *Manager) Withdrawls(username string) ([]model.WithdrawlDto, error) {
	err := m.updateOrdersInfo(username)
	if err != nil {
		logger.Error("fail update orders info: %v", err)
	}

	transs, err := m.s.GetTransactions(username)
	if err != nil {
		return nil, err
	}

	sort.Slice(transs, func(i, j int) bool {
		return transs[i].Created.Before(transs[j].Created)
	})

	var result []model.WithdrawlDto

	for _, t := range transs {
		if t.Sum < 0 {
			withdrawl := model.GetWithdrawFront(t.OrderId, t.Sum*(-1), t.Created)
			result = append(result, *withdrawl)
		}
	}

	return result, nil
}

func (m *Manager) convertToFrontOrder(o *model.Order) (*model.OrderDto, error) {
	if o.Status == model.NEW || o.Status == model.PROCESSING {
		return nil, ErrOrderNotReady
	}
	if o.Status == model.INVALID {
		return nil, ErrOrderInvalid
	}

	var orderFront model.OrderDto
	var accrual int

	transs, err := m.r.GetTransactionsByOrder(o.Id)
	if err != nil {
		return nil, err
	}

	for _, t := range transs {
		accrual += t.Sum
	}

	err = orderFront.ConvertOrder(o, accrual)
	if err != nil {
		return nil, err
	}

	return &orderFront, nil
}

func (m *Manager) updateOrdersInfo(username string) error {
	orders, err := m.s.GetNotComplitedOrders(username)
	if err != nil {
		return err
	}

	for _, o := range orders {
		if err = m.updateOrderInfo(&o); err != nil {
			// TODO может поломанные игнорить или менять статус?
			logger.Error("fail while check and update processing order: %v error: %v", o, err)
			return err
		}
	}

	return nil
}

func (m *Manager) updateOrderInfo(o *model.Order) error {
	var newOrder model.Order
	var accrual, status int
	var isUpdate bool

	info, err := m.a.AccrualReq(o.Id)
	if err != nil {
		return err
	}

	if info.Status != o.Status {
		accrual, status, err = newOrder.ConvertAccrual(info)
		if err != nil {
			return err
		}
		err = m.r.UpdateOrderStatus(o.Id, status)
		if err != nil {
			return err
		}
		isUpdate = true
	}

	if info.Status == model.PROCESSED && isUpdate {
		err = m.r.InsertBalanceTransaction(o.UserId, o.Id, accrual)
		if err != nil {
			return err
		}
	}

	return nil
}

func checkLuhn(orderId int) error {
	var sum int
	shouldDouble := false

	for orderId > 0 {
		digit := orderId % 10

		if shouldDouble {
			digit <<= 1
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		shouldDouble = !shouldDouble
		orderId /= 10
	}

	if sum%10 == 0 {
		return nil
	}

	return ErrLuhnFail
}
