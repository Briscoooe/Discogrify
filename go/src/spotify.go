package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os/exec"
	"strings"

	"github.com/zmb3/spotify"
)

const redirectURI = "http://localhost:8080/callback"

var (
	auth       = spotify.NewAuthenticator(redirectURI, spotify.ScopePlaylistModifyPrivate, spotify.ScopePlaylistModifyPublic)
	ch         = make(chan *spotify.Client)
	characters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
)

var playlists []Playlist

func init() {
	RepoCreatePlaylist(Playlist{Title: "Write presentation"})
	RepoCreatePlaylist(Playlist{Title: "Host meetup"})
}

func RepoCreatePlaylist(p Playlist) Playlist {
	p.Followers = 10
	playlists = append(playlists, p)
	return p
}

type Artist struct {
	Name string `json:"name"`
}

func GetAllSongsFromSpotify(decoder *json.Decoder) Playlist {

	var artist Artist
	err := decoder.Decode(&artist)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println(artist.Name)

	artistId := GetArtistId(artist.Name)

	fmt.Println("Artist ID")
	fmt.Println(artistId)
	playlist := Playlist{ID: ""}
	return playlist
}

type ArtistSearchResponse struct {
	Artists struct {
		Items []struct {
			Href string `json:"href"`
		} `json:"items"`
	} `json:"artists"`
}

func GetArtistId(artistName string) string {
	// curl -s 'https://api.spotify.com/v1/search?q=Daft+Punk&type=artist' | jq -r '.artists.items[0].href'

	artistName = url.QueryEscape(artistName)
	fmt.Println(artistName)
	curlURL := fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&type=artist", artistName)

	out, _ := exec.Command("curl", "-s", curlURL).Output()

	var response ArtistSearchResponse
	json.Unmarshal(out, &response)

	splitString := strings.Split(response.Artists.Items[0].Href, "/")
	return splitString[len(splitString)-1]
}

func GenerateStateString() string {
	b := make([]rune, 16)
	for i := range b {
		b[i] = characters[rand.Intn(len(characters))]
	}
	return string(b)
}

func LoginToSpotify(w http.ResponseWriter, r *http.Request, state string) {


}
