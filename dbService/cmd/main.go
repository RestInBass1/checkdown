package main

import (
	"checkdown/dbService/internal/config"
	"checkdown/dbService/internal/repository"
	"checkdown/dbService/internal/server"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx := context.Background()
	cfg := config.LoadConfig()
	storage, err := repository.NewStorage(ctx, cfg.POSTGRESCONFIG)
	if err != nil {
		log.Fatal(err)
	}
	server, err := server.NewServer(ctx, cfg, storage)
	if err != nil {
		log.Fatal(err)
	}

	doneChan := make(chan os.Signal)
	signal.Notify(doneChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err = server.Start()
		if err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	}()
	<-doneChan
	server.Stop()

}
