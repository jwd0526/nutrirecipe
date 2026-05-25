package db

import (
	"context"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	pgkitx "github.com/jwd0526/pgkitx"
)

func Connect(ctx context.Context) (*pgkitx.Pool, error) {
	cfg, err := pgkitx.LoadConfigFromEnv()
	if err != nil {
		return nil, fmt.Errorf("db config: %w", err)
	}

	pool, err := pgkitx.NewPool(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("db pool: %w", err)
	}

	return pool, nil
}

func Migrate(databaseURL string) error {
	m, err := migrate.New("file://db/migrations", databaseURL)
	if err != nil {
		return fmt.Errorf("migrate init: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migrate up: %w", err)
	}

	return nil
}
