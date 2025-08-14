package infrastructure

import (
	"fmt"
	"log"

	"github.com/Vidgimka/LoginTracking/internal/config"
	"github.com/Vidgimka/LoginTracking/internal/domain"
	"github.com/kelseyhightower/envconfig"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// функция подключения к БД  т
func Init() (*gorm.DB, error) {
	var cfg config.Config

	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("failed to parsing: %v", err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s  dbname=%s  port=%d  sslmode=%s", cfg.DataBase.Host, cfg.DataBase.User, cfg.DataBase.Password, cfg.DataBase.Name, cfg.DataBase.Port, cfg.DataBase.SSLMode)
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
