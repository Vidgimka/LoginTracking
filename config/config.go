package config

type Config struct {
	DataBase DataBaseConfig
}

type DataBaseConfig struct {
	Host     string `envconfig:"DB_HOST" required:"true"`
	User     string `envconfig:"DB_USER" required:"true"`
	Password string `envconfig:"DB_PASSWORD" required:"true"`
	Name     string `envconfig:"DB_NAME" required:"true"`
	Port     int    `envconfig:"DB_PORT" required:"true"`
	SSLMode  string `envconfig:"DB_SSLMODE" required:"true"`
}
