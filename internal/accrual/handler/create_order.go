package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/errors"
	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/model"
)

type createOrderRequest struct {
	Number uint64 `json:"number"`
}

func (h *Handler) CreateOrderHandler(c *gin.Context) {
	var req createOrderRequest

	//if ok := bindData(c, &req); !ok {
	//    return
	//}

	o := &model.Order{
		Number: req.Number,
	}

	ctx := c.Request.Context()
	err := h.OrderService.CreateOrder(ctx, o)

	if err != nil {
		log.Printf("Failed to sign in user: %v\n", err.Error())
		c.JSON(errors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"number": o.Number,
	})
}
