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

func (r *RedisClient) Get(artistId string) (artistTracks []byte) {
	result := r.redisClient.Get(artistId)
	returnBytes, _ := result.Bytes()

	return returnBytes
}

func (r *RedisClient) Set(artistId string, artistTracks string, expireIn time.Duration) (setOrNot bool) {
	r.redisClient.Incr(artistId + ":searched")
	result := r.redisClient.Set(artistId+":tracks", artistTracks, time.Hour*168)
	return result.Val() != ""
}
