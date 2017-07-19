package main

import (
	"fmt"
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
	done		 = make(chan bool)
	mutex 		 = &sync.Mutex{}
	token 	     = ""
	ch           = make(chan *spotify.Client)
	client       spotify.Client
	auth         = spotify.NewAuthenticator(redirectURI, spotify.ScopePlaylistModifyPrivate, spotify.ScopePlaylistModifyPublic)
	unauthClient = spotify.DefaultClient
	characters   = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
)

func GetDiscographyFromSpotify(artistId string) []*spotify.FullAlbum {
	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		TokenURL:     spotify.TokenURL,
	}
	token, err := config.Token(context.Background())
	if err != nil {
		rollingLog.Fatalf("couldn't get token: %v", err)
	}

	client = spotify.Authenticator{}.NewClient(token)
	mutex = &sync.Mutex{}

	uniqueAlbums := getUniqueAlbums(artistId)

	allTracks := getAllTracksFromAlbums(artistId, uniqueAlbums)

	return allTracks
}

func getUniqueAlbums(artistId string) []spotify.ID {
	allAlbums := make(map[spotify.ID]spotify.SimpleAlbum)
	limit := 50
	albumTypes := []int{1,2,3,4,5}
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
	rollingLog.Printf("%s: Unique albums found: %v\n", artistId, len(uniqueAlbums))

	return uniqueAlbumsArray
}

func getAllTracksFromAlbums(artistId string, uniqueAlbums []spotify.ID) []*spotify.FullAlbum {
	count := 0
	var allTracks []*spotify.FullAlbum
	for i := 0; i < len(uniqueAlbums); i += 20 {
		i := i
		go func () {
			limit := i + 20
			if limit >= len(uniqueAlbums) {
				limit = i + (i - len(uniqueAlbums)) * -1
			}
			results, err := client.GetAlbums(uniqueAlbums[i:limit]...)
			if err != nil {
				rollingLog.Fatal(err)
			}
			for _, album := range results{
				var tempTrackList []spotify.SimpleTrack
				for _, track := range album.Tracks.Tracks{
					track := track
					for _, artist := range track.Artists {
						if artist.ID == spotify.ID(artistId) {
							mutex.Lock()
							tempTrackList = append(tempTrackList, track)
							count ++
							mutex.Unlock()
						}
					}
				}
				album.Tracks = spotify.SimpleTrackPage{
					Tracks: tempTrackList,
				}
				allTracks = append(allTracks, album)
			}
			done <- true
		}()
	}
	for i := 0; i < len(uniqueAlbums); i += 20 {
		<- done
	}
	rollingLog.Printf("%s: Tracks found: %v\n", artistId, count)

	return allTracks
}

type CustomAlbum struct {
	AlbumName	string `json:"name"`
	AlbumId		string `json:"id"`
	Tracks 		[]spotify.SimpleTrack `json:"tracks"`
}

type CustomAlbum1 struct {
	Tracks map[string][]spotify.SimpleTrack
}

func SearchForArtist(artistName string) []spotify.FullArtist {
	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		TokenURL:     spotify.TokenURL,
	}
	token, err := config.Token(context.Background())
	if err != nil {
		rollingLog.Fatalf("couldn't get token: %v", err)
	}

	client := spotify.Authenticator{}.NewClient(token)

	rollingLog.Printf("Searching for arist: %s\n", artistName)
	result, err := client.Search(artistName, spotify.SearchTypeArtist)

	if err != nil {
		rollingLog.Fatal(err)
	}

	var artistsArray []spotify.FullArtist
	if result.Artists != nil {
		for _, item := range result.Artists.Artists{
			artistsArray = append(artistsArray, item)
		}
	}

	rollingLog.Printf("Artists found: %v\n", len(artistsArray))
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
		rollingLog.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		rollingLog.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	// use the token to get an authenticated client
	client := auth.NewClient(tok)
	fmt.Println(w, "Login Completed!")

	client.CurrentUser()
}