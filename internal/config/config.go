package config

import "os"

type Config struct {
	Port               string
	JwtSecret          string
	JwtExpirationHours string
	DBPath             string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		Port:               os.Getenv("PORT"),
		JwtSecret:          os.Getenv("JWT_SECRET"),
		JwtExpirationHours: os.Getenv("JWT_EXPIRATION_HOURS"),
		DBPath:             os.Getenv("DB_PATH"),
	}

	if cfg.DBPath == "" {
		cfg.DBPath = "data/app.db"
	}

	return cfg, nil
}
