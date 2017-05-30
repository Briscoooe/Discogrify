package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

var stateString string

// Index ...
func Index(w http.ResponseWriter, r *http.Request) {
	// wait for auth to complete
	client := <-ch

	fmt.Println("HERE")
	// use the client to make calls that require authorization
	user, err := client.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}

	type SpotifyUser struct {
		ID string `json:"id"`
	}
	spotifyUser := SpotifyUser{
		ID: user.ID,
	}
	if err := json.NewEncoder(w).Encode(spotifyUser); err != nil {
		panic(err)
	}
	fmt.Println("You are logged in as:", user.ID)
}

// GetPlaylists ...
func GetPlaylists(w http.ResponseWriter, r *http.Request) {
	stateString = auth.AuthURL(GenerateStateString())
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	type LoginInformation struct {
		LoginURL  string     `json:"loginUrl"`
		Playlists []Playlist `json:"playlists"`
	}
	loginStuff := LoginInformation{
		LoginURL:  stateString,
		Playlists: playlists,
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
	var playlist Playlist

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048476))

	if err != nil {
		panic(err)
	}

	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &playlist); err != nil {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	stateString = auth.AuthURL(GenerateStateString())

	url := auth.AuthURL(stateString)

	x := map[string]string { "url" : url}

	if err := json.NewEncoder(w).Encode(x); err != nil {
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
}

func FollowPlaylist(w http.ResponseWriter, r *http.Request) {

	// Follow playlist on spotify
}
