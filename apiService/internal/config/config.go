package config

import (
	"checkdown/apiService/internal/pkg/logger"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	GRPCAddr string `env:"GRPCAddr" env-required:"true"`
	HTTPPort int    `env:"HTTP_PORT" env-default:"8080"`

	KafkaAddr  string `env:"KAFKA_ADDR" env-required:"true"` // например "kafka:9092"
	KafkaTopic string `env:"KafkaTopic" env-required:"true"` // например "events"

	// параметры логов
	LogLevel string `env:"LOG_LEVEL" env-default:"info"`
	Env      string `env:"ENV" env-default:"local"`
}

// LoadConfig читает .env и инициализирует zap‑логер
func LoadConfig() *Config {
	cfg := &Config{}
	if err := cleanenv.ReadConfig("./.env", cfg); err != nil {
		logger.Init("error", "local")
		logger.Log.Fatalw("config read error", "err", err)
	}
	logger.Init(cfg.LogLevel, cfg.Env)
	return cfg
}
