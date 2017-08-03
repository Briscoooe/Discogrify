package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/Briscooe/Discogrify/go/logging"
	//"os/user"
	"github.com/gorilla/mux"
	"time"
	"github.com/Briscooe/Discogrify/go/caching"
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
			fmt.Println(string(tracksJson))
		} else {
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

		err, msg := spotify.ValidateCallback(r)

		if err != nil {
			logger.Fatal(err, msg)
		}

		http.Redirect(w, r, "/#", http.StatusAccepted)
	})
}

func searchArtistHandler(cacheClient caching.Client, spotify Spotify) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		vars := mux.Vars(r)

		var artists = spotify.SearchForArtist(vars["name"], cacheClient)

		if err := json.NewEncoder(w).Encode(artists); err != nil {
			panic(err)
		}
	})
}

func publishPlaylistHandle(cacheClient caching.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		fmt.Println(r.Body)
		fmt.Println(r.GetBody())
	})
}
