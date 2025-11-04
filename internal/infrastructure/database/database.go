package database

import (
	"context"
	"fmt"
	"time"

	"github.com/Vidgimka/LoginTracking/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetPool(ctx context.Context, cfg *config.DataBaseConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s  dbname=%s  port=%d  sslmode=%s", cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port, cfg.SSLMode)
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse config failed: %w", err)
	}

	config.MaxConns = 20
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("create pool failed: %w", err)
	}
	return pool, nil
}
