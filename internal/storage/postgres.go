// Package storage responsible for connecting with the database and send connection object to the ragger.
package storage

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // ...
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	ErrDBInteracting     = errors.New("error interacting with the database")
	ErrNotAvailableICCID = errors.New("this iccid not available")
)

// Config stores configs for postgres storage.
type Config struct {
	Host    string `config:"DATABASE_HOST,required"`
	Name    string `config:"DATABASE_NAME,required"`
	User    string `config:"DATABASE_USER,required"`
	Pass    string `config:"DATABASE_PASSWORD,required"`
	SSLMode string `config:"DATABASE_SSLMODE,required"`
}

// Postgres returns postgres dsn.
func (c Config) Postgres() string {
	return fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=%s",
		c.Host, c.Name, c.User, c.Pass, c.SSLMode)
}

// Connect return database connection and error if postgres connection don't open
func Connect(dsn string, logger *logrus.Logger) (*sqlx.DB, error) {
	var err error

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		logger.Fatal("opening db connection:", err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		logger.Fatal("pinging db connection:", err)
		return nil, err
	}

	return db, nil
}

// MakeMigrations provides an opportunity to work with migrations
func (c Config) MakeMigrations(logger *logrus.Logger) error {
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		c.User, c.Pass, c.Host, c.Name, c.SSLMode)

	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		logger.Fatal("failed to creating new migrations:", err)
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Fatal("failed to upping migrations:", err)
		return err
	}

	return nil
}
