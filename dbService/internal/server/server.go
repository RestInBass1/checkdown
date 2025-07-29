package server

import (
	"checkdown/dbService/internal/pkg/logger"
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	dto "checkdown/dbService/internal/DTO"
	"checkdown/dbService/internal/config"
	"checkdown/dbService/internal/transport/grpcHandlers"
	pb "checkdown/dbService/pkg/api"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"google.golang.org/grpc"
)

type PostgresRepository interface {
	CreateTask(ctx context.Context, task dto.Task) (int, error)
	GetTasks(ctx context.Context) ([]*dto.Task, error)
	UpdateTask(ctx context.Context, id int64) error
	DeleteTask(ctx context.Context, id int64) error
}

type Server struct {
	srv  *grpc.Server
	addr string
}

// NewServer собирает gRPC‑сервер, регистрирует имплементацию и возвращает обёртку.
func NewServer(ctx context.Context, cfg *config.Config, repo PostgresRepository) (*Server, error) {
	addr := fmt.Sprintf(":%d", cfg.GRPCPORT)

	// zap‑интерсептор для автоматического логирования RPC‑вызовов
	zapInter := grpc_zap.UnaryServerInterceptor(logger.Log.Desugar())
	s := grpc.NewServer(grpc.UnaryInterceptor(zapInter))

	// регистрируем сгенерированный gRPC‑сервис
	pb.RegisterDBServiceServer(s, grpcHandlers.NewDBService(repo))

	return &Server{srv: s, addr: addr}, nil
}

// Start запускает gRPC‑сервер (блокирующий вызов).
func (s *Server) Start() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("listen %s: %w", s.addr, err)
	}
	logger.Log.Infow("grpc listen", "addr", s.addr)

	if err := s.srv.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		logger.Log.Errorw("grpc serve error", "err", err)
		return err
	}
	logger.Log.Infow("grpc server stopped")
	return nil
}

// Stop завершает работу сервера корректно, с таймаутом на graceful shutdown.
func (s *Server) Stop() {
	logger.Log.Infow("grpc graceful stop")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() { // fallback: принудительная остановка
		<-ctx.Done()
		if ctx.Err() == context.DeadlineExceeded {
			logger.Log.Errorw("grpc stop timeout — forcing stop")
			s.srv.Stop()
		}
	}()

	s.srv.GracefulStop()
}
