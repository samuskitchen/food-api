package database

import (
	"food-api/infrastructure/auth"
	"github.com/go-redis/redis/v8"
)

type RedisService struct {
	Auth   auth.InterfaceAuth
	Client *redis.Client
}

func NewRedisDB(host, port, password string) (*RedisService, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       0,
	})

	redisService := &RedisService{
		Auth:   auth.NewAuth(redisClient),
		Client: redisClient,
	}

	return redisService, nil
}
