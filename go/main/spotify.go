package main

import (
	"encoding/json"
	"errors"
	"github.com/Briscooe/Discogrify/go/caching"
	"github.com/Briscooe/Discogrify/go/logging"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"math/rand"
	"net/http"
)

var (
	characters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
)

type SpotifyClient interface {
	CurrentUser() (*spotify.PrivateUser, error)
	Search(query string, t spotify.SearchType) (*spotify.SearchResult, error)
	CreatePlaylistForUser(userID, playlistName string, public bool) (*spotify.FullPlaylist, error)
	AddTracksToPlaylist(userID string, playlistID spotify.ID, trackIDs ...spotify.ID) (snapshotID string, err error)
	GetArtistAlbumsOpt(artistID spotify.ID, options *spotify.Options, t *spotify.AlbumType) (*spotify.SimpleAlbumPage, error)
	GetAlbums(ids ...spotify.ID) ([]*spotify.FullAlbum, error)
}

type Authenticator interface {
	NewClient(token *oauth2.Token) spotify.Client
	Token(state string, r *http.Request) (*oauth2.Token, error)
	AuthURL(state string) string
}

type Spotify struct {
	Client SpotifyClient
	Auth   Authenticator
}

func InitSpotifyClient(redirectUri string) *Spotify {
	return &Spotify{
		Auth: spotify.NewAuthenticator(redirectUri, spotify.ScopePlaylistModifyPublic),
	}
}

func (s Spotify) NewClient(tokenStr string) SpotifyClient {
	token := oauth2.Token{AccessToken: tokenStr}
	client := s.Auth.NewClient(&token)
	client.AutoRetry = true
	return &client
}

func GenerateStateString() string {
	b := make([]rune, 16)
	for i := range b {
		b[i] = characters[rand.Intn(len(characters))]
	}
	return string(b)
}

func GenerateLoginUrl(s *Spotify) string {
	stateString = GenerateStateString()
	url := s.Auth.AuthURL(stateString)
	return url
}

func ValidateCallback(r *http.Request, log logging.Logger, s *Spotify) (token *oauth2.Token, err error) {
	tok, err := s.Auth.Token(stateString, r)
	if err != nil {
		log.Fatal(err)
		log.Fatal("Could not get token")
		return nil, err
	}

	if st := r.FormValue("state"); st != stateString {
		log.Printf("State mismatch: %s != %s\n", st, stateString)
		return nil, errors.New("State mismatch")
	}

	return tok, nil
}

func GetDiscography(id string, log logging.Logger, s SpotifyClient) []*spotify.FullAlbum {
	var allTracks []*spotify.FullAlbum
	log.Printf("%s: Getting unique albums...", id)
	uniqueAlbums := getUniqueAlbums(id, s, log)
	if len(uniqueAlbums) > 0 {
		log.Printf("%s: Unique albums found: %v\n", id, len(uniqueAlbums))
		log.Printf("%s: Getting unique tracks...", id)
		allTracks = getAllTracksFromAlbums(id, uniqueAlbums, s, log)
		log.Printf("%s: Unique tracks found: %v\n", id, len(allTracks))
	}

	return allTracks
}

func SearchForArtist(name string, c caching.Client, s SpotifyClient, log logging.Logger) []spotify.FullArtist {
	log.Println("Searching for artist:", name)
	result, err := s.Search(name, spotify.SearchTypeArtist)

	if err != nil {
		log.Printf("%s: Could not get results", name)
		log.Println(err)
	}

	var artistsArray []spotify.FullArtist
	if result != nil {
		for _, item := range result.Artists.Artists {
			artistsArray = append(artistsArray, item)
		}
	}

	artistsJson, _ := json.Marshal(artistsArray)
	AddToCache(name, string(artistsJson), c, log, formatSearchArtist)
	log.Printf("Artists found: %v\n", len(artistsArray))
	return artistsArray
}

func PublishPlaylist(tracks []string, name string, log logging.Logger, s SpotifyClient) (string, bool) {
	result := true

	user, err := s.CurrentUser()
	if err != nil {
		log.Println("Could not get user")
		result = false
	}

	if len(tracks) == 0 {
		log.Println("No tracks present")
		return "No tracks present", false
	}

	playlist, err := s.CreatePlaylistForUser(user.ID, name+" - By Discogrify", true)

	if err != nil {
		log.Printf("Could not create playlist")
		log.Println(err)
		result = false
	}

	log.Println("Playlist created: " + playlist.ID)

	var ids []spotify.ID
	for _, track := range tracks {
		ids = append(ids, spotify.ID(track))
	}

	startIndex := 0
	endIndex := 49
	added := 0
	if endIndex > len(tracks) {
		endIndex = len(tracks)
	}
	for added != len(tracks) {
		_, err := s.AddTracksToPlaylist(user.ID, playlist.ID, ids[startIndex:endIndex+1]...)
		if err != nil {
			log.Fatal("Error adding tracks to playlist")
			log.Fatal(err)
			result = false
		}
		added += endIndex - startIndex
		if endIndex >= 49 && endIndex != len(tracks) {
			added += 1
		}
		startIndex += 50
		if len(tracks)-added < 50 {
			endIndex = len(tracks) - 1
		} else {
			endIndex += 50
		}
	}

	log.Printf("Added %d tracks to playlist ID: %s", added, playlist.ID)
	return string(playlist.ExternalURLs["spotify"]), result
}

func getUniqueAlbums(id string, s SpotifyClient, log logging.Logger) []spotify.ID {
	allAlbums := make(map[spotify.ID]spotify.SimpleAlbum)
	limit := 50
	albumTypes := []int{1, 2, 3, 4, 5}
	for _, albumType := range albumTypes {
		albumType := albumType
		go func() {
			offset := 0
			for {
				options := &spotify.Options{Limit: &limit, Offset: &offset}
				albumTypes := spotify.AlbumType(albumType)
				results, err := s.GetArtistAlbumsOpt(spotify.ID(id), options, &albumTypes)
				if err != nil {
					log.Printf("%s: Could not get albums", id)
					log.Println(err)
					break
				}
				for _, album := range results.Albums {
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
	for range albumTypes {
		<-done
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
	return uniqueAlbumsArray
}

func getAllTracksFromAlbums(artistId string, uniqueAlbums []spotify.ID, s SpotifyClient, logger logging.Logger) []*spotify.FullAlbum {
	count := 0
	var allTracks []*spotify.FullAlbum
	for i := 0; i < len(uniqueAlbums); i += 20 {
		i := i
		go func() {
			limit := i + 20
			if limit >= len(uniqueAlbums) {
				limit = i + (i-len(uniqueAlbums))*-1
			}
			results, err := s.GetAlbums(uniqueAlbums[i:limit]...)
			if err != nil {
				logger.Printf("%s: Could not get albums", artistId)
				logger.Println(err)
			}
			for _, album := range results {
				var tempTrackList []spotify.SimpleTrack
				for _, track := range album.Tracks.Tracks {
					track := track
					for _, artist := range track.Artists {
						if artist.ID == spotify.ID(artistId) {
							mutex.Lock()
							tempTrackList = append(tempTrackList, track)
							count++
							mutex.Unlock()
						}
					}
				}
				if len(tempTrackList) > 0 {
					album.Tracks = spotify.SimpleTrackPage{
						Tracks: tempTrackList,
					}
					allTracks = append(allTracks, album)
				}
			}
			done <- true
		}()
	}
	for i := 0; i < len(uniqueAlbums); i += 20 {
		<-done
	}
	return allTracks
}
