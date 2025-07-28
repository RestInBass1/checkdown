package config

import (
	"checkdown/dbService/internal/repository"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	POSTGRESCONFIG repository.Config
	GRPCPORT       int `env:"GRPC_PORT" env-required:"true"`
}

func LoadConfig() *Config {
	cfg := &Config{}
	err := cleanenv.ReadConfig("./.env", cfg)
	if err != nil {
		log.Fatalf("error reading config: %s", err.Error())
	}
	return cfg
}
