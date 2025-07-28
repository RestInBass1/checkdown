package main

import (
	grpcadapter "checkdown/apiService/internal/adapter/grpc"
	"checkdown/apiService/internal/config"
	"checkdown/apiService/internal/server"
	"checkdown/apiService/internal/transport/httpHandlers"
	"checkdown/dbService/pkg/api"
	"google.golang.org/grpc"
	"log"
)

func main() {
	cfg := config.LoadConfig()
	conn, err := grpc.Dial(
		cfg.GRPCAddr)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	dbClient := api.NewDBServiceClient(conn)
	taskSvc := grpcadapter.New(dbClient)
	h := httpHandlers.New(taskSvc)
	server := server.New(cfg.GRPCAddr, h.NewRouter())
	server.Start()
}
