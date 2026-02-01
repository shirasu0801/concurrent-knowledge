package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"concurrent-knowledge/internal/config"
	"concurrent-knowledge/internal/database"
)

func main() {
	// Parse command line flags
	force := flag.Bool("force", false, "Force re-initialization (deletes existing database)")
	flag.Parse()

	log.Println("========================================")
	log.Println("Concurrent Knowledge - DB Initialization")
	log.Println("========================================")
	log.Println()

	// Load configuration
	cfg := config.Load()

	// Check if database file exists
	if _, err := os.Stat(cfg.DatabasePath); err == nil {
		if !*force {
			log.Printf("Database already exists at: %s", cfg.DatabasePath)
			log.Println("Use -force flag to re-initialize (WARNING: This will delete all data)")
			os.Exit(1)
		}

		// Delete existing database
		log.Printf("Deleting existing database: %s", cfg.DatabasePath)
		if err := os.Remove(cfg.DatabasePath); err != nil {
			log.Fatalf("Failed to delete existing database: %v", err)
		}
	}

	// Ensure data directory exists
	dataDir := "./data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	// Initialize database
	log.Printf("Creating database at: %s", cfg.DatabasePath)
	db, err := database.InitDB(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Run migrations
	migrationsDir := "./migrations"
	if err := database.RunMigrations(db, migrationsDir); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Verify database
	log.Println()
	log.Println("Verifying database...")

	genreCount, err := database.GetGenreCount(db)
	if err != nil {
		log.Fatalf("Failed to count genres: %v", err)
	}
	log.Printf("Genres: %d", genreCount)

	questionCount, err := database.GetQuestionCount(db)
	if err != nil {
		log.Fatalf("Failed to count questions: %v", err)
	}
	log.Printf("Questions: %d", questionCount)

	log.Println()
	log.Println("========================================")
	log.Println("Database initialization completed!")
	log.Println("========================================")
	fmt.Printf("Database location: %s\n", cfg.DatabasePath)
	fmt.Printf("Total genres: %d\n", genreCount)
	fmt.Printf("Total questions: %d\n", questionCount)
	log.Println()
	log.Println("You can now start the server with:")
	log.Println("  go run cmd/server/main.go")
	log.Println()
}
