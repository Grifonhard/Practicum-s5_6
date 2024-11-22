package service

import (
	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/model"
	"github.com/Grifonhard/Practicum-s5_6/internal/lib/math"
)

func CalculateReward(good model.Good, accrual model.AccrualProgram) float64 {
	switch accrual.RewardType {
	case model.RewardTypePoints:
		return accrual.Reward
	default:
		return math.Percent(accrual.Reward, good.Price)
	}
}
