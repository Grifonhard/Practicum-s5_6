package storage

import (
	"context"
	"errors"
	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/model"
	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/storage/postgres"
	"io"
)

type Type string

const (
	TypePostgres Type = "postgres"
)

type Config struct {
	StorageType Type
	Postgres    *postgres.Config
}

type Storager interface {
	CreateOrder(ctx context.Context, u model.Order) error
	io.Closer
}

func NewStorage(cfg *Config) (Storager, error) {
	switch cfg.StorageType {
	case TypePostgres:
		return postgres.NewPostgresStorage(cfg.Postgres)
	default:
		return nil, errors.New("unknown storage type")
	}
}
