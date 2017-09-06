package main

import (
	"encoding/json"
	"fmt"
	"github.com/Briscooe/Discogrify/go/caching"
	"github.com/Briscooe/Discogrify/go/logging"
	"github.com/zmb3/spotify"
	"regexp"
	"strings"
)

const formatArtistTracks = "artist:%s:tracks"
const formatArtistSearched = "artist:%s:searched"
const formatSearchArtist = "artist:search:%s"

func GetSearchResultsFromCache(query string, client caching.Client, logger logging.Logger) []spotify.FullArtist {
	key := fmt.Sprintf(formatSearchArtist, query)
	result := client.Get(key)

	fmt.Println(key)

	if len(result) == 0 {
		logger.Printf("%s: Query not found", query)
		return nil
	}

	logger.Printf("%s: Query results found ", query)

	var artists []spotify.FullArtist
	err := json.Unmarshal(result, &artists)

	if err != nil {
		logger.Fatal(err)
	}

	return artists
}
func GetTracksFromCache(id string, client caching.Client, logger logging.Logger) []*spotify.FullAlbum {
	key := fmt.Sprintf(formatArtistTracks, id)
	result := client.Get(key)

	fmt.Println(key)
	if len(result) == 0 {
		logger.Printf("%s: Artist ID not found", id)
		return nil
	}

	logger.Printf("%s: Artist ID found in cache", id)
	IncrementKeyInCache(id, client)

	var tracks []*spotify.FullAlbum
	err := json.Unmarshal(result, &tracks)

	if err != nil {
		logger.Fatal(err)
	}

	return tracks
}

func IncrementKeyInCache(key string, client caching.Client) bool {
	result := false
	if validateKey(key) {
		key = fmt.Sprintf(formatArtistSearched, key)
		result = client.Increment(key)
	}

	return result
}

func AddToCache(key string, value string, client caching.Client, logger logging.Logger, format string) {
	formattedKey := fmt.Sprintf(format, key)
	if validateKey(formattedKey) && validateValue(value) {
		if client.Set(formattedKey, value) {
			logger.Printf("%s: Successfully added key to cache", formattedKey)
		} else {
			logger.Printf("%s: Could not add key to cache", formattedKey)
		}
	} else {
		logger.Printf("%s: Incorrect format\nKey: %s\nValue: %s", formattedKey, formattedKey, value)
	}
}

func validateKey(key string) bool {
	stringSlice := strings.Split(key, ":")
	if len(stringSlice) < 3 {
		return false
	}
	for _, str := range stringSlice {
		if !validateValue(str) {
			return false
		}
	}
	return true
}

func validateValue(value string) bool {
	var js map[string]interface{}
	result := json.Unmarshal([]byte(value), &js) != nil

	if !result {
		result, _ = regexp.MatchString("^[A-Za-z0-9\\S]+$", value)
	}
	return result
}
