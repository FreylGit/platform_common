package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type redisCash struct {
	client *redis.Client
}

type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	GetModel(ctx context.Context, key string, model interface{}) error
	SetModel(ctx context.Context, key string, value interface{}, expiration time.Duration) error
}

func NewCash(ctx context.Context, address string) (Cache, error) {
	oprs, err := redis.ParseURL(address)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(oprs)
	err = client.FlushDB(ctx).Err()
	if err != nil {
		return nil, err
	}
	return &redisCash{client: client}, nil
}

func (r *redisCash) Get(ctx context.Context, key string) (string, error) {
	cachedData, err := r.client.Get(context.Background(), key).Result()
	if err != nil {
		return "", err
	}

	return cachedData, nil
}

func (r *redisCash) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	err := r.client.Set(ctx, key, value, expiration).Err()
	return err
}

func (r *redisCash) GetModel(ctx context.Context, key string, model interface{}) error {
	err := r.client.Get(ctx, key).Scan(model)
	return err
}

func (r *redisCash) SetModel(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := r.client.Set(ctx, key, value, expiration).Err()
	return err
}
