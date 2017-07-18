package main

import (
	"github.com/go-redis/redis"
	"log"
	"github.com/zmb3/spotify"
	"encoding/json"
	"fmt"
)

var (
	redisClient  *redis.Client
)
func setupClient() {
	fmt.Printf("%v:%v", config.Redis.Host, config.Redis.Port)
	redisClient = redis.NewClient(&redis.Options{
		Addr:		fmt.Sprintf("%s:%s", config.Redis.Host, config.Redis.Port),
		Password: 	config.Redis.Password,
		DB: 		config.Redis.Db,
	})

	_, err := redisClient.Ping().Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis host %s:%s\n%v\n", config.Redis.Host, config.Redis.Port, err)
	}
	log.Printf("Successfully connected to Redis host %s:%s\n", config.Redis.Host, config.Redis.Port)
}

func AddDiscographyToCache(artistId, artistTracks string) bool {

	result := redisClient.Set(artistId, artistTracks, 0)

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
		log.Fatal(err)
	}

	return tracks
}