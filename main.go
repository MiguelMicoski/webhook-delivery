package main

import (
	"awesomeProject/internal/config"
	"awesomeProject/internal/database"
	"awesomeProject/internal/handler"
	"awesomeProject/internal/service"
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := gin.Default()

	userService := service.NewUserService()
	userHandler := handler.NewUserHandler(userService)

	router.GET("/users", userHandler.FindAll)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
