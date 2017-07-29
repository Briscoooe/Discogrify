package main

import (
	"encoding/json"
	"github.com/zmb3/spotify"
	"time"
	"strings"
	"regexp"
	"github.com/Briscooe/Discogrify/go/caching"
)


func GetTracksFromCache(artistId string, client caching.Client) []*spotify.FullAlbum {
	result := client.Get(artistId)

	if len(result) == 0 {
		return nil
	}

	var tracks []*spotify.FullAlbum
	err := json.Unmarshal(result, &tracks)

	if err != nil {
		//rollingLog.Fatal(err)
	}

	return tracks
}

func IncrementKeyInCache(key string, client caching.Client) bool {
	result := false
	if validateKey(key) {
		result = client.Increment(key)
	}

	return result
}

func AddToCache(key string, value string, expiration time.Duration, client caching.Client) bool {
	result := false
	if validateKey(key) {
		if validateValue(value) {
			if validateExpiration(expiration) {
				result = client.Set(key, value, expiration)
			}
		}
	}
	return result
}

func validateKey(key string) bool {
	stringSlice := strings.Split(key, ":")
	if len(stringSlice) < 3 {
		return false
	}
	for _, str := range stringSlice {
		if !validateValue(str){
			return false
		}
	}
	return true
}

func validateValue(value string) bool {
	result, _ := regexp.MatchString("^[A-Za-z0-9\\S]+$", value)
	return result
}

func validateExpiration(expiration time.Duration) bool {
	if expiration < time.Hour * 24 {
		return false
	}
	return true
}