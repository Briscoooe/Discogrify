package main

import (
	"encoding/json"
	"github.com/zmb3/spotify"
	"time"
)

type CacheClient interface {
	Get(artistId string) []byte
	Set(artistId string, artistTracks string, expireIn time.Duration) bool
}

func GetDiscographyFromCache(artistId string, client CacheClient) []*spotify.FullAlbum {
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

func AddDiscographyToCache(artistId, artistTracks string, client CacheClient) bool {
	return client.Set(artistId, artistTracks, time.Hour*168)
}
