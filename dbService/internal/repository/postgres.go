package repository

import (
	"checkdown/dbService/internal/pkg/logger"
	"context"
	"fmt"
	"time"

	dto "checkdown/dbService/internal/DTO" // поправил регистр пакета
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	UserName string `env:"POSTGRES_USER"     env-required:"true"`
	Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
	Host     string `env:"POSTGRES_HOST"     env-required:"true"`
	Port     string `env:"POSTGRES_PORT"     env-required:"true"`
	DBName   string `env:"POSTGRES_DB"       env-required:"true"`
}

type Storage struct {
	pool *pgxpool.Pool
}

// NewStorage открывает пул соединений к БД и логирует успех / ошибку.
func NewStorage(ctx context.Context, cfg Config) (*Storage, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.UserName, cfg.Password, cfg.Host, cfg.Port, cfg.DBName,
	)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}
	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	logger.Log.Infow("pgx pool created", "dsn", dsn)
	return &Storage{pool: pool}, nil
}

// Close — аккуратное закрытие пула.
func (s *Storage) Close() { s.pool.Close() }

// ───────────────────────── CRUD ───────────────────────────────────────────

func (s *Storage) CreateTask(ctx context.Context, task dto.Task) (int, error) {
	const q = `
		INSERT INTO tasks (title, description, is_done, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id;
	`
	now := time.Now()

	start := time.Now()
	var id int
	err := s.pool.QueryRow(ctx, q,
		task.Title,
		task.Description,
		task.IsDone,
		now,
		now,
	).Scan(&id)
	ms := time.Since(start).Milliseconds()

	if err != nil {
		logger.Log.Errorw("sql insert failed", "query", "CREATE", "err", err, "duration_ms", ms)
		return 0, err
	}
	logger.Log.Debugw("sql insert", "id", id, "duration_ms", ms)
	return id, nil
}

func (s *Storage) GetTasks(ctx context.Context) ([]*dto.Task, error) {
	const q = `SELECT id, title, description, is_done, created_at, updated_at FROM tasks;`

	start := time.Now()
	rows, err := s.pool.Query(ctx, q)
	ms := time.Since(start).Milliseconds()
	if err != nil {
		logger.Log.Errorw("sql select failed", "err", err, "duration_ms", ms)
		return nil, fmt.Errorf("get tasks query: %w", err)
	}
	logger.Log.Debugw("sql select", "duration_ms", ms)
	defer rows.Close()

	var tasks []*dto.Task
	for rows.Next() {
		task := &dto.Task{}
		if err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.IsDone,
			&task.CreatedAt,
			&task.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}
	return tasks, nil
}

func (s *Storage) UpdateTask(ctx context.Context, id int64) error {
	const q = `UPDATE tasks SET is_done = $1, updated_at = $2 WHERE id = $3;`
	now := time.Now()

	start := time.Now()
	_, err := s.pool.Exec(ctx, q, "true", now, id)
	ms := time.Since(start).Milliseconds()

	if err != nil {
		logger.Log.Errorw("sql update failed", "id", id, "err", err, "duration_ms", ms)
		return fmt.Errorf("update task: %w", err)
	}
	logger.Log.Debugw("sql update", "id", id, "duration_ms", ms)
	return nil
}

func (s *Storage) DeleteTask(ctx context.Context, id int64) error {
	const q = `DELETE FROM tasks WHERE id = $1;`

	start := time.Now()
	_, err := s.pool.Exec(ctx, q, id)
	ms := time.Since(start).Milliseconds()

	if err != nil {
		logger.Log.Errorw("sql delete failed", "id", id, "err", err, "duration_ms", ms)
		return fmt.Errorf("delete task: %w", err)
	}
	logger.Log.Debugw("sql delete", "id", id, "duration_ms", ms)
	return nil
}
