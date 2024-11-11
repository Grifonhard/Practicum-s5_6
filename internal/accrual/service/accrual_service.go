package service

import (
	"context"
	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/model"
)

type accrualRepository interface {
	RegisterAccrual(context.Context, model.AccrualProgram) error
}

type AccrualService struct {
	repo accrualRepository
}

func NewAccrualService(repo accrualRepository) *AccrualService {
	return &AccrualService{repo: repo}
}

func (u *AccrualService) RegisterAccrual(ctx context.Context, accrual model.AccrualProgram) error {
	return u.repo.RegisterAccrual(ctx, accrual)
}
