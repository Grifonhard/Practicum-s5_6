package service

import (
	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/model"
	"github.com/Grifonhard/Practicum-s5_6/internal/lib/math"
)

func CalculateReward(good model.Good, accrual model.AccrualProgram) float64 {
	var result float64

	switch accrual.RewardType {
	case model.RewardTypePoints:
		result = float64(accrual.Reward)
	case model.RewardTypePercent:
		reward := math.Percent(int(accrual.Reward), int(good.Price))
		result = float64(reward)
	}

	return result
}
