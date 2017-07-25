package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	//"os/user"
	"github.com/gorilla/mux"
	"time"
)

var (
	stateString string
)

func indexHandlerFunc(w http.ResponseWriter, r *http.Request) {
	// Do something
}

func loginToSpotifyHandlerFunc(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(stateString, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		rollingLog.Fatal(err)
	}
	if st := r.FormValue("state"); st != stateString {
		http.NotFound(w, r)
		rollingLog.Fatalf("State mismatch: %s != %s\n", st, stateString)
	}
	// use the token to get an authenticated client
	client := auth.NewClient(tok)
	fmt.Println(w, "Login Completed!")

	client.CurrentUser()
}

func getTracksHandler(cacheClient CacheClient) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		vars := mux.Vars(r)
		artistId := vars["artistId"]

		rollingLog.Printf("%s: Checking cache for artist ID", artistId)
		artistTracks := GetTracksFromCache("artist:" + artistId + ":tracks", cacheClient)

		if artistTracks == nil {
			rollingLog.Printf("%s: Artist ID not found", artistId)
			artistTracks = GetDiscographyFromSpotify(artistId)
			tracksJson, _ := json.Marshal(artistTracks)
			if AddToCache("artist:" + artistId + ":tracks", string(tracksJson), time.Hour * 168, cacheClient) {
				IncrementKeyInCache("artist:" + artistId + ":searched", cacheClient)
				rollingLog.Printf("%s: Successfully added artist to cache", artistId)
			} else {
				rollingLog.Printf("%s: Could not add artist to cache", artistId)
			}
			fmt.Println(string(tracksJson))
		} else {
			rollingLog.Printf("%s: Artist ID found in cache", artistId)
		}

		rollingLog.Printf("%s: Returning tracks", artistId)
		if err := json.NewEncoder(w).Encode(artistTracks); err != nil {
			panic(err)
		}
	})
}

func callbackHandler(cacheClient CacheClient) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tok, err := auth.Token(stateString, r)
		if err != nil {
			http.Error(w, "Couldn't get token", http.StatusForbidden)
			rollingLog.Fatal(err)
		}
		if st := r.FormValue("state"); st != stateString {
			http.NotFound(w, r)
			rollingLog.Fatalf("State mismatch: %s != %s\n", st, stateString)
		}
		// 	rollingLog.Printf("Adding user %s to cache", userId)

		// use the token to get an authenticated client
		client := auth.NewClient(tok)

		token = tok.AccessToken
		user, err := client.CurrentUser()
		if err != nil {
			rollingLog.Fatal(err)
		}

		fmt.Println("Logged in as: ", user.ID)

		http.Redirect(w, r, "/#", http.StatusAccepted)
	})
}

func searchArtistHandler(cacheClient CacheClient) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		vars := mux.Vars(r)

		var artists = SearchForArtist(vars["name"], cacheClient)

		if err := json.NewEncoder(w).Encode(artists); err != nil {
			panic(err)
		}
	})
}

func publishPlaylistHandle(cacheClient CacheClient) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		fmt.Println(r.Body)
		fmt.Println(r.GetBody())
	})
}
