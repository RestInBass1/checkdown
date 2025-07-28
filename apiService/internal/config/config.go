package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	GRPCAddr string `env:"GRPCAddr" env-required:"true"`
	HTTPPort int    `env:"HTTP_PORT" env-default:"8080"`
}

func LoadConfig() *Config {
	cfg := &Config{}
	err := cleanenv.ReadConfig("./.env", cfg)
	if err != nil {
		log.Fatalf("error reading config: %s", err.Error())
	}
	return cfg
}
