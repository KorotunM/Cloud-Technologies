package database

import (
	"context"
	"fmt"
	"log"

	"pragma/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB is a shared connection pool.
var DB *pgxpool.Pool

// ConnectDB establishes a connection pool to PostgreSQL using env-driven config.
func ConnectDB(ctx context.Context, cfg config.DatabaseConfig) error {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.SSLMode)

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("parse postgres config: %w", err)
	}

	poolCfg.MaxConns = cfg.MaxConns
	poolCfg.MinConns = cfg.MinConns
	poolCfg.MaxConnLifetime = cfg.MaxConnLifetime
	poolCfg.MaxConnIdleTime = cfg.MaxConnIdleTime

	DB, err = pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return fmt.Errorf("init postgres pool: %w", err)
	}

	if err = DB.Ping(ctx); err != nil {
		return fmt.Errorf("ping postgres: %w", err)
	}

	log.Println("Connected to PostgreSQL.")
	return nil
}

// CloseDB releases the PostgreSQL connection pool.
func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("PostgreSQL pool closed.")
	}
}
