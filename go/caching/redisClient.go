package caching

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/Briscooe/Discogrify/go/logging"
	"time"
)

type RedisClient struct {
	RedisClient *redis.Client
	Logger logging.Logger
	Expiration time.Duration
}

func NewRedisClient(logger logging.Logger, host, port, password string, db int, expiration int) *RedisClient {
	logger = logger
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       db,
	})

	_, err := redisClient.Ping().Result()
	if err != nil {
		logger.LogErr(err,"Could not connect to Redis host %s:%s\n%v\n", host, port)
	}

	logger.Logf("Successfully connected to Redis host %s:%s", host, port)
	return &RedisClient{
		RedisClient: redisClient,
		Logger: logger,
		Expiration: time.Hour * time.Duration(expiration),
	}
}

func (r *RedisClient) Get(key string) []byte {
	result := r.RedisClient.Get(key)
	bytes, _ := result.Bytes()
	return bytes
}

func (r *RedisClient) Set(key string,  value string) bool {
	result := r.RedisClient.Set(key, value, r.Expiration)
	return result.Val() != ""
}

func (r *RedisClient) Increment(key string) bool {
	result := r.RedisClient.Incr(key)
	return result.Val() != 0
}

