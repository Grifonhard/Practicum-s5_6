package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/services/http/accrual"
	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/logger"
	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/repository"
	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/services/auth"
	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/services/order"
	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/services/transactions"
	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/web"
	"github.com/caarlos0/env/v10"
	"github.com/gin-gonic/gin"
)

const (
	NODATA = "no data"
)

type CFG struct {
	Address *string `env:"RUN_ADDRESS"`
	DBURI *string `env:"DATABASE_URI"`
}

func main() {
	address := flag.String("a", NODATA, "адрес гофемарта")
	uri := flag.String("d", NODATA, "адрес db")

	flag.Parse()

	var cfg CFG
	err := env.Parse(&cfg)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	err = logger.InitLogger(logger.DEBUG, "")
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	if cfg.Address != nil {
		address = cfg.Address
	}
	if cfg.DBURI != nil {
		uri = cfg.DBURI
	}

	am, om, err := initServices(uri)
	if err != nil {
		logger.Error("init services error: %v", err)
		log.Fatal(err)
	}

	router := initRouter(am, om)

	logger.Debug("env addr: %s, uri: %s; flag addr: %s, uri: %s", *cfg.Address, *cfg.DBURI, *address, *uri)
	fmt.Println("Server start")
	log.Fatal(router.Run(*address))
}

func initRouter(am *auth.Manager, om *order.Manager) *gin.Engine {
	router := gin.Default()

	router.POST("/api/user/register", web.Registration(am))
	router.POST("/api/user/login", web.Login(am))

	router.POST("/api/user/orders", web.Authentication(am), web.AddOrder(om))
	router.GET("/api/user/orders", web.Authentication(am), web.ListOrders(om))
	router.GET("/api/user/balance", web.Authentication(am), web.Balance(om))
	router.POST("/api/user/balance/withdraw", web.Authentication(am), web.Withdraw(om))
	router.GET("/api/user/withdrawals", web.Authentication(am), web.Withdrawals(om))

	return router
}

func initServices(uri *string) (*auth.Manager, *order.Manager, error) {
	transMu, err := transactions.New()
	if err != nil {
		return nil, nil, err
	}
	db, err := repository.New(*uri)
	if err != nil {
		return nil, nil, err
	}
	acc, err := accrual.New(*uri)
	if err != nil {
		return nil, nil, err
	}

	authM, err := auth.New(db, transMu)
	if err != nil {
		return nil, nil, err
	}
	orderM, err := order.New(db, transMu, acc)
	if err != nil {
		return nil, nil, err
	}

	return authM, orderM, nil
}
