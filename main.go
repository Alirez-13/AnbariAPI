package main

import (
	"AnbariAPI/Internal/auth/handler"
	"AnbariAPI/Internal/auth/middleware"
	"AnbariAPI/Internal/auth/repository"
	"AnbariAPI/Internal/auth/service"
	"AnbariAPI/api/routes"
	"AnbariAPI/shared/database"
	"log"
	"log/slog"
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
	authService := service.NewAuthService(userRepo, sessionRepo, slog.Default())
	authHandler := handler.NewAuthHandler(authService)

	r := gin.Default()

	routes.SetupRoutes(r, db, authHandler)

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
