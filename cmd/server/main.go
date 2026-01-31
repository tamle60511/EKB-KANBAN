package main

import (
	"cqs-kanban/config"
	"cqs-kanban/database"
	"cqs-kanban/internal/app"
	"cqs-kanban/logger"
)

func main() {
	// Load configuration
	cfg := config.MustConfig()
	log := logger.NewLogger(cfg)
	// Connect to database
	db := database.MustNewDatabase(cfg, log)

	// Create application
	application := app.New(cfg, db)

	// Setup routes
	application.SetupRoutes()

	// Start application
	application.Start()
}
