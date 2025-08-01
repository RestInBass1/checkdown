package main

import (
	grpcadapter "checkdown/apiService/internal/adapter/grpc"
	"checkdown/apiService/internal/config"
	"checkdown/apiService/internal/kafka"
	"checkdown/apiService/internal/pkg/logger"
	"checkdown/apiService/internal/server"
	"checkdown/apiService/internal/transport/httpHandlers"
	"checkdown/dbService/pkg/api"
	"google.golang.org/grpc"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	cfg := config.LoadConfig()
	conn, err := grpc.Dial(cfg.GRPCAddr, grpc.WithInsecure())
	if err != nil {
		logger.Log.Fatalw("grpc dial failed", "addr", cfg.GRPCAddr, "err", err)
	}
	defer conn.Close()
	logger.Log.Infow("grpc connected", "addr", cfg.GRPCAddr)
	dbClient := api.NewDBServiceClient(conn)
	brokers := []string{cfg.KafkaAddr}
	prod, err := kafka.NewProducer(brokers, cfg.KafkaTopic)
	if err != nil {
		logger.Log.Fatalw("new producer failed", "err", err)
	}
	defer prod.Close()
	taskSvc := grpcadapter.New(dbClient, prod)
	h := httpHandlers.New(taskSvc)
	srv := server.New(":"+strconv.Itoa(cfg.HTTPPort), h.NewRouter())

	// запускаем server.Start() в отдельной горутине
	go func() {
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatalw("http server error", "err", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // ждём Ctrl‑C или SIGTERM от Kubernetes

	srv.Stop() // корректно гасим HTTP‑сервер

}
