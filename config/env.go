package config

import (
	"log"

	"github.com/joho/godotenv"
)

// Проверяем наличие файла окружения
func LoadEnv() {
	// загружаем значения из .env в систему
	if err := godotenv.Load("d:/LoginTracking/config/.env"); err != nil {
		log.Printf("no .env file found: %v", err)
	}
}
