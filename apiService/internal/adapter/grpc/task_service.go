package grpcadapter

import (
	"checkdown/apiService/internal/pkg/logger"
	"checkdown/apiService/internal/usecase"
	"checkdown/dbService/pkg/api"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

// в main напишем че то типо conn, err := grpc.Dial("db_service:50051", grpc.WithInsecure())
// cli := api.NewDBServiceClient(conn)
// New(cli) -> service
// service.c.AddTask
type service struct{ c api.DBServiceClient }

func New(c api.DBServiceClient) usecase.TaskService {
	return &service{c: c}
}

func (s *service) Create(ctx context.Context, t usecase.Task) (int64, error) {
	logger.Log.Debugw("grpc call AddTask", "title", t.Title)
	resp, err := s.c.AddTask(ctx, &api.TaskRequest{
		Title:       t.Title,
		Description: t.Description,
	})
	if err != nil {
		logger.Log.Errorw("AddTask failed", "err", err)
		return 0, err // ошибка сети, таймаута и проч.
	}
	logger.Log.Infow("task created via gRPC", "id", resp.Id)
	return resp.Id, nil // вернули сгенерированный id
}
func (s *service) List(ctx context.Context) ([]usecase.Task, error) {
	resp, err := s.c.GetTasks(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	out := make([]usecase.Task, 0, len(resp.Tasks))
	for _, v := range resp.Tasks {
		out = append(out, usecase.Task{
			ID:          v.Id,
			Title:       v.Title,
			Description: v.Description,
			IsDone:      v.IsDone,
			CreatedAt:   v.CreatedAt.AsTime(),
		})
	}
	return out, nil
}

func (s *service) Delete(ctx context.Context, id int64) error {
	_, err := s.c.DeleteTask(ctx, &api.TaskIdRequest{Id: id})
	if err != nil {
		logger.Log.Errorw("DeleteTask failed", "id", id, "err", err)
	} else {
		logger.Log.Infow("task deleted", "id", id)
	}
	return err
}

func (s *service) MarkDone(ctx context.Context, id int64) error {
	_, err := s.c.MarkDoneTask(ctx, &api.TaskIdRequest{Id: id})
	return err
}
