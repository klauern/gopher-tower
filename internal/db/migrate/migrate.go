package migrate

import (
	"errors"
	"fmt"
	"io/fs"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

// FS is an interface that both embed.FS and testing filesystems can implement
type FS interface {
	fs.FS
}

// MigrateDB performs database migrations
func MigrateDB(dbPath string, migrations FS) error {
	// We use "." as the root of the migrations
	d, err := iofs.New(migrations, ".")
	if err != nil {
		return fmt.Errorf("error creating migration source: %w", err)
	}

	dbDSN := fmt.Sprintf("sqlite3://%s", dbPath)
	m, err := migrate.NewWithSourceInstance("iofs", d, dbDSN)
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
func RollbackDB(dbPath string, migrations FS) error {
	d, err := iofs.New(migrations, ".")
	if err != nil {
		return fmt.Errorf("error creating migration source: %w", err)
	}

	dbDSN := fmt.Sprintf("sqlite3://%s", dbPath)
	m, err := migrate.NewWithSourceInstance("iofs", d, dbDSN)
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
func GetMigrationVersion(dbPath string, migrations FS) (uint, bool, error) {
	d, err := iofs.New(migrations, ".")
	if err != nil {
		return 0, false, fmt.Errorf("error creating migration source: %w", err)
	}

	dbDSN := fmt.Sprintf("sqlite3://%s", dbPath)
	m, err := migrate.NewWithSourceInstance("iofs", d, dbDSN)
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
