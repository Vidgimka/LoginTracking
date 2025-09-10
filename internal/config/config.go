package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	DataBase DataBaseConfig `yaml:"data_base"`
	Client   Client         `yaml:"client"`
}

type DataBaseConfig struct {
	Host     string `yaml:"host" envconfig:"DB_HOST" required:"true"`
	User     string `yaml:"user" envconfig:"DB_USER" required:"true"`
	Password string `            envconfig:"DB_PASSWORD" required:"true"`
	Name     string `            envconfig:"DB_NAME" required:"true"`
	Port     int    `yaml:"port" envconfig:"DB_PORT" required:"true"`
	SSLMode  string `yaml:"ssl_mode" envconfig:"DB_SSLMODE" required:"true"`
}

type Client struct {
	ClientTimeOut int    `yaml:"client_timeout" envconfig:"CLIENT_TIMEOUT" env-default:"10"`
	BaseUrl       string `yaml:"base_url"           envconfig:"BASE_URL" required:"true"`
}

func NewConfig(envPath, yamlPath string) (*Config, error) {
	cfg := &Config{}

	if envPath == "" {
		envPath = ".env"
	}
	if err := godotenv.Load(envPath); err != nil {
		fmt.Errorf("no .env file found: %w", err)
	}

	if yamlPath == "" {
		yamlPath = "./config/config.yml"
	}
	if err := cleanenv.ReadConfig(yamlPath, cfg); err != nil {
		fmt.Errorf("yaml config file  not found: %w", err)
	}
	return cfg, nil
}
