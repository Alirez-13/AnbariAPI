package main

import (
	"AnbariAPI/database"
	"log"
)

func main() {

	// Initialize database
	config := database.DefaultConfig()

	config.DBPath = "./data/inventory.db" // Custom path
	config.Debug = true                   // Enable SQL logging

	if err := database.InitDatabase(config); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.CloseDatabase()

	database.GetDB()

}
