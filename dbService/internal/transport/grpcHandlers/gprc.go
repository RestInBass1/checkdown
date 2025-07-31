package grpcHandlers

import (
	"checkdown/dbService/internal/pkg/logger"
	"context"
	"github.com/redis/go-redis/v9"
	"time"

	dto "checkdown/dbService/internal/DTO"
	api2 "checkdown/dbService/pkg/api"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PostgresRepository interface {
	CreateTask(ctx context.Context, task dto.Task) (int, error)
	GetTasks(ctx context.Context) ([]*dto.Task, error)
	UpdateTask(ctx context.Context, id int64) error
	DeleteTask(ctx context.Context, id int64) error
}

type RedisRepository interface {
	GetTasks(ctx context.Context) ([]*dto.Task, error)
	SetTasks(ctx context.Context, tasks []*dto.Task) error
	DeleteTasks(ctx context.Context) error
}

type DBService struct {
	api2.UnimplementedDBServiceServer
	repo      PostgresRepository
	repoRedis RedisRepository
}

func NewDBService(repo PostgresRepository, repoRed RedisRepository) *DBService {
	return &DBService{
		repo:      repo,
		repoRedis: repoRed,
	}
}

func (s *DBService) AddTask(ctx context.Context, req *api2.TaskRequest) (*api2.CreateTaskResponse, error) {
	start := time.Now()
	logger.Log.Debugw("RPC AddTask start",
		"title", req.Title,
		"description", req.Description,
	)

	id, err := s.repo.CreateTask(ctx, dto.Task{
		Title:       req.Title,
		Description: req.Description,
	})
	if err != nil {
		logger.Log.Errorw("RPC AddTask error",
			"error", err,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return nil, err
	}
	_ = s.repoRedis.DeleteTasks(ctx)
	logger.Log.Infow("RPC AddTask success",
		"id", id,
		"duration_ms", time.Since(start).Milliseconds(),
	)
	return &api2.CreateTaskResponse{Id: int64(id)}, nil
}

func (s *DBService) GetTasks(ctx context.Context, _ *emptypb.Empty) (*api2.GetTasksResponse, error) {
	start := time.Now()
	logger.Log.Debugw("RPC GetTasks start")

	tasks, err := s.repoRedis.GetTasks(ctx)
	if err == redis.Nil {
		logger.Log.Debugw("RPC GetTasks redis error")
		bro, err := s.repo.GetTasks(ctx)
		if err != nil {
			logger.Log.Errorw("RPC GetTasks postgres error",
				"error", err,
				"duration_ms", time.Since(start).Milliseconds(),
			)
			return nil, err
		}
		logger.Log.Infow("RPC GetTasks success")
		err = s.repoRedis.SetTasks(ctx, bro)
		if err != nil {
			logger.Log.Debugw("RPC GetTasks redis error")
			return nil, err
		}
		logger.Log.Infow("RPC GetTasks  redis set success")
		resp := &api2.GetTasksResponse{Tasks: make([]*api2.Task, 0, len(bro))}
		for _, t := range bro {
			if t == nil {
				continue
			}
			resp.Tasks = append(resp.Tasks, &api2.Task{
				Id:          int64(t.ID),
				Title:       t.Title,
				Description: t.Description,
				// конвертация bool→string, как в твоём proto
				IsDone:    t.IsDone,
				CreatedAt: timestamppb.New(t.CreatedAt),
				UpdatedAt: timestamppb.New(t.UpdatedAt),
			})
		}

		logger.Log.Infow("RPC GetTasks success",
			"count", len(resp.Tasks),
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return resp, nil
	}

	resp := &api2.GetTasksResponse{Tasks: make([]*api2.Task, 0, len(tasks))}
	for _, t := range tasks {
		if t == nil {
			continue
		}
		resp.Tasks = append(resp.Tasks, &api2.Task{
			Id:          int64(t.ID),
			Title:       t.Title,
			Description: t.Description,
			// конвертация bool→string, как в твоём proto
			IsDone:    t.IsDone,
			CreatedAt: timestamppb.New(t.CreatedAt),
			UpdatedAt: timestamppb.New(t.UpdatedAt),
		})
	}

	logger.Log.Infow("RPC GetTasks success",
		"count", len(resp.Tasks),
		"duration_ms", time.Since(start).Milliseconds(),
	)
	return resp, nil
}

func (s *DBService) DeleteTask(ctx context.Context, req *api2.TaskIdRequest) (*api2.DeleteTaskResponse, error) {
	start := time.Now()
	logger.Log.Debugw("RPC DeleteTask start",
		"id", req.Id,
	)

	err := s.repo.DeleteTask(ctx, req.Id)
	if err != nil {
		logger.Log.Errorw("RPC DeleteTask error",
			"id", req.Id,
			"error", err,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return nil, err
	}
	_ = s.repoRedis.DeleteTasks(ctx)
	logger.Log.Infow("RPC DeleteTask success",
		"id", req.Id,
		"duration_ms", time.Since(start).Milliseconds(),
	)
	return &api2.DeleteTaskResponse{}, nil
}

func (s *DBService) MarkDoneTask(ctx context.Context, req *api2.TaskIdRequest) (*api2.DeleteTaskResponse, error) {
	start := time.Now()
	logger.Log.Debugw("RPC MarkDoneTask start",
		"id", req.Id,
	)

	err := s.repo.UpdateTask(ctx, req.Id)
	if err != nil {
		logger.Log.Errorw("RPC MarkDoneTask error",
			"id", req.Id,
			"error", err,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return nil, err
	}
	_ = s.repoRedis.DeleteTasks(ctx)
	logger.Log.Infow("RPC MarkDoneTask success",
		"id", req.Id,
		"duration_ms", time.Since(start).Milliseconds(),
	)
	return &api2.DeleteTaskResponse{}, nil
}
