package storage

import (
	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/repository"
	"github.com/Grifonhard/Practicum-s5_6/internal/model"
)

// TODO service layers
type Storage struct {
	db *repository.DB
}

func New(db *repository.DB) (*Storage, error) {
	var stor Storage
	stor.db =db
	return &stor, nil
}

func (s *Storage) NewOrder(username string, orderId int) error {
	user, err := s.db.GetUser(username)
	if err != nil {
		return err
	}
	err = s.db.InsertOrder(user.Id, orderId)
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