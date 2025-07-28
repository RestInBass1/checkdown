package grpcHandlers

import (
	"checkdown/dbService/internal/DTO"
	api2 "checkdown/dbService/pkg/api"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

type PostgresRepository interface {
	CreateTask(ctx context.Context, task dto.Task) (int, error)
	GetTasks(ctx context.Context) ([]*dto.Task, error)
	UpdateTask(ctx context.Context, id int64) error
	DeleteTask(ctx context.Context, id int64) error
}

type DBService struct {
	api2.UnimplementedDBServiceServer
	repo PostgresRepository
}

func NewDBService(repo PostgresRepository) *DBService {
	return &DBService{repo: repo}
}

func (s *DBService) AddTask(ctx context.Context, task_r *api2.TaskRequest) (*api2.CreateTaskResponse, error) {
	task := dto.Task{
		Title:       task_r.Title,
		Description: task_r.Description,
	}
	id, err := s.repo.CreateTask(ctx, task)
	if err != nil {
		return nil, err
	}
	createdTask := &api2.CreateTaskResponse{
		Id:    int64(id),
		Error: err.Error(),
	}
	return createdTask, nil
}
func (s *DBService) GetTasks(context.Context, *emptypb.Empty) (*api2.GetTasksResponse, error) {
	tasks, err := s.repo.GetTasks(context.Background())
	if err != nil {
		return nil, err
	}
	resp := &api2.GetTasksResponse{}
	for _, task := range tasks {
		resp.Tasks = append(resp.Tasks, &api2.Task{
			Title:       task.Title,
			Description: task.Description,
		})
	}
	return resp, nil
}
func (s *DBService) DeleteTask(ctx context.Context, id *api2.TaskIdRequest) (*api2.DeleteTaskResponse, error) {
	err := s.repo.DeleteTask(ctx, id.Id)
	if err != nil {
		return nil, err
	}
	return &api2.DeleteTaskResponse{
		Error: err.Error(),
	}, nil
}
func (s *DBService) MarkDoneTask(ctx context.Context, id *api2.TaskIdRequest) (*api2.DeleteTaskResponse, error) {
	err := s.repo.UpdateTask(ctx, id.Id)
	if err != nil {
		return nil, err
	}
	return &api2.DeleteTaskResponse{
		Error: err.Error(),
	}, nil
}
