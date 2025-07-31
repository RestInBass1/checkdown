package repository

import (
	dto "checkdown/dbService/internal/DTO"
	"checkdown/dbService/internal/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisConfig struct {
	Host     string `env:"REDIS_HOST" env-required:"true"`
	Port     string `env:"REDIS_PORT" env-required:"true"`
	Password string `env:"REDIS_PASSWORD" env-required:"true"`
	TTL      int    `env:"REDIS_TTL"      env-default:"300"` // секунд
}

const tasksKey = "tasks:all" // один глобальный ключ, если пока нет фильтров

func NewRedis(cfg RedisConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		panic(fmt.Errorf("failed to connect Redis:%w", err))
	}
	logger.Log.Infow("redis connected", "addr", cfg.Host+":"+cfg.Port)
	return client
}

type RedisRepository struct {
	Client *redis.Client
	ttl    time.Duration
}

func NewRedisRepository(redisClient *redis.Client, ttlSeconds int) *RedisRepository {
	return &RedisRepository{
		Client: redisClient,
		ttl:    time.Second * time.Duration(ttlSeconds),
	}
}

func (r *RedisRepository) GetTasks(ctx context.Context) ([]*dto.Task, error) {
	start := time.Now()
	raw, err := r.Client.Get(ctx, tasksKey).Bytes()
	switch {
	case err == nil:
		var tasks []*dto.Task
		if err := json.Unmarshal(raw, &tasks); err != nil {
			logger.Log.Warnw("redis unmarshal failed", "key", tasksKey, "err", err)
			// битые данные — лучше удалить ключ и заставить сходить в БД
			_ = r.Client.Del(ctx, tasksKey).Err()
			return nil, redis.Nil
		}
		logger.Log.Debugw("redis cache hit", "key", tasksKey, "len", len(tasks), "ms", time.Since(start).Milliseconds())
		return tasks, nil
	case err == redis.Nil:
		logger.Log.Debugw("redis cache miss", "key", tasksKey)
		return nil, err
	default:
		logger.Log.Warnw("redis get failed", "key", tasksKey, "err", err)
		return nil, err
	}
}

func (r *RedisRepository) SetTasks(ctx context.Context, tasks []*dto.Task) error {
	raw, err := json.Marshal(tasks)
	if err != nil {
		logger.Log.Warnw("redis marshal failed", "key", tasksKey, "err", err)
		return err
	}
	if err := r.Client.Set(ctx, tasksKey, raw, r.ttl).Err(); err != nil {
		logger.Log.Warnw("redis set failed", "key", tasksKey, "err", err)
	}
	logger.Log.Debugw("redis set", "key", tasksKey, "len", len(tasks), "ttl", r.ttl.Seconds())
	return nil
}

func (r *RedisRepository) DeleteTasks(ctx context.Context) error {
	return r.Client.Del(ctx, tasksKey).Err()
}
