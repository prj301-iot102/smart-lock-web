package config

type Jwt struct {
	JwtSecret int `env:"JWT_SECRET"`
	ExpiresIn int `env:"JWT_EXPIRES_IN"`
}
