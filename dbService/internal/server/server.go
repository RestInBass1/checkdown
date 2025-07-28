package server

import (
	"checkdown/dbService/internal/DTO"
	"checkdown/dbService/internal/config"
	"checkdown/dbService/internal/transport/grpcHandlers"
	"checkdown/dbService/pkg/api"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
)

type PostgresRepository interface {
	CreateTask(ctx context.Context, task dto.Task) (int, error)
	GetTasks(ctx context.Context) ([]*dto.Task, error)
	UpdateTask(ctx context.Context, id int64) error
	DeleteTask(ctx context.Context, id int64) error
}
type Server struct {
	grpcServer   *grpc.Server
	grpcListener net.Listener
}

func NewServer(ctx context.Context, cfg *config.Config, repo PostgresRepository) (*Server, error) {
	gRPCaddr := fmt.Sprintf(":%d", cfg.GRPCPORT)
	grpcListener, err := net.Listen("tcp", gRPCaddr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on %s: %w", gRPCaddr, err)
	}
	opts := []grpc.ServerOption{}
	grpcSrv := grpc.NewServer(opts...)
	dbservice := grpcHandlers.NewDBService(repo)
	api.RegisterDBServiceServer(grpcSrv, dbservice)
	return &Server{
		grpcServer:   grpcSrv,
		grpcListener: grpcListener,
	}, nil
}

func (s *Server) Start() error {
	return s.grpcServer.Serve(s.grpcListener)
}

func (s *Server) Stop() {
	s.grpcServer.GracefulStop()
}
