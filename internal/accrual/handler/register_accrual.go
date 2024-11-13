package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"

	errs "github.com/Grifonhard/Practicum-s5_6/internal/accrual/errors"
	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/model"
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

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == errs.ErrPostgresUniqueViolation {
			c.JSON(http.StatusConflict, gin.H{
				"error": errs.NewConflict(err),
			})
			return
		}

		c.JSON(errs.Status(err), gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"accrual": accrual,
	})
}
