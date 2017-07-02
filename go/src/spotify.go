package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	//"os/exec"
	//"strings"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
	"os"
	"context"
)

const redirectURI = "http://localhost:8080/callback"

var (
	token 	     = ""
	ch           = make(chan *spotify.Client)
	client       spotify.Client
	auth         = spotify.NewAuthenticator(redirectURI, spotify.ScopePlaylistModifyPrivate, spotify.ScopePlaylistModifyPublic)
	unauthClient = spotify.DefaultClient
	characters   = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
)

type Artist struct {
	Name string `json:"name"`
}

func GetAllSongsFromSpotify(decoder *json.Decoder) spotify.FullPlaylist {
	var artist Artist
	err := decoder.Decode(&artist)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	artistID := GetArtistId(artist.Name)

	if artistID != "" {
		var spotifyID = spotify.ID(artistID)
		user, err := client.CurrentUser()

		if err != nil {
			fmt.Println(err)
		}

		playlist, err := client.CreatePlaylistForUser(user.ID, artistID, true)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(playlist)
		fmt.Println(user.ID)
		fmt.Println(artistID)
		fmt.Println(spotifyID)

		//getSongsFromAlbums(spotifyID, &playlist)
		//getSongsFromSingles(spotifyID, &playlist)
		//getSongsFromCompilations(spotifyID, &playlist)
		//getSongsFromAppearsOn(spotifyID, &playlist)

	}

	var playlist spotify.FullPlaylist
	return playlist
}

func getSongsFromAlbums(ID spotify.ID, p **spotify.FullPlaylist) {
	fmt.Println("Getting album songs")
	//albums, _ := unauthClient.GetArtistAlbums(ID)

	fmt.Println(p)
	/*
		for _, albumReference := range albums.Albums {
			albumDetails, _ := unauthClient.GetAlbumTracks(albumReference.ID)
			for _, track := range albumDetails.Tracks {
				fmt.Println(track)
			}
		}*/
}

/*
func (c *spotify.Client) GetArtistSingles(artistID spotify.ID) (*SimpleAlbumPage, error) {

}
func getSongsFromSingles(ID spotify.ID, p **spotify.FullPlaylist) {

}

func getSongsFromCompilations(ID spotify.ID, p **spotify.FullPlaylist) {

}

func getSongsFromAppearsOn(ID spotify.ID, p *spotify.FullPlaylist) {

}*/

type ArtistSearchResponse struct {
	Artists struct {
		Items []struct {
			Href string `json:"href"`
		} `json:"items"`
	} `json:"artists"`
}

func GetArtistId(artistName string) string {
	artistName = url.QueryEscape(artistName)
	//curlURL := fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&type=artist", artistName)


	//tokenFmt := fmt.Sprintf("\"Authorization: Bearer %s\"", token)

	//out, _ := exec.Command("curl", "-X", "GET", "-H", "\"Accept: application/json\"", "-H", tokenFmt,  curlURL).Output()

	//fmt.Println(out)
	//var response ArtistSearchResponse
	//json.Unmarshal(out, &response)

	//splitString := strings.Split(response.Artists.Items[0].Href, "/")
	//return splitString[len(splitString)-1]


	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		TokenURL:     spotify.TokenURL,
	}
	token, err := config.Token(context.Background())
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	client := spotify.Authenticator{}.NewClient(token)

	result, err := client.Search(artistName, spotify.SearchTypeArtist)

	if (err != nil) {
		log.Fatal(err)
	}

	if(result.Artists != nil) {
		fmt.Println("Artists")
		for _, item := range result.Artists.Artists{
			fmt.Println("   ", item.Name)
		}
	}

	return ""
}

func GenerateStateString() string {
	b := make([]rune, 16)
	for i := range b {
		b[i] = characters[rand.Intn(len(characters))]
	}
	return string(b)
}

func LoginToSpotify(w http.ResponseWriter, r *http.Request, state string) {
	tok, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	// use the token to get an authenticated client
	client := auth.NewClient(tok)
	fmt.Println(w, "Login Completed!")

	client.CurrentUser()
}
