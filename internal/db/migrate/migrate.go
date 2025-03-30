package migrate

import (
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// MigrateDB performs database migrations
func MigrateDB(dbPath string, migrationsPath string) error {
	migrationsDSN := fmt.Sprintf("file://%s", migrationsPath)
	dbDSN := fmt.Sprintf("sqlite3://%s", dbPath)

	m, err := migrate.New(migrationsDSN, dbDSN)
	if err != nil {
		return fmt.Errorf("error creating migrate instance: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("error running migrations: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return fmt.Errorf("error getting migration version: %w", err)
	}

	log.Printf("Migration completed. Version: %d, Dirty: %v", version, dirty)
	return nil
}

// RollbackDB rolls back the last migration
func RollbackDB(dbPath string, migrationsPath string) error {
	migrationsDSN := fmt.Sprintf("file://%s", migrationsPath)
	dbDSN := fmt.Sprintf("sqlite3://%s", dbPath)

	m, err := migrate.New(migrationsDSN, dbDSN)
	if err != nil {
		return fmt.Errorf("error creating migrate instance: %w", err)
	}
	defer m.Close()

	if err := m.Steps(-1); err != nil {
		return fmt.Errorf("error rolling back migration: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return fmt.Errorf("error getting migration version: %w", err)
	}

	log.Printf("Rollback completed. Version: %d, Dirty: %v", version, dirty)
	return nil
}

// GetMigrationVersion returns the current migration version
func GetMigrationVersion(dbPath string, migrationsPath string) (uint, bool, error) {
	migrationsDSN := fmt.Sprintf("file://%s", migrationsPath)
	dbDSN := fmt.Sprintf("sqlite3://%s", dbPath)

	m, err := migrate.New(migrationsDSN, dbDSN)
	if err != nil {
		return 0, false, fmt.Errorf("error creating migrate instance: %w", err)
	}
	defer m.Close()

	version, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return 0, false, fmt.Errorf("error getting migration version: %w", err)
	}

	return version, dirty, nil
}
