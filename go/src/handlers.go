package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/zmb3/spotify"
)

var stateString string

// Index ...
func Index(w http.ResponseWriter, r *http.Request) {
	user, err := client.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(user.ID)
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
func GetSongsByArtist(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	decoder := json.NewDecoder(r.Body)

	fmt.Println(r.Body)
	var playlist = GetAllSongsFromSpotify(decoder)

	if err := json.NewEncoder(w).Encode(playlist); err != nil {
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
	LoginToSpotify(w, r, stateString)
}

func FollowPlaylist(w http.ResponseWriter, r *http.Request) {

	// Follow playlist on spotify
}
