package handler

import (
	errs "github.com/Grifonhard/Practicum-s5_6/internal/lib/errors"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

func (h *Handler) GetOrdersHandler(c *gin.Context) {
	ctx := c.Request.Context()
	orders, err := h.OrderService.GetRegisteredOrdersWithGoods(ctx)

	if err != nil {
		slog.ErrorContext(ctx, "get orders handle", "err", err)

		c.JSON(errs.Status(err), gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
	})
}
