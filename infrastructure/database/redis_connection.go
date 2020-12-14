package database

import (
	"food-api/infrastructure/auth"
	"github.com/go-redis/redis/v8"
	"os"
	"sync"
)

var (
	clientRedis *RedisService
	onceRedis   sync.Once
)

type RedisService struct {
	Auth   auth.InterfaceAuth
	Client *redis.Client
}

func NewRedisDB() *RedisService {
	onceRedis.Do(initRedisDB)
	return clientRedis
}

func initRedisDB() {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: redisPassword,
		DB:       0,
	})

	clientRedis = &RedisService{
		Auth:   auth.NewAuth(redisClient),
		Client: redisClient,
	}
}
