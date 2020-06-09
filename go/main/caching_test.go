package main

import (
	"encoding/json"
	"github.com/Briscoooe/Discogrify/go/logging"
	"github.com/Briscoooe/spotify"
	"testing"
)

type FakeCacheClient struct {
}

func (f FakeCacheClient) Get(artistId string) []byte {
	validArtists := []string{"artist:ArtistID1:tracks",
		"artist:ArtistID2:tracks"}

	for _, artist := range validArtists {
		if artistId == artist {
			var trackPage spotify.SimpleTrackPage
			trackPage.Tracks = append(trackPage.Tracks, spotify.SimpleTrack{
				ID: "TrackID"},
			)
			var album spotify.FullAlbum
			album.Tracks = trackPage
			var albums []*spotify.FullAlbum
			albums = append(albums, &album)
			bytes, _ := json.Marshal(albums)
			return bytes
		}
	}
	return nil
}

func (f FakeCacheClient) Set(artistId string, artistTracks string) bool {
	return true
}

func (f FakeCacheClient) Increment(key string) bool {
	return true
}
func TestGetTracksFromCache(t *testing.T) {
	var testData = []struct {
		artistId               string
		expectedTracksReturned int
	}{
		{"ArtistID1", 1},
		{"ArtistID2", 1},
		{"ArtistID3", 0},
	}
	var cacheClient FakeCacheClient
	for _, test := range testData {
		result := GetTracksFromCache(test.artistId, cacheClient, logging.NewFakeLogger())
		if len(result) != test.expectedTracksReturned {
			t.Errorf("Input: %s\nExpected %d\nOutput%d\n", test.artistId, test.expectedTracksReturned, len(result))
		}
	}
}

func TestIncrementKeyInCache(t *testing.T) {
	var testData = []struct {
		key            string
		expectedOutput bool
	}{
		{"artist:ID1:tracks", true},
		{"user:ID2:loggedin", true},
		{"artist:ID3:searched", true},
		{"artist:helloworld:tracks", true},
		{"::", false},
		{"0:0", false},
		{"false", false},
		{"null", false},
		{"0", false},
	}

	var cacheClient FakeCacheClient
	for _, test := range testData {
		result := IncrementKeyInCache(test.key, cacheClient)
		if result != test.expectedOutput {
			t.Errorf("Input: %s\nExpected %t\nOutput %t\n", test.key, test.expectedOutput, result)
		}
	}
}

func TestAddToCache(t *testing.T) {
	var testData = []struct {
		key            string
		value          string
		expectedOutput bool
	}{
		{"artist:ID1:tracks", "testString", true},
		{"artist:ID2:tracks", "testString", true},
		{"artist:ID3:tracks", "hello world", true},
		{"::", "validString", false},
	}

	var cacheClient FakeCacheClient
	for _, test := range testData {
		result := AddToCache(test.key, test.value, cacheClient, logging.NewFakeLogger(), "")
		if result != test.expectedOutput {
			t.Errorf("Input: Key: %s\nValue: %s\nExpected %t\nOutput %t\n", test.key, test.value, test.expectedOutput, result)
		}
	}
}
