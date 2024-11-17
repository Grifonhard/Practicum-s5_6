package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/model"
	errs "github.com/Grifonhard/Practicum-s5_6/internal/lib/errors"
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
		"number": req.Order,
	})
}
