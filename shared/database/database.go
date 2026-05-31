package database

import (
	"AnbariAPI/shared/migration"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Config holds database configuration
type Config struct {
	DBPath string
	Debug  bool
}

// DefaultConfig returns default configuration
func DefaultConfig() Config {
	return Config{
		DBPath: "inventory.db",
		Debug:  true, // Set to false in production
	}
}

// InitDatabase initializes the database connection and runs migrations
func InitDatabase(config Config) error {
	// Ensure the directory for the database file exists
	dbDir := filepath.Dir(config.DBPath)
	if dbDir != "." {
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return fmt.Errorf("failed to create database directory: %w", err)
		}
	}

	// Configure GORM logger
	var logLevel logger.LogLevel
	if config.Debug {
		logLevel = logger.Info
	} else {
		logLevel = logger.Error
	}

	// SQLite specific configuration
	db, err := gorm.Open(sqlite.Open(config.DBPath), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
		// DisableForeignKeyConstraintWhenMigrating: true, // Uncomment if needed for SQLite
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Enable WAL mode for better concurrent access (SQLite specific)
	if err := db.Exec("PRAGMA journal_mode=WAL").Error; err != nil {
		log.Printf("Warning: Could not enable WAL mode: %v", err)
	}

	// Enable foreign keys support (important for SQLite!)
	if err := db.Exec("PRAGMA foreign_keys=ON").Error; err != nil {
		log.Printf("Warning: Could not enable foreign keys: %v", err)
	}

	// Optimize SQLite settings
	pragmas := []string{
		"PRAGMA busy_timeout=5000",
		"PRAGMA synchronous=NORMAL",
		"PRAGMA cache_size=-2000", // 2MB cache
		"PRAGMA temp_store=MEMORY",
		"PRAGMA mmap_size=268435456", // 256MB memory map
	}

	for _, pragma := range pragmas {
		if err := db.Exec(pragma).Error; err != nil {
			log.Printf("Warning: Could not set %s: %v", pragma, err)
		}
	}

	// Assign to global variable
	DB = db

	// Run migrations
	if err := migration.Migrate(DB); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database initialized successfully")
	return nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}

// CloseDatabase closes the database connection
func CloseDatabase() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
