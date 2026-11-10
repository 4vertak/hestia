package database

import (
	"embed"
	"fmt"
	"net"
	"net/url"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/4vertak/hestia/internal/config"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func NewPostgresDB(cfg *config.Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", postgresDSN(cfg))
	if err != nil {
		return nil, fmt.Errorf("failed to open db connection: %w", err)
	}

	// Настройка пула соединений
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	return db, nil
}

func postgresDSN(cfg *config.Config) string {
	dsn := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.DBUser, cfg.DBPassword),
		Host:     net.JoinHostPort(cfg.DBHost, cfg.DBPort),
		Path:     cfg.DBName,
		RawQuery: url.Values{"sslmode": {cfg.DBSSLMode}}.Encode(),
	}
	return dsn.String()
}

func newMigrator(db *sqlx.DB) (*migrate.Migrate, error) {
	source, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		return nil, err
	}
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return nil, err
	}
	return migrate.NewWithInstance("iofs", source, "postgres", driver)
}

// RunMigrations applies embedded migrations, including on application startup.
func RunMigrations(db *sqlx.DB) error {
	m, err := newMigrator(db)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

// RollbackMigrations reverts all migrations when explicitly requested by the CLI.
func RollbackMigrations(db *sqlx.DB) error {
	m, err := newMigrator(db)
	if err != nil {
		return err
	}
	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
