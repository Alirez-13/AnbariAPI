package main

import (
	"AnbariAPI/database"
	"AnbariAPI/handler"
	"AnbariAPI/middleware"
	"AnbariAPI/repository"
	"AnbariAPI/routes"
	"AnbariAPI/service"
	"log"
	"os"

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

	routes.SetupRoutes(r)

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

	port := getPort()
	log.Printf("Server starting on port %s", port)
	r.Run(":" + port)
}

func getPort() string {
	if port := os.Getenv("PORT"); port != "" {
		return port
	}
	return "8080"
}
