package main

import (
	"fmt"
	"log"

	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/auth"
	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/storage"
	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/web"
	"github.com/gin-gonic/gin"
)

func main() {
	stor, err := storage.New()
	if err != nil {
		log.Fatal(err)
	}
	authManager, err := auth.New(stor)
	if err != nil {
		log.Fatal(err)
	}

	router := initRouter(authManager)

	fmt.Println("Server start")
	log.Fatal(router.Run("localhost:8080"))
}

func initRouter(am *auth.Manager) *gin.Engine {
	router := gin.Default()

	router.POST("/api/user/register", web.Registration(am))
	router.POST("/api/user/login", web.Login(am))

	return router
}
