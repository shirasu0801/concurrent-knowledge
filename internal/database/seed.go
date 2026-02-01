package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// RunMigrations runs all SQL migration files in order
func RunMigrations(db *sql.DB, migrationsDir string) error {
	log.Println("Running database migrations...")

	// Read migration files
	migrations := []string{
		"001_schema.sql",
		"002_seed.sql",
	}

	for _, migration := range migrations {
		migrationPath := filepath.Join(migrationsDir, migration)
		log.Printf("Executing migration: %s", migration)

		// Read migration file
		sqlBytes, err := os.ReadFile(migrationPath)
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", migration, err)
		}

		// Execute migration
		_, err = db.Exec(string(sqlBytes))
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", migration, err)
		}

		log.Printf("Migration %s completed successfully", migration)
	}

	log.Println("All migrations completed successfully")
	return nil
}

// CheckDatabase checks if the database is initialized
func CheckDatabase(db *sql.DB) (bool, error) {
	// Check if genres table exists
	query := `SELECT name FROM sqlite_master WHERE type='table' AND name='genres';`
	var tableName string
	err := db.QueryRow(query).Scan(&tableName)

	if err == sql.ErrNoRows {
		return false, nil // Table doesn't exist
	}
	if err != nil {
		return false, err
	}

	return true, nil // Table exists
}

// GetQuestionCount returns the number of questions in the database
func GetQuestionCount(db *sql.DB) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM questions").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetGenreCount returns the number of genres in the database
func GetGenreCount(db *sql.DB) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM genres").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
