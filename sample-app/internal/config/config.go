package config

// AppConfig holds the application configuration.
type AppConfig struct {
	Port        string `env:"PORT" default:"3000"`
	Environment string `env:"APP_ENV" default:"development"`
	DBHost      string `env:"DB_HOST" default:"localhost"`
	DBPort      int    `env:"DB_PORT" default:"5432"`
	DBUser      string `env:"DB_USER" default:"postgres"`
	DBPassword  string `env:"DB_PASSWORD"`
	DBName      string `env:"DB_NAME" default:"sample-app"`
	JWTSecret   string `env:"JWT_SECRET" default:"change-me-in-production"`
}
