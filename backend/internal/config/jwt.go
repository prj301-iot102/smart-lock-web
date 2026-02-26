package config

import "time"

type Jwt struct {
	JwtSecret []byte        `env:"JWT_SECRET"`
	ExpiresIn time.Duration `env:"JWT_EXPIRES_IN"`
}
