package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Vidgimka/LoginTracking/internal/config"
	"github.com/Vidgimka/LoginTracking/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kelseyhightower/envconfig"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// функция подключения к БД  т
func NewDatabase(cfg *config.DataBaseConfig) (*gorm.DB, error) {
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("failed to parsing: %v", err)
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s  dbname=%s  port=%d  sslmode=%s", cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port, cfg.SSLMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	err = db.AutoMigrate(&domain.Data{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}
	return db, nil
}

func SetPool(ctx context.Context, cfg *config.DataBaseConfig) (*pgxpool.Pool, error) {
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("failed to parsing: %v", err)
	}
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
