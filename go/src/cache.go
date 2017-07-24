package main

import (
	"encoding/json"
	"github.com/zmb3/spotify"
	"time"
	"strings"
	"regexp"
	"fmt"
)

type CacheClient interface {
	Get(cacheKey string) []byte
	Set(key string, value string, expiration time.Duration) bool
	Increment(key string) bool
}

func GetTracksFromCache(artistId string, client CacheClient) []*spotify.FullAlbum {
	rollingLog.Printf("%s: Checking cache for artist ID", artistId)
	result := client.Get(artistId)

	if len(result) == 0 {
		rollingLog.Printf("%s: Artist ID not found", artistId)
		return nil
	}

	rollingLog.Printf("%s: Artist ID found in cache", artistId)
	var tracks []*spotify.FullAlbum
	err := json.Unmarshal(result, &tracks)

	if err != nil {
		rollingLog.Fatal(err)
	}

	return tracks
}

func IncrementKeyInCache(key string, client CacheClient) bool {
	result := false
	if validateKey(key) {
		result = client.Increment(key)
	}
	return result
}

func validateKey(key string) bool {
	stringSlice := strings.Split(key, ":")
	fmt.Println("KEY = ", key)
	fmt.Println("Slice length = ", len(stringSlice))
	if len(stringSlice) < 2 {
		return false
	}

	for _, str := range stringSlice {
		fmt.Println("Sting = ", str)
		result, _ := regexp.MatchString("[^A-Za-z0-9]+$", str)
		if !result {
			return false
		}35
	}
	return true
}

func AddToCache(key string, value string, expiration time.Duration, client CacheClient) bool {
	result := false
	if validateKey(key) {
		result = client.Set(key, value, expiration)
	}
	return result
}
