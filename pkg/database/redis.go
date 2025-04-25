package database

import (
	"context"
	"fmt"
	"product_recommendation/pkg/config"
	"sync"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client
var once sync.Once

type RedisClientInterface interface {
	RedisPing(ctx context.Context) error
	GetRedisClient() *redis.Client
}

type RedisClient struct {
	RedisClient *redis.Client
}

func NewRedisClient() RedisClientInterface {
	config := config.GetConfig()
	return &RedisClient{
		RedisClient: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", config.Redis.Host, config.Redis.Port),
			Password: config.Redis.Password,
			DB:       config.Redis.DB,
		}),
	}
}

func (r *RedisClient) GetRedisClient() *redis.Client {
	return r.RedisClient
}

func (r *RedisClient) RedisPing(ctx context.Context) error {
	return r.RedisClient.Ping(ctx).Err()
}
