package main

import (
	"AnbariAPI/database"
	"AnbariAPI/handler"
	"AnbariAPI/middleware"
	"AnbariAPI/repository"
	"AnbariAPI/service"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	config := database.DefaultConfig()
	config.DBPath = "./data/inventory.db"
	config.Debug = true

	if err := database.InitDatabase(config); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.CloseDatabase()

	db := database.GetDB()

	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	authService := service.NewAuthService(userRepo, sessionRepo)
	authHandler := handler.NewAuthHandler(authService)

	r := gin.Default()

	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/logout", authHandler.Logout)
	}

	protected := r.Group("/api")
	protected.Use(middleware.SessionAuth(authService))
	{
		protected.GET("/me", authHandler.GetCurrentUser)
	}

	r.Run(":8080")
}