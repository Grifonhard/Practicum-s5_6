package handler

import (
	"github.com/Grifonhard/Practicum-s5_6/internal/lib/errors"
	"github.com/gin-gonic/gin"
	"log/slog"
	"math"
	"net/http"
	"strconv"
)

func (h *Handler) GetOrderHandler(c *gin.Context) {
	num := c.Param("number")
	ctx := c.Request.Context()

	number, err := strconv.ParseUint(num, 10, 64)
	if err != nil {
		slog.ErrorContext(ctx, "get order handle", "err", err)
		c.JSON(errors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	order, err := h.OrderService.GetOrderByNumber(ctx, number)

	if err != nil {
		slog.ErrorContext(ctx, "get order handle", "err", err)
		c.JSON(errors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	resp := gin.H{
		"order":  strconv.FormatUint(order.Number, 10),
		"status": order.Status,
	}

	if order.Accrual != nil {
		orderAccrual := math.Round(*order.Accrual*100) / 100
		resp["accrual"] = orderAccrual
	}

	c.JSON(http.StatusOK, resp)
}
