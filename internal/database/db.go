package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// DB is the database connection
var DB *sql.DB

// InitDB initializes the database connection with optimizations
func InitDB(dbPath string) (*sql.DB, error) {
	log.Printf("Initializing database connection: %s", dbPath)

	// Open database connection
	// Parameters:
	// - cache=shared: Allow multiple connections to share cache
	// - mode=rwc: Read-write-create mode
	// - _journal_mode=WAL: Write-Ahead Logging for better concurrency
	db, err := sql.Open("sqlite3", fmt.Sprintf("%s?cache=shared&mode=rwc&_journal_mode=WAL", dbPath))
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool for concurrent requests
	db.SetMaxOpenConns(25)                 // Maximum number of open connections
	db.SetMaxIdleConns(5)                  // Maximum number of idle connections
	db.SetConnMaxLifetime(5 * time.Minute) // Maximum lifetime of a connection

	// Verify database connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Enable WAL mode for better concurrent read performance
	// This allows multiple readers and one writer at the same time
	if _, err := db.Exec("PRAGMA journal_mode=WAL;"); err != nil {
		log.Printf("Warning: Failed to set WAL mode: %v", err)
	}

	// Set busy timeout to 5 seconds to handle concurrent writes
	if _, err := db.Exec("PRAGMA busy_timeout=5000;"); err != nil {
		log.Printf("Warning: Failed to set busy timeout: %v", err)
	}

	// Enable foreign key constraints
	if _, err := db.Exec("PRAGMA foreign_keys=ON;"); err != nil {
		log.Printf("Warning: Failed to enable foreign keys: %v", err)
	}

	log.Printf("Database connection established successfully")

	DB = db
	return db, nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		log.Println("Closing database connection")
		return DB.Close()
	}
	return nil
}

// GetDB returns the database connection
func GetDB() *sql.DB {
	return DB
}
