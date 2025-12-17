package cache

import (
	"context"
	"log"
	"time"

	"pragma/internal/config"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(cfg config.RedisConfig) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		DB:       cfg.DB,
		Password: cfg.Password,
	})

	var lastErr error
	for i := 0; i < 5; i++ {
		if _, lastErr = client.Ping(ctx).Result(); lastErr == nil {
			log.Println("Connected to Redis.")
			return &RedisClient{Client: client}, nil
		}
		log.Printf("redis ping failed (attempt %d/5): %v", i+1, lastErr)
		time.Sleep(2 * time.Second)
	}

	return nil, lastErr
}

func (r *RedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	return r.Client.Set(ctx, key, value, expiration).Err()
}

func (r *RedisClient) Get(key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

func (r *RedisClient) Delete(key string) error {
	return r.Client.Del(ctx, key).Err()
}

func (r *RedisClient) Close() error {
	if r.Client == nil {
		return nil
	}
	return r.Client.Close()
}
