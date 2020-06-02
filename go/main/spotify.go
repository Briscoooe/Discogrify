package main

import (
	"encoding/json"
	"errors"
	"../caching"
	"../logging"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"math/rand"
	"net/http"
	"fmt"
)

var (
	characters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
)

type SpotifyClient interface {
	CurrentUser() (*spotify.PrivateUser, error)
	Search(query string, t spotify.SearchType) (*spotify.SearchResult, error)
	CreatePlaylistForUser(, playlistName string, public bool) (*spotify.FullPlaylist, error)
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

func GetUserInfo(l logging.Logger, s SpotifyClient) *spotify.PrivateUser {
	user, err := s.CurrentUser()
	if err != nil {
		l.Log("Could not get user info")
		l.Log(err)
		return nil
	}
	return user
}

func GenerateStateString() string {
	b := make([]rune, 16)
	for i := range b {
		b[i] = characters[rand.Intn(len(characters))]
	}
	return string(b)
}

func GenerateLoginUrl(s *Spotify, state string) string {
	url := s.Auth.AuthURL(state)
	return url
}

func ValidateCallback(r *http.Request, l logging.Logger, s *Spotify, state string) (token *oauth2.Token, err error) {
	tok, err := s.Auth.Token(state, r)
	if err != nil {
		l.LogErr(err,"Could not get token")
		return nil, err
	}

	if st := r.FormValue("state"); st != state {
		l.Logf("State mismatch: %s != %s\n", st, state)
		return nil, errors.New(fmt.Sprintf("State mismatch: %s != %s\n", st, state))
	}

	return tok, nil
}

func GetDiscography(id string, l logging.Logger, s SpotifyClient) []*spotify.FullAlbum {
	var allTracks []*spotify.FullAlbum
	l.Logf("%s: Getting unique albums...", id)
	uniqueAlbums := getUniqueAlbums(id, s, l)
	if len(uniqueAlbums) > 0 {
		l.Logf("%s: Unique albums found: %v\n", id, len(uniqueAlbums))
		l.Logf("%s: Getting unique tracks...", id)
		allTracks = getAllTracksFromAlbums(id, uniqueAlbums, s, l)
		l.Logf("%s: Unique tracks found: %v\n", id, len(allTracks))
	}

	return allTracks
}

func SearchForArtist(name string, c caching.Client, s SpotifyClient, l logging.Logger) []spotify.FullArtist {
	l.Log("Searching for artist:", name)
	result, err := s.Search(name, spotify.SearchTypeArtist)

	if err != nil {
		l.LogErrf(err,"%s: Could not get results", name)
	}

	var artistsArray []spotify.FullArtist
	if result != nil && len(result.Artists.Artists) > 0 {
		for _, item := range result.Artists.Artists {
			artistsArray = append(artistsArray, item)
		}
		artistsJson, _ := json.Marshal(artistsArray)
		AddToCache(name, string(artistsJson), c, l, formatSearchArtist)
		l.Logf("Artists found: %v\n", len(artistsArray))
	}

	return artistsArray
}

func PublishPlaylist(tracks []string, name string, l logging.Logger, s SpotifyClient) (string, int) {
	user, err := s.CurrentUser()
	if err != nil {
		l.Log("Could not get user")
		return "", http.StatusUnauthorized
	}

	if len(tracks) == 0 {
		l.Log("No tracks present")
		return "", http.StatusNotFound
	}

	playlist, err := s.CreatePlaylistForUser(user.ID, name+" - By Discogrify", true)

	if err != nil {
		l.LogErr(err, "Could not create playlist")
		return "", http.StatusNotModified
	}

	l.Log("Playlist created: " + playlist.ID)

	var ids []spotify.ID
	for _, track := range tracks {
		ids = append(ids, spotify.ID(track))
	}

	tracksPerRequest := 50
	startIndex := 0
	endIndex := tracksPerRequest -1
	added := 0
	if endIndex > len(tracks) {
		endIndex = len(tracks)
	}
	status := http.StatusOK
	for added != len(tracks) {
		if added == 0 && len(tracks) == endIndex {
			endIndex -= 1
			added += 1
		}
		_, err := s.AddTracksToPlaylist(user.ID, playlist.ID, ids[startIndex:endIndex+1]...)
		if err != nil {
			if err.Error() != "Invalid track uri: spotify:track:" {
				l.LogErr(err, "Error adding tracks to playlist")
				l.Log(ids[startIndex:endIndex+1])
				status = http.StatusBadRequest
			}
		}
		added += endIndex - startIndex
		if endIndex >= tracksPerRequest - 1 && endIndex != len(tracks) {
			added += 1
		}
		startIndex += tracksPerRequest
		if len(tracks)-added < tracksPerRequest {
			endIndex = len(tracks) - 1
		} else {
			endIndex += tracksPerRequest
		}
	}
	l.Logf("Added %d tracks to playlist ID: %s", added, playlist.ID)
	return string(playlist.ExternalURLs["spotify"]), status
}

func getUniqueAlbums(id string, s SpotifyClient, l logging.Logger) []spotify.ID {
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
					l.LogErrf(err, "%s: Could not get albums", id)
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

func getAllTracksFromAlbums(id string, albums []spotify.ID, s SpotifyClient, l logging.Logger) []*spotify.FullAlbum {
	count := 0
	var allTracks []*spotify.FullAlbum
	for i := 0; i < len(albums); i += 20 {
		i := i
		go func() {
			limit := i + 20
			if limit >= len(albums) {
				limit = i + (i-len(albums))*-1
			}
			results, err := s.GetAlbums(albums[i:limit]...)
			if err != nil {
				l.LogErrf(err, "%s: Could not get albums", id)
			}
			for _, album := range results {
				var tempTrackList []spotify.SimpleTrack
				for _, track := range album.Tracks.Tracks {
					track := track
					for _, artist := range track.Artists {
						if artist.ID == spotify.ID(id) {
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
	for i := 0; i < len(albums); i += 20 {
		<-done
	}
	return allTracks
}