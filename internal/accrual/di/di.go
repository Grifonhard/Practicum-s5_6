package di

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/handler"
	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/repository"
	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/service"
)

// Init initialize a handler starting from data source
// which inject into repository layer
// which inject into service layer
// which inject into handler layer
func Init(db *pgx.Conn) (*gin.Engine, error) {
	log.Println("Injecting data sources")

	orderRepository := repository.NewUserRepository(db)

	orderService := service.NewOrderService(&service.OrderServiceConfig{
		OrderRepository: orderRepository,
	})

	router := gin.Default()

	handlerTimeout := os.Getenv("HANDLER_TIMEOUT")
	ht, err := strconv.ParseInt(handlerTimeout, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse HANDLER_TIMEOUT as int: %w", err)
	}

	handler.NewHandler(&handler.Config{
		Router:          router,
		OrderService:    orderService,
		TimeoutDuration: time.Duration(ht) * time.Second,
	})

	return router, nil
}
