package redis

import (
	"fmt"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	Client *redis.Client
}

func Connect(addr, password string, idx int) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       idx,
	})
	_, err := client.Ping(client.Context()).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}
	return &Redis{Client: client}, nil
}

func (r *Redis) Close() error {
	return r.Client.Close()
}
