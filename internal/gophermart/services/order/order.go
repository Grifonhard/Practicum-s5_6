package order

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/logger"
	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/model"
	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/repository"
	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/services/http/accrual"
	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/services/transactions"
)

// TODO поля полными именами

type Manager struct {
	accrual       *accrual.Manager
	repository    *repository.DB
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

func (m *Manager) AddOrder(userID int, orderID int) error {

	err := checkLuhn(orderID)
	if err != nil {
		return err
	}

	err = m.repository.InsertOrder(context.Background(), userID, orderID)
	if errors.Is(err, repository.ErrOrderExist) {
		order, err := m.repository.GetOrder(context.Background(), orderID)
		if err != nil {
			logger.Error("fail while get order: %v", err)
			return err
		}
		if order.UserID == userID {
			return ErrOrderExistThis
		} else {
			return fmt.Errorf("%w user id: %d", ErrOrderExistAnother, order.UserID)
		}
	}

	return nil
}

func (m *Manager) ListOrders(userID int) ([]model.OrderDto, error) {

	orders, err := m.repository.GetOrders(context.Background(), userID)
	if err != nil {
		return nil, err
	}

	// часть получения инфы о заказах, по которым ещё нет данных
	// собираем недостающую инфу
	var ordersFront []model.OrderDto
	for _, o := range orders {
		order, err := m.convertToFrontOrder(&o)
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

func (m *Manager) Balance(userID int) (*model.BalanceDto, error) {

	ts, err := m.repository.GetTransactions(context.Background(), userID)
	if err != nil {
		return nil, err
	}

	var sum, withdrawn float64

	for _, t := range ts {
		sum += t.Sum
		if t.Sum < 0 {
			withdrawn -= t.Sum
		}
	}

	return &model.BalanceDto{
		Current:   sum,
		Withdrawn: withdrawn,
	}, nil
}

func (m *Manager) Withdraw(userID int, order string, sum float64) error {

	m.muTransaction.Lock(strconv.Itoa(userID))
	defer m.muTransaction.Unlock(strconv.Itoa(userID))

	balance, err := m.Balance(userID)
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

	// списания - это транзакции со знаком -
	sum *= (-1)

	ts, err := m.repository.GetTransactionsByOrder(context.Background(), orderInt)
	if errors.Is(err, repository.ErrTransNotFound) {
		err = m.repository.InsertOrder(context.Background(), userID, orderInt)
		if err != nil {
			return err
		}
		return m.repository.InsertBalanceTransaction(context.Background(), userID, orderInt, sum)
	}
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

	return m.repository.InsertBalanceTransaction(context.Background(), userID, orderInt, sum)
}

func (m *Manager) Withdrawls(userID int) ([]model.WithdrawlDto, error) {

	transs, err := m.repository.GetTransactions(context.Background(), userID)
	if err != nil {
		return nil, err
	}

	var result []model.WithdrawlDto

	for _, t := range transs {
		if t.Sum < 0 {
			withdrawl := model.GetWithdrawFront(t.OrderID, t.Sum*(-1), t.Created)
			result = append(result, *withdrawl)
		}
	}

	return result, nil
}

func (m *Manager) convertToFrontOrder(o *model.Order) (*model.OrderDto, error) {
	var orderFront model.OrderDto
	var accrual float64
	var err error
	if o.Status == model.NEW || o.Status == model.PROCESSING {
		err = orderFront.ConvertOrder(o, accrual)
		if err != nil {
			return nil, err
		}
		return &orderFront, nil
	}

	if o.Status == model.INVALID {
		return nil, ErrOrderInvalid
	}

	transs, err := m.repository.GetTransactionsByOrder(context.Background(), o.ID)
	if err != nil && !errors.Is(err, repository.ErrTransNotFound) {
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

func checkLuhn(orderID int) error {
	var sum int
	shouldDouble := false

	for orderID > 0 {
		digit := orderID % 10

		if shouldDouble {
			digit <<= 1
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		shouldDouble = !shouldDouble
		orderID /= 10
	}

	if sum%10 == 0 {
		return nil
	}

	return ErrLuhnFail
}
