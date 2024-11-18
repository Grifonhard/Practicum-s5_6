package service

import (
	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/model"
	"github.com/Grifonhard/Practicum-s5_6/internal/lib/math"
)

func CalculateReward(good model.Good, accrual model.AccrualProgram) uint64 {
	var result uint64

	switch accrual.RewardType {
	case model.RewardTypePoints:
		result = uint64(accrual.Reward)
	case model.RewardTypePercent:
		reward := math.Percent(int(accrual.Reward), int(good.Price))
		result = uint64(reward)
	}

	return result
}
