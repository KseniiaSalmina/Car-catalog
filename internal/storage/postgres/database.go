package postgres

import (
	"context"
	"embed"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"github.com/KseniiaSalmina/Car-catalog/internal/config"
)

type DB struct {
	db *pgxpool.Pool
}

func NewDB(ctx context.Context, cfg config.Postgres) (*DB, error) {
	connstr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	db, err := pgxpool.New(ctx, connstr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err = db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	database := &DB{db: db}

	/*if err := database.migrate(cfg.Migration); err != nil {
		return nil, fmt.Errorf("failed to setup database: %w", err)
	}*/

	return database, nil
}

func (db *DB) migrate(embedMigration embed.FS) error {
	goose.SetBaseFS(embedMigration)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to migrate: %w", err)
	}

	d := stdlib.OpenDBFromPool(db.db)
	if err := goose.Up(d, "schema"); err != nil {
		return fmt.Errorf("failed to migrate: %w", err)
	}

	return nil
}

func (db *DB) Close() {
	db.db.Close()
}

func (db *DB) NewTransaction(ctx context.Context) (*Transaction, error) {
	tx, err := db.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	return &Transaction{tx: tx}, nil
}
