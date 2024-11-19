package order

import (
	"errors"
	"fmt"
	"sort"
	"strconv"

	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/logger"
	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/model"
	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/repository"
	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/services/http/accrual"
	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/services/transactions"
)

// TODO поля полными именами

type Manager struct {
	accrual   *accrual.Manager
	repository   *repository.DB
	muTransaction *transactions.Mutex
}

func New(r *repository.DB, t *transactions.Mutex, acm *accrual.Manager) (*Manager, error) {
	var m Manager

	m.muTransaction = t
	m.accrual = acm
	m.repository = r

	go m.updateOrdersInfoLoop()

	return &m, nil
}

func (m *Manager) AddOrder(userId int, orderID int) error {
	err := checkLuhn(orderID)
	if err != nil {
		return err
	}
	err = m.repository.InsertOrder(userId, orderID)

	if errors.Is(err, repository.ErrOrderExist) {
		order, err := m.repository.GetOrder(orderID)
		if err != nil {
			logger.Error("fail while get order: %v", err)
			return err
		}
		if order.UserId == userId {
			return ErrOrderExistThis
		} else {
			return fmt.Errorf("%w user id: %d", ErrOrderExistAnother, order.UserId)
		}
	}

	return nil
}

func (m *Manager) ListOrders(userId int) ([]model.OrderDto, error) {
	orders, err := m.repository.GetOrders(userId)
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

func (m *Manager) Balance(userId int) (*model.BalanceDto, error) {

	ts, err := m.repository.GetTransactions(userId)
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

func (m *Manager) Withdraw(userId int, order string, sum int) error {
	m.muTransaction.Lock(strconv.Itoa(userId))
	defer m.muTransaction.Unlock(strconv.Itoa(userId))

	balance, err := m.Balance(userId)
	if err != nil {
		return err
	}
	if balance.Current < sum {
		return ErrNotEnoughBalance
	}

	orderInt, err := strconv.Atoi(order)
	if err != nil {
		return err
	}

	ts, err := m.repository.GetTransactionsByOrder(orderInt)
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

	// списания - это транзакции со знаком -
	sum *= (-1)

	return m.repository.InsertBalanceTransaction(userId, orderInt, sum)
}

func (m *Manager) Withdrawls(userId int) ([]model.WithdrawlDto, error) {

	transs, err := m.repository.GetTransactions(userId)
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

	transs, err := m.repository.GetTransactionsByOrder(o.Id)
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
