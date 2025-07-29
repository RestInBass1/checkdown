package config

import (
	"checkdown/dbService/internal/pkg/logger"
	"checkdown/dbService/internal/repository"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config агрегирует все переменные окружения dbService.
type Config struct {
	// PostgreSQL
	POSTGRESCONFIG repository.Config

	// gRPC
	GRPCPORT int `env:"GRPC_PORT" env-default:"50051"`

	// логирование
	LogLevel string `env:"LOG_LEVEL" env-default:"info"`
	Env      string `env:"ENV" env-default:"local"`
}

// LoadConfig читает .env / OS env и инициализирует zap‑логер.
func LoadConfig() *Config {
	var cfg Config

	// пытаемся прочитать .env
	if err := cleanenv.ReadConfig("./.env", &cfg); err != nil {
		// сделаем «резервный» логер, чтобы вывести ошибку
		logger.Init("error", "local")
		logger.Log.Fatalw("config load error", "err", err)
	}

	// полноценный логер
	logger.Init(cfg.LogLevel, cfg.Env)

	logger.Log.Infow("config loaded",
		"grpc_port", cfg.GRPCPORT,
		"pg_host", cfg.POSTGRESCONFIG.Host,
	)
	return &cfg
}
