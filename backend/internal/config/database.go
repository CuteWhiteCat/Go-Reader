package config

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

// InitDatabase initializes the database connection and runs migrations
func InitDatabase(dbPath string) (*sqlx.DB, error) {
	// Create database directory if it doesn't exist
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open database connection
	db, err := sqlx.Connect("sqlite3", dbPath+"?_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Run migrations
	if err := runMigrations(db.DB); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	logrus.Info("Database initialized successfully")
	return db, nil
}

// runMigrations runs database migrations
func runMigrations(db *sql.DB) error {
	migrationFiles := []string{
		"001_initial.sql",
		"002_add_volume_columns.sql",
	}

	pathsToTry := []string{
		"./backend/migrations",
		"./migrations",
	}

	for _, file := range migrationFiles {
		if file == "002_add_volume_columns.sql" {
			hasVolumeCols, err := columnExists(db, "chapters", "volume_number")
			if err != nil {
				return fmt.Errorf("failed to inspect schema before %s: %w", file, err)
			}
			if hasVolumeCols {
				logrus.Infof("Skipping migration %s (columns already exist)", file)
				continue
			}
		}

		var migrationSQL []byte
		var readErr error

		for _, base := range pathsToTry {
			migrationPath := filepath.Join(base, file)
			migrationSQL, readErr = os.ReadFile(migrationPath)
			if readErr == nil {
				break
			}
		}

		if readErr != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, readErr)
		}

		if _, err := db.Exec(string(migrationSQL)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file, err)
		}
		logrus.Infof("Migration %s executed successfully", file)
	}

	return nil
}

func columnExists(db *sql.DB, tableName, columnName string) (bool, error) {
	rows, err := db.Query(fmt.Sprintf("PRAGMA table_info(%s)", tableName))
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		var cid int
		var name string
		var ctype string
		var notnull int
		var dfltValue interface{}
		var pk int

		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dfltValue, &pk); err != nil {
			return false, err
		}
		if strings.EqualFold(name, columnName) {
			return true, nil
		}
	}

	return false, rows.Err()
}
