package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Storage  StorageConfig  `mapstructure:"storage"`
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Path string `mapstructure:"path"`
}

// StorageConfig holds storage-related configuration
type StorageConfig struct {
	BooksDir  string `mapstructure:"books_dir"`
	CoversDir string `mapstructure:"covers_dir"`
}

// Load loads configuration from file and environment variables
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("$HOME/.go-reader")

	// Set defaults
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("database.path", "./data/database.db")
	viper.SetDefault("storage.books_dir", "./data/books")
	viper.SetDefault("storage.covers_dir", "./data/covers")

	// Allow overriding with environment variables
	viper.SetEnvPrefix("GOREADER")
	viper.AutomaticEnv()

	// Read config file (ignore error if file doesn't exist)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Create necessary directories
	if err := createDirectories(&config); err != nil {
		return nil, fmt.Errorf("failed to create directories: %w", err)
	}

	return &config, nil
}

// createDirectories creates all necessary directories for the application
func createDirectories(config *Config) error {
	dirs := []string{
		filepath.Dir(config.Database.Path),
		config.Storage.BooksDir,
		config.Storage.CoversDir,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}
