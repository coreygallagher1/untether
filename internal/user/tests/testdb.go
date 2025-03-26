package tests

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	// Get test database configuration from environment or use defaults
	dbHost := getEnvOrDefault("TEST_DB_HOST", "localhost")
	dbPort := getEnvOrDefault("TEST_DB_PORT", "5432")
	dbUser := getEnvOrDefault("TEST_DB_USER", "untether")
	dbPassword := getEnvOrDefault("TEST_DB_PASSWORD", "untether")
	dbName := getEnvOrDefault("TEST_DB_NAME", "untether_test")

	// Create test database connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// Connect to the database
	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)

	// Test the connection
	err = db.Ping()
	require.NoError(t, err)

	// Run migrations
	err = runMigrations(db)
	require.NoError(t, err)

	// Register cleanup function
	t.Cleanup(func() {
		cleanupTestDB(t, db)
	})

	return db
}

func runMigrations(db *sql.DB) error {
	// Read migration files
	files, err := filepath.Glob("../../migrations/*.up.sql")
	if err != nil {
		return fmt.Errorf("failed to read migration files: %v", err)
	}

	// Sort files to ensure correct order
	sort.Strings(files)

	// Read and execute each migration
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %v", file, err)
		}

		// Execute migration
		_, err = db.Exec(string(content))
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %v", file, err)
		}
	}

	return nil
}

func cleanupTestDB(t *testing.T, db *sql.DB) {
	// Clean up test data
	_, err := db.Exec("TRUNCATE TABLE users CASCADE")
	require.NoError(t, err)

	// Close the connection
	err = db.Close()
	require.NoError(t, err)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 