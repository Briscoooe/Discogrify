package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
	"os"
	"context"
	"sync"
//	"time"
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

func GetAllArtistTracks(artistId string) []spotify.SimpleTrack {
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

	mutex := &sync.Mutex{}

	done := make(chan bool)
	allAlbums := make(map[spotify.ID]spotify.SimpleAlbum)
	limit := 50
	albumTypes := []int{1,2,4,5}
	for _, albumType := range albumTypes {
		albumType := albumType
		go func () {
			offset := 0
			for {
				options := &spotify.Options{Limit: &limit, Offset: &offset}
				albumTypes := spotify.AlbumType(albumType)
				results, _ := client.GetArtistAlbumsOpt(spotify.ID(artistId), options, &albumTypes)
				for _, album := range results.Albums{
					mutex.Lock()
					allAlbums[album.ID] = album
					mutex.Unlock()
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

	var uniqueAlbumsArray []spotify.ID

	for _, track := range uniqueAlbums {
		uniqueAlbumsArray = append(uniqueAlbumsArray, track.ID)
	}

	var allTracks []spotify.SimpleTrack
	for i := 0; i < len(uniqueAlbumsArray); i += 20 {
		i := i
		go func () {
			limit := i + 20
			if limit >= len(uniqueAlbumsArray) {
				limit = i + (i - len(uniqueAlbumsArray)) * -1
			}
			results, err := client.GetAlbums(uniqueAlbumsArray[i:limit]...)
			if err != nil {
				log.Fatal(err)
			}
			for _, album := range results{
				for _, track := range album.Tracks.Tracks{
					for _, artist := range track.Artists {
						if artist.ID == spotify.ID(artistId) {
							mutex.Lock()
							allTracks = append(allTracks, track)
							mutex.Unlock()
						}
					}
				}
			}
			done <- true
		}()
	}
	// Possibly have to change
	for i := 0; i < len(uniqueAlbumsArray); i += 20 {
		<- done
	}

	return allTracks
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
