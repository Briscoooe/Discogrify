package main

import (
	"encoding/json"
	"net/http"
	"github.com/Briscooe/Discogrify/go/logging"
	"github.com/gorilla/mux"
	"time"
	"github.com/Briscooe/Discogrify/go/caching"
	"io/ioutil"
	"strings"
)

var (
	stateString string
)

func indexHandlerFunc(logger logging.Logger) http.Handler {
	// Do something
	return nil
}

func loginToSpotifyHandlerFunc(logger logging.Logger, spotify Spotify) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := spotify.GenerateLoginUrl()

		if err := json.NewEncoder(w).Encode(url); err != nil {
			panic(err)
		}
	})
}

func getTracksHandler(cacheClient caching.Client, logger logging.Logger, spotify Spotify) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		vars := mux.Vars(r)
		artistId := vars["artistId"]

		logger.Printf("%s: Checking cache for artist ID", artistId)
		artistTracks := GetTracksFromCache("artist:" + artistId + ":tracks", cacheClient)

		if artistTracks == nil {
			logger.Printf("%s: Artist ID not found", artistId)
			artistTracks = spotify.GetDiscography(artistId)
			tracksJson, _ := json.Marshal(artistTracks)
			if AddToCache("artist:" + artistId + ":tracks", string(tracksJson), time.Hour * 168, cacheClient) {
				IncrementKeyInCache("artist:" + artistId + ":searched", cacheClient)
				logger.Printf("%s: Successfully added artist to cache", artistId)
			} else {
				logger.Printf("%s: Could not add artist to cache", artistId)
			}
		} else {
			IncrementKeyInCache("artist:" + artistId + ":searched", cacheClient)
			logger.Printf("%s: Artist ID found in cache", artistId)
		}

		logger.Printf("%s: Returning tracks", artistId)
		if err := json.NewEncoder(w).Encode(artistTracks); err != nil {
			panic(err)
		}
	})
}

func callbackHandler(cacheClient caching.Client, logger logging.Logger, spotify Spotify) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err, msg := spotify.ValidateCallback(r)

		if err != nil {
			logger.Fatal(err, msg)
		}

		http.Redirect(w, r, "/", 302)

		if err := json.NewEncoder(w).Encode(user); err != nil {
			panic(err)
		}
	})
}

func searchArtistHandler(cacheClient caching.Client, logger logging.Logger, spotify Spotify) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		w.WriteHeader(http.StatusOK)

		vars := mux.Vars(r)

		query := vars["name"]
		logger.Printf("%s: Checking cache for search query ", query)
		results := GetSearchResultsFromCache("artist:search:" + query, cacheClient)
		if len(results) == 0 {
		logger.Printf("%s: Query not found", query)
			results = spotify.SearchForArtist(query, cacheClient)
		} else {
			logger.Printf("%s: Query results found ", query)
		}

		if err := json.NewEncoder(w).Encode(results); err != nil {
			panic(err)
		}
	})
}

func publishPlaylistHandler(cacheClient caching.Client, s Spotify) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		body, _ := ioutil.ReadAll(r.Body)
		tracks := strings.Split(string(body), ",")

		result := s.PublishPlaylist(tracks)

		if err := json.NewEncoder(w).Encode(result); err != nil {
			panic(err)
		}

	})
}
