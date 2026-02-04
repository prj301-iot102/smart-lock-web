package config

type Server struct {
	Port int `env:"PORT" envDefault:"8080"`
}
