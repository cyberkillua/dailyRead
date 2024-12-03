package main

import (
	"database/sql"
	"log"

	"github.com/cyberkillua/dailyread/internal/config"
	"github.com/cyberkillua/dailyread/internal/database"
	"github.com/cyberkillua/dailyread/internal/server"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Database connection
	connection, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize database queries
	db := database.New(connection)

	// Create and start server
	srv := server.New(cfg, db)
	if err := srv.Start(); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
