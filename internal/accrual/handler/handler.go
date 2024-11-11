package handler

import (
	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	Router         *gin.Engine
	OrderService   *service.OrderService
	AccrualService *service.AccrualService
}

func NewHandler(
	router *gin.Engine,
	orderService *service.OrderService,
	accrualService *service.AccrualService,
) {
	h := &Handler{
		Router:         router,
		OrderService:   orderService,
		AccrualService: accrualService,
	}

	g := h.Router.Group("/api")

	g.GET("/orders/:number", h.GetAccrualsCalculationHandler)
	g.POST("/orders", h.OrderRegistrationHandler)
	g.POST("/goods", h.AccrualRegistrationHandler)
}
