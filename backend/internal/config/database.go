package config

import "fmt"

type Database struct {
	Database         string `env:"DB_DATABASE"`
	DatabaseUser     string `env:"DB_USERNAME"`
	DatabasePassword string `env:"DB_PASSWORD"`
	DatabasePort     int    `env:"DB_PORT"`
	DatabaseHost     string `env:"DB_HOST"`
}

func (cfg *Database) DatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.DatabaseUser, cfg.DatabasePassword, cfg.DatabaseHost, cfg.DatabasePort, cfg.Database)
}
