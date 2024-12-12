package migrator

import (
	"dns-server/internal/storage"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

func RunMigrations(db *storage.SQLite, migrationsPath string) error {
	driver, err := sqlite3.WithInstance(db.Conn, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("failed to init sqlite3 driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"sqlite3",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to init migrator: %w", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to applied migrations: %w", err)
	}

	log.Println("Migrations applied successfully")
	return nil
}
