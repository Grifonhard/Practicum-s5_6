package server

import (
	"context"
	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/config"
	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/handler"
	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/repository"
	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/service"
	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/storage/postgres"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	parseFlags()
	initConfig()

	httpCfg := config.HTTPConfig{RunAddress: flagRunAddress}
	pgCfg := config.PostgresConfig{
		DatabaseURI:    flagDatabaseURI,
		ConnectTimeout: 5 * time.Second,
	}

	db, err := postgres.NewConnection(&pgCfg)
	if err != nil {
		panic(err)
	}

	err = postgres.CreateSchema(ctx, db)
	if err != nil {
		panic(err)
	}

	err = postgres.CreateTables(ctx, db)
	if err != nil {
		panic(err)
	}

	orderRepository := repository.NewOrderRepository(db)
	accrualRepository := repository.NewAccrualRepository(db)
	goodRepository := repository.NewGoodRepository(db)
	orderService := service.NewOrderService(orderRepository)
	accrualService := service.NewAccrualService(accrualRepository, orderRepository, goodRepository)

	router := gin.Default()

	handler.NewHandler(router, orderService, accrualService)

	slog.Info("Connected to database", "uri", pgCfg.DatabaseURI)
	slog.Info("Start server", "address", httpCfg.RunAddress)

	accrualService.CalculateAccruals(ctx)

	err = http.ListenAndServe(httpCfg.RunAddress, router)

	if err != nil {
		panic(err)
	}
}

func initConfig() {
	if envRunAddress := os.Getenv("RUN_ADDRESS"); envRunAddress != "" {
		flagRunAddress = envRunAddress
	}

	if envDatabaseURI := os.Getenv("DATABASE_URI"); envDatabaseURI != "" {
		flagDatabaseURI = envDatabaseURI
	}
}
