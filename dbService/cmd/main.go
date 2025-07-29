package main

import (
	"checkdown/dbService/internal/pkg/logger"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"checkdown/dbService/internal/config"
	"checkdown/dbService/internal/repository"
	"checkdown/dbService/internal/server"
)

func main() {
	// ── конфиг + инициализация логгера ───────────────────────────────────
	cfg := config.LoadConfig() // внутри уже вызван logger.Init()
	defer logger.Log.Sync()    // flush перед выходом

	ctx := context.Background()

	// ── инициализируем хранилище ─────────────────────────────────────────
	storage, err := repository.NewStorage(ctx, cfg.POSTGRESCONFIG)
	if err != nil {
		logger.Log.Fatalw("storage init failed", "err", err)
	}

	// ── поднимаем gRPC‑сервер ────────────────────────────────────────────
	srv, err := server.NewServer(ctx, cfg, storage)
	if err != nil {
		logger.Log.Fatalw("server init failed", "err", err)
	}

	go func() {
		if err := srv.Start(); err != nil {
			logger.Log.Fatalw("grpc server error", "err", err)
		}
	}()

	addr := fmt.Sprintf(":%d", cfg.GRPCPORT)
	logger.Log.Infow("grpc server started", "addr", addr)

	// ── ловим SIGINT / SIGTERM ───────────────────────────────────────────
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	logger.Log.Infow("shutting down grpc server")
	srv.Stop()
}
