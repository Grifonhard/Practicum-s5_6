package handler

import (
	"time"

	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/model"
	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/service"
	"github.com/gin-gonic/gin"
)

// Handler struct holds required services for handler to function
type Handler struct {
	OrderService model.OrderService
}

// Config will hold services that will eventually be injected into this
// handler layer on handler initialization
type Config struct {
	Router          *gin.Engine
	OrderService    service.OrderService
	TimeoutDuration time.Duration
}

// NewHandler initializes the handler with required injected services along with http routes
// Does not return as it deals directly with a reference to the gin Engine
func NewHandler(c *Config) {
	h := &Handler{
		OrderService: &c.OrderService,
	}

	g := c.Router.Group("/api")

	g.POST("/user/orders", h.CreateOrderHandler)
}
