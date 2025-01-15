package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/cyberkillua/dailyread/internal/config"
	"github.com/cyberkillua/dailyread/internal/database"
	"github.com/cyberkillua/dailyread/internal/server"
	"github.com/cyberkillua/dailyread/internal/utils"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {

	godotenv.Load()
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Verify critical environment variables
	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	// Database connection
	connection, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize database queries
	db := database.New(connection)

	go utils.StartScrapping(db, 10, 12*time.Hour)

	// Create and start server
	srv := server.New(cfg, db)
	if err := srv.Start(); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
