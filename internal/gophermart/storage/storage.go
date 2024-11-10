package storage

import (
	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/drivers/psql"
	"github.com/Grifonhard/Practicum-s5_6/internal/model"
)

type Storage struct {
	db *psql.DB
}

func New(uri string) (*Storage, error) {
	var stor Storage
	var err error
	stor.db, err = psql.New(uri)
	if err != nil {
		return nil, err
	}
	return &stor, nil
}

func (s *Storage) NewUser(user model.User) error {
	err := s.db.InsertUser(user.Username, user.Password_hash)
	return err
}

func (s *Storage) GetUser(uname string) (*model.User, error) {
	user, err := s.db.GetUser(uname)
	if err != nil {
		return nil, err
	}
	return user, nil
}

