package main

import (
	"github.com/go-redis/redis"
	"github.com/zmb3/spotify"
	"encoding/json"
	"fmt"
	"time"
)

var (
	redisClient  *redis.Client
)
func setupRedisClient() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:		fmt.Sprintf("%s:%s", config.Redis.Host, config.Redis.Port),
		Password: 	config.Redis.Password,
		DB: 		config.Redis.Db,
	})

	_, err := redisClient.Ping().Result()
	if err != nil {
		rollingLog.Fatalf("Could not connect to Redis host %s:%s\n%v\n", config.Redis.Host, config.Redis.Port, err)
	}
	rollingLog.Printf("Successfully connected to Redis host %s:%s\n", config.Redis.Host, config.Redis.Port)
}

func AddDiscographyToCache(artistId, artistTracks string) bool {

	result := redisClient.Set(artistId, artistTracks, time.Hour * 168)

	if result.Val() == "" {
		return false
	}
	return true
}

func GetDiscographyFromCache(artistId string) []spotify.SimpleTrack {
	result := redisClient.Get(artistId)

	if result.Val() == "" {
		return nil
	}

	var tracks []spotify.SimpleTrack
	bytes, _ := result.Bytes()
	err := json.Unmarshal(bytes, &tracks)

	if err != nil {
		rollingLog.Fatal(err)
	}

	return tracks
}