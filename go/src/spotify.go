package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	//"os/exec"
	//"strings"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
	"os"
	"context"
	"sync"
	"time"
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
	Id   string `json:"id"`
}

func GetAllSongsFromSpotify(decoder *json.Decoder) spotify.FullPlaylist {
	var artist Artist
	err := decoder.Decode(&artist)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	//artistID := GetArtistId(artist.Name)

	/*if artistID != "" {
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
*/
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

func GetAllSongsByArtist(artistId string) []spotify.SimpleTrack {
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

	// 1 = Album
	// 2 = Single
	// 4 = AppearsOn
	// 5 = Compilations

	done := make(chan bool)
	allAlbums := make(map[spotify.ID]spotify.SimpleAlbum)
	limit := 50
	albumTypes := []int{2}
	for _, albumType := range albumTypes {
		albumType := albumType
		go func () {
			offset := 0
			for {
				options := &spotify.Options{Limit: &limit, Offset: &offset}
				albumTypes := spotify.AlbumType(albumType)
				results, _ := client.GetArtistAlbumsOpt(spotify.ID(artistId), options, &albumTypes)
				for _, album := range results.Albums{
					allAlbums[album.ID] = album
				}
				if len(results.Albums) == 50 {
					offset += 50
				} else {
					break
				}
			}
			done <- true
		}()
	}
	for range albumTypes{
		<- done
	}

	uniqueAlbums := make(map[spotify.ID]spotify.SimpleAlbum)
	for id, album := range allAlbums {
		if _, ok := uniqueAlbums[id]; !ok {
			uniqueAlbums[id] = album
		}
	}

	count := 0
	mutex := &sync.Mutex{}
	rate := time.Millisecond * 100
	throttle := time.Tick(rate)
	allTracks := make(map[spotify.ID]spotify.SimpleTrack)
	for _, album := range uniqueAlbums {
		<- throttle
		album := album
		go func () {
				results, err := client.GetAlbumTracks(album.ID)
				if err != nil {
					log.Fatal(err)
				}
				for _, track := range results.Tracks{
					for _, artist := range track.Artists {
						if artist.ID == spotify.ID(artistId) {
							count ++
							fmt.Println(count)
							fmt.Println(track.Name)
							mutex.Lock()
							allTracks[track.ID] = track
							mutex.Unlock()
						}
					}
				}
			done <- true
		}()
	}
	for range uniqueAlbums{
		<- done
	}

	var returnTracks []spotify.SimpleTrack

	for _, track := range allTracks {
		returnTracks = append(returnTracks, track)
	}
	return returnTracks
}

func SearchForArtist(artistName string) []spotify.FullArtist {
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

	if err != nil {
		log.Fatal(err)
	}

	var artistsArray []spotify.FullArtist
	if result.Artists != nil {
		for _, item := range result.Artists.Artists{
			artistsArray = append(artistsArray, item)
		}
	}

	return artistsArray
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
