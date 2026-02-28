package service

import (
	"embed"
	"fmt"
	"log"
	"sort"
	"strings"

	"gorm.io/gorm"
)

// MigrationService handles database schema migrations
type MigrationService struct {
	db *gorm.DB
}

// NewMigrationService creates a new migration service
func NewMigrationService(db *gorm.DB) *MigrationService {
	return &MigrationService{db: db}
}

// Migration represents a database migration record
type Migration struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"uniqueIndex;size:255;not null"`
	AppliedAt int64  `gorm:"autoCreateTime"`
}

// RunMigrations executes all pending SQL migrations from the embedded files
func (m *MigrationService) RunMigrations(migrations embed.FS) error {
	// Create migrations tracking table if it doesn't exist
	if err := m.db.AutoMigrate(&Migration{}); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Read all migration files from root directory
	entries, err := migrations.ReadDir(".")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Sort migration files by name (they should be numbered)
	var migrationFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			migrationFiles = append(migrationFiles, entry.Name())
		}
	}
	sort.Strings(migrationFiles)

	log.Printf("Found %d migration files", len(migrationFiles))

	// Execute each migration if not already applied
	for _, filename := range migrationFiles {
		if err := m.runMigration(migrations, filename); err != nil {
			return fmt.Errorf("failed to run migration %s: %w", filename, err)
		}
	}

	log.Println("All migrations completed successfully")
	return nil
}

// runMigration executes a single migration file if not already applied
func (m *MigrationService) runMigration(migrations embed.FS, filename string) error {
	// Check if migration was already applied
	var existing Migration
	result := m.db.Where("name = ?", filename).First(&existing)
	if result.Error == nil {
		log.Printf("Migration %s already applied, skipping", filename)
		return nil
	}

	// Read migration file from root directory
	content, err := migrations.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	sql := string(content)
	if strings.TrimSpace(sql) == "" {
		log.Printf("Migration %s is empty, skipping", filename)
		return nil
	}

	log.Printf("Applying migration: %s", filename)

	// Execute migration in a transaction
	err = m.db.Transaction(func(tx *gorm.DB) error {
		// Split SQL by semicolons and execute each statement
		statements := strings.Split(sql, ";")
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" || strings.HasPrefix(stmt, "--") {
				continue
			}

			if err := tx.Exec(stmt).Error; err != nil {
				// Ignore "Duplicate column" errors for ALTER TABLE ADD COLUMN
				if strings.Contains(err.Error(), "Duplicate column") ||
					strings.Contains(err.Error(), "duplicate key") ||
					strings.Contains(err.Error(), "already exists") {
					log.Printf("Warning: %v (continuing)", err)
					continue
				}
				return fmt.Errorf("failed to execute statement: %w", err)
			}
		}

		// Record migration as applied
		migration := Migration{Name: filename}
		if err := tx.Create(&migration).Error; err != nil {
			return fmt.Errorf("failed to record migration: %w", err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	log.Printf("Migration %s applied successfully", filename)
	return nil
}
