package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/errors"
	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/model"
)

type newOrderRegistrationRequest struct {
	Number uint64 `json:"number"`
}

func (h *Handler) NewOrderRegistrationHandler(c *gin.Context) {
	var req newOrderRegistrationRequest

	//if ok := bindData(c, &req); !ok {
	//    return
	//}

	o := &model.Order{
		Number: req.Number,
	}

	ctx := c.Request.Context()
	err := h.OrderService.CreateOrder(ctx, o)

	if err != nil {
		slog.ErrorContext(ctx, "new order registration", "err", err)
		c.JSON(errors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"number": o.Number,
	})
}
