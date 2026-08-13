package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore interface {
	Set(
		ctx context.Context,
		key string,
		value string,
		ttl time.Duration,
	) error

	Get(
		ctx context.Context,
		key string,
	) (string, error)

	Delete(
		ctx context.Context,
		key string,
	) error

	Exists(
		ctx context.Context,
		key string,
	) (bool, error)
}
type Store struct {
	client *redis.Client
}

func NewStore(client *redis.Client) *Store {
	return &Store{
		client: client,
	}
}

func (s *Store) Set(
	ctx context.Context,
	key string,
	value string,
	ttl time.Duration,
) error {
	return s.client.Set(ctx, key, value, ttl).Err()
}

func (s *Store) Get(
	ctx context.Context,
	key string,
) (string, error) {
	return s.client.Get(ctx, key).Result()
}

func (s *Store) Delete(
	ctx context.Context,
	key string,
) error {
	return s.client.Del(ctx, key).Err()
}

func (s *Store) Exists(
	ctx context.Context,
	key string,
) (bool, error) {
	result, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return result > 0, nil
}
