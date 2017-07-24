package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

type RedisClient struct {
	redisClient *redis.Client
}

func (r *RedisClient) SetupClient() {
	r.redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Redis.Host, config.Redis.Port),
		Password: config.Redis.Password,
		DB:       config.Redis.Db,
	})

	_, err := r.redisClient.Ping().Result()
	if err != nil {
		rollingLog.Fatalf("Could not connect to Redis host %s:%s\n%v\n", config.Redis.Host, config.Redis.Port, err)
	}
	rollingLog.Printf("Successfully connected to Redis host %s:%s\n", config.Redis.Host, config.Redis.Port)
}

func (r *RedisClient) Get(key string) []byte {
	result := r.redisClient.Get(key)
	artistTracks, _ := result.Bytes()

	return artistTracks
}

func (r *RedisClient) Set(key string,  value string, expireIn time.Duration) bool {
	result := r.redisClient.Set(key, value, expireIn)
	return result.Val() != ""
}

func (r *RedisClient) Increment(key string) bool {
	result := r.redisClient.Incr(key)
	return result.Val() != 0
}

