package storage

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/logger"
	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/repository"
	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/model"
)

// TODO service layers
type Storage struct {
	db *repository.DB
}

func New(db *repository.DB) (*Storage, error) {
	var stor Storage
	stor.db = db
	return &stor, nil
}

func (s *Storage) NewOrder(username string, orderId int) error {
	user, err := s.db.GetUser(username)
	if err != nil {
		return err
	}
	err = s.db.InsertOrder(user.Id, orderId)

	if errors.Is(err, repository.ErrOrderExist) {
		order, err := s.db.GetOrder(orderId)
		if err != nil {
			logger.Error("fail while get order: %v", err)
			return err
		}
		if order.UserId == user.Id {
			return ErrOrderExistThis
		} else {
			return fmt.Errorf("%w user id: %d", ErrOrderExistAnother, order.UserId)
		}
	}

	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetOrders(username string) ([]model.Order, error) {
	user, err := s.db.GetUser(username)
	if err != nil {
		return nil, err
	}
	orders, err := s.db.GetOrders(user.Id)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *Storage) GetNotComplitedOrders(username string) ([]model.Order, error) {
	user, err := s.db.GetUser(username)
	if err != nil {
		return nil, err
	}
	orders, err := s.db.GetNotComplitedOrders(user.Id)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *Storage) GetTransactions(username string) ([]model.BalanceTransactions, error) {
	u, err := s.db.GetUser(username)
	if err != nil {
		return nil, err
	}

	transs, err := s.db.GetTransactions(u.Id)
	if err != nil {
		return nil, err
	}
	return transs, nil
}

func (s *Storage) GetTransactionsByOrder(order string) ([]model.BalanceTransactions, error) {
	orderint, err := strconv.Atoi(order)
	if err != nil {
		return nil, err
	}

	ts, err := s.db.GetTransactionsByOrder(orderint)
	if err != nil {
		return nil, err
	}

	return ts, nil
}

func (s *Storage) Withdraw(username, order string, sum int) error {
	user, err := s.db.GetUser(username)
	if err != nil {
		return err
	}

	orderint, err := strconv.Atoi(order)
	if err != nil {
		return err
	}

	sum *= (-1)

	return s.db.InsertBalanceTransaction(user.Id, orderint, sum)
}
