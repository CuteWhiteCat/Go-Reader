package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

// NewTestDatabase creates a new in-memory SQLite database for testing and runs migrations.
func NewTestDatabase(t *testing.T) *sqlx.DB {
	t.Helper()

	// Use in-memory SQLite database for tests
	db, err := sqlx.Connect("sqlite3", ":memory:")
	require.NoError(t, err, "Failed to connect to in-memory database")

	// Find the project root to locate the migrations directory
	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok, "Failed to get current file path")
	projectRoot := filepath.Join(filepath.Dir(currentFile), "..", "..") // Navigate up from config/internal/backend
	migrationsDir := filepath.Join(projectRoot, "migrations")

	// Read and execute migration files
	migrations, err := os.ReadDir(migrationsDir)
	require.NoError(t, err, "Failed to read migrations directory at %s", migrationsDir)

	for _, migrationFile := range migrations {
		if filepath.Ext(migrationFile.Name()) == ".sql" {
			migrationPath := filepath.Join(migrationsDir, migrationFile.Name())
			content, err := os.ReadFile(migrationPath)
			require.NoError(t, err, "Failed to read migration file: %s", migrationFile.Name())

			_, err = db.Exec(string(content))
			require.NoError(t, err, "Failed to execute migration: %s", migrationFile.Name())
		}
	}

	// Add a cleanup function to close the database connection
	t.Cleanup(func() {
		err := db.Close()
		require.NoError(t, err, "Failed to close database connection")
	})

	return db
}
