package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/errors"
	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/model"
)

type orderRegistrationRequest struct {
	Order uint64       `json:"order"`
	Goods []model.Good `json:"goods"`
}

func (h *Handler) OrderRegistrationHandler(c *gin.Context) {
	var req orderRegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	err := h.OrderService.RegisterOrder(ctx, req.Order, req.Goods)

	if err != nil {
		slog.ErrorContext(ctx, "order registration handle", "err", err)
		c.JSON(errors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"number": req.Order,
	})
}
