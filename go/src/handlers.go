package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/zmb3/spotify"
	//"os/user"
	"github.com/gorilla/mux"
)

var (
	stateString string
	cacheClient CacheClient
)

func registerCacheClient(client CacheClient) {
	cacheClient = client
}

func Callback(w http.ResponseWriter, r *http.Request) {
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

	token = tok.AccessToken
	user, err := client.CurrentUser()
	if err != nil {
		rollingLog.Fatal(err)
	}

	fmt.Println("Logged in as: ", user.ID)

	http.Redirect(w, r, "/#", http.StatusAccepted)

}

// Index ...
func Index(w http.ResponseWriter, r *http.Request) {

}

// GetPlaylists ...
func GetPlaylists(w http.ResponseWriter, r *http.Request) {
	stateString = auth.AuthURL(GenerateStateString())
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	type LoginInformation struct {
		LoginURL  string                   `json:"loginUrl"`
		Playlists []spotify.SimplePlaylist `json:"playlists"`
	}
	loginStuff := LoginInformation{
		LoginURL:  stateString,
		Playlists: nil,
	}

	if err := json.NewEncoder(w).Encode(loginStuff); err != nil {
		panic(err)
	}
}

// GetSongsByArtist ...
func SearchForArtistHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	vars := mux.Vars(r)

	var artists = SearchForArtist(vars["name"])

	if err := json.NewEncoder(w).Encode(artists); err != nil {
		panic(err)
	}
	// Get songs from albums, singles, appears on, compilations
}

func GetPlaylist(w http.ResponseWriter, r *http.Request) {

	// Get playlist from Redis
}

func PublishPlaylist(w http.ResponseWriter, r *http.Request) {

}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	stateString = GenerateStateString()

	url := auth.AuthURL(stateString)

	urlJson := map[string]string{"url": url}

	if err := json.NewEncoder(w).Encode(urlJson); err != nil {
		panic(err)
	}
	//http.Redirect(w, r, url, 400)
	//
	//tok, err := auth.Token(stateString, r)
	//if err != nil {
	//	http.Error(w, "Couldn't get token", http.StatusForbidden)
	//	log.Fatal(err)
	//}
	//if st := r.FormValue("state"); st != stateString {
	//	http.NotFound(w, r)
	//	log.Fatalf("State mismatch: %s != %s\n", st, stateString)
	//}
	//// use the token to get an authenticated client
	//client := auth.NewClient(tok)
	//fmt.Fprintf(w, "Login Completed!")
	//ch <- &client
	////LoginToSpotify(w, r, stateString)
	//LoginToSpotify(w, r, stateString)
}

func FollowPlaylist(w http.ResponseWriter, r *http.Request) {

	// Follow playlist on spotify
}

func GetTracksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	vars := mux.Vars(r)
	artistId := vars["artistId"]

	artistTracks := GetDiscographyFromCache(artistId, cacheClient)

	if artistTracks == nil {
		artistTracks = GetDiscographyFromSpotify(artistId)
		tracksJson, _ := json.Marshal(artistTracks)
		if AddDiscographyToCache(artistId, string(tracksJson), cacheClient) {
			rollingLog.Printf("%s: Successfully added artist to cache", artistId)
		} else {
			rollingLog.Printf("%s: Could not add artist to cache", artistId)
		}
	}

	rollingLog.Printf("%s: Returning tracks", artistId)
	if err := json.NewEncoder(w).Encode(artistTracks); err != nil {
		panic(err)
	}
}
