package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"todo-lits-DMARK/app/internal/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type Database struct {
	DB *sql.DB
}

func New(cfg *config.Config) (*Database, error) {
	db, err := sql.Open("postgres", cfg.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	database := &Database{DB: db}

	if err := database.RunMigrations(); err != nil {
		log.Printf("Warning: failed to run migrations: %v", err)
	}

	return database, nil
}

func (d *Database) RunMigrations() error {
	driver, err := postgres.WithInstance(d.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create postgres driver: %w", err)
	}

	migrationPath := os.Getenv("MIGRATION_PATH")
	if migrationPath == "" {

		if wd, err := os.Getwd(); err == nil {

			appMigrations := filepath.Join(wd, "app", "migrations")
			if _, err := os.Stat(appMigrations); err == nil {
				migrationPath = "file://app/migrations"
			} else {

				migrationPath = "file://migrations"
			}
		} else {
			migrationPath = "file://migrations"
		}
	}

	log.Printf("Using migration path: %s", migrationPath)

	m, err := migrate.NewWithDatabaseInstance(
		migrationPath,
		"postgres", driver,
	)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not run migrations: %w", err)
	}

	log.Println("Migrations completed successfully")
	return nil
}

func (d *Database) Close() error {
	return d.DB.Close()
}

func (d *Database) Health() error {
	return d.DB.Ping()
}
