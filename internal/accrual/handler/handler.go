package handler

import (
	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	Router       *gin.Engine
	OrderService *service.OrderService
}

func NewHandler(router *gin.Engine, orderService *service.OrderService) {
	h := &Handler{
		Router:       router,
		OrderService: orderService,
	}

	g := h.Router.Group("/api")

	g.GET("/orders/:number", h.OrderRegistrationHandler)
	g.POST("/orders", h.OrderRegistrationHandler)
	g.POST("/goods", h.OrderRegistrationHandler)
}
