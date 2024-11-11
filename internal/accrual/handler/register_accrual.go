package handler

import (
	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/errors"
	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/model"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type accrualRegistrationRequest struct {
	Match      string `json:"match" required:"true"`
	Reward     int64  `json:"reward" required:"true"`
	RewardType string `json:"reward_type" required:"true"`
}

func (h *Handler) AccrualRegistrationHandler(c *gin.Context) {
	var req accrualRegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accrual := model.AccrualProgram{
		Match:      req.Match,
		Reward:     req.Reward,
		RewardType: req.RewardType,
	}

	ctx := c.Request.Context()
	err := h.AccrualService.RegisterAccrual(ctx, accrual)

	if err != nil {
		slog.ErrorContext(ctx, "accrual registration handle", "err", err)
		c.JSON(errors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"accrual": accrual,
	})
}
