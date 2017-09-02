package main

import (
	"fmt"
	"github.com/zmb3/spotify"
	"net/http"
	//	"time"
	"github.com/Briscooe/Discogrify/go/logging"
	"github.com/Briscooe/Discogrify/go/caching"
	"encoding/json"
	"bytes"
	"io/ioutil"
	"time"
	"golang.org/x/oauth2"
)

type Spotify interface {
	GetCurrentUser() (spotify.User, error)
	GenerateLoginUrl() map[string]string
	ValidateCallback(r *http.Request) (*oauth2.Token, error, string)
	GetDiscography(artistId string) []*spotify.FullAlbum
	SearchForArtist(artistName string, cacheClient caching.Client) []spotify.FullArtist
	PublishPlaylist(tracks []string) bool
}

type SpotifyClient struct {
	Client spotify.Client
	StateString string
	Authenticator spotify.Authenticator
	Logger logging.Logger
}

func NewSpotifyClient(logger logging.Logger) *SpotifyClient {
	auth := spotify.NewAuthenticator(redirectURI, spotify.ScopePlaylistModifyPrivate, spotify.ScopePlaylistModifyPublic)
	return &SpotifyClient{
		Authenticator: auth,
		Logger: logger,
	}
}

func (s *SpotifyClient) GetCurrentUser() (spotify.User, error) {
	user, err := s.Client.CurrentUser()
	if err != nil {
		return user.User, err
	}
	return user.User, nil
}
func (s *SpotifyClient) GenerateLoginUrl() map[string]string {
	s.StateString = s.Authenticator.AuthURL(GenerateStateString())

	url := s.Authenticator.AuthURL(s.StateString)

	urlJson := map[string]string{"url": url}

	return urlJson
}

func (s *SpotifyClient) ValidateCallback(r *http.Request) (token *oauth2.Token, err error, errMsg string) {
	tok, err := s.Authenticator.Token(s.StateString, r)
	if err != nil {
		s.Logger.Println(err)
		return tok, err, "Could not get token"
	}
	if st := r.FormValue("state"); st != s.StateString {
		return tok, err, fmt.Sprintf("State mismatch: %s != %s\n", st, stateString)
	}

	s.Client = s.Authenticator.NewClient(tok)
	s.Client.AutoRetry = true

	return tok, nil, ""
}

func (s *SpotifyClient) GetDiscography(artistId string) []*spotify.FullAlbum {
	var allTracks []*spotify.FullAlbum
	s.Logger.Printf("%s: Getting unique albums...", artistId)
	uniqueAlbums := s.getUniqueAlbums(artistId)
	if len(uniqueAlbums) > 0 {
		s.Logger.Printf("%s: Unique albums found: %v\n", artistId, len(uniqueAlbums))
		s.Logger.Printf("%s: Getting unique tracks...", artistId)
		allTracks = s.getAllTracksFromAlbums(artistId, uniqueAlbums)
		s.Logger.Printf("%s: Unique tracks found: %v\n", artistId, len(allTracks))
	}

	return allTracks
}

func (s *SpotifyClient) SearchForArtist(artistName string, cacheClient caching.Client) []spotify.FullArtist {
	s.Logger.Println("Searching for artist:", artistName)
	result, err := s.Client.Search(artistName, spotify.SearchTypeArtist)

	if err != nil {
		s.Logger.Printf("%s: Could not get results", artistName)
		s.Logger.Println(err)
	}

	var artistsArray []spotify.FullArtist
	if result.Artists != nil {
		for _, item := range result.Artists.Artists {
			artistsArray = append(artistsArray, item)
		}
	}

	artistsJson, _ := json.Marshal(artistsArray)
	if cacheClient.Set("artist:search:" + artistName , string(artistsJson), time.Hour * 168) {
		s.Logger.Printf("Added query to cache: " + artistName)
	}
	s.Logger.Printf("Artists found: %v\n", len(artistsArray))
	return artistsArray
}

func (s *SpotifyClient) PublishPlaylist(tracks []string) bool {
	user, _ := s.Client.CurrentUser()
	playlist, err := s.Client.CreatePlaylistForUser(user.ID, "My Playlist", true)

	if err != nil {
		s.Logger.Printf("Could not create playlist")
		s.Logger.Println(err)
		return false
	}

	s.Logger.Println("Playlist created: " + playlist.ID)

	uris := make([]string, len(tracks))
	for i, id := range tracks {
		uris[i] = fmt.Sprintf("spotify:track:%s", id)
	}

	var bodyStruct = struct {
		Uris []string `json:"uris"`
	} {
		Uris:uris,
	}

	tok, err := s.Client.Token()
	client := &http.Client{}
	spotifyUrl := fmt.Sprintf("https://api.spotify.com/v1/users/%s/playlists/%s/tracks", user.ID, string(playlist.ID))
	startIndex := 0
	endIndex := 99
	tracksAdded := 0
	if endIndex > len(tracks) {
		endIndex = len(tracks)
	}
	result := true
	for tracksAdded != len(tracks) {
		bodyJSON, err := json.Marshal(bodyStruct.Uris[startIndex:endIndex+1])
		if err != nil {
			s.Logger.Println(err)
		}
		tracksAdded += endIndex - startIndex + 1
		// If less than 100 tracks left
		if len(tracks) - tracksAdded < 100 {
			startIndex += 100
			endIndex = len(tracks) - 1
		} else {
			startIndex += 100
			endIndex += 100
		}
		req, err := http.NewRequest("POST", spotifyUrl, bytes.NewReader(bodyJSON))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer " + string(tok.AccessToken))

		resp, err := client.Do(req)

		if err != nil {
			s.Logger.Println("Request failed", err)
		}
		defer resp.Body.Close()

		if resp.Status != "201 Created" {
			s.Logger.Println("Response Status:", resp.Status)
			s.Logger.Println("Response Headers:", resp.Header)
			body, _ := ioutil.ReadAll(resp.Body)
			s.Logger.Println("Response Body:", string(body))
			result = false
		}
	}
	s.Logger.Printf("Added %s tracks to playlist ID: %s", tracksAdded, playlist.ID)
	return result
}

func (s *SpotifyClient) getUniqueAlbums(artistId string) []spotify.ID {
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
				results, err := s.Client.GetArtistAlbumsOpt(spotify.ID(artistId), options, &albumTypes)
				if err != nil {
					s.Logger.Printf("%s: Could not get albums", artistId)
					s.Logger.Println(err)
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

func (s *SpotifyClient) getAllTracksFromAlbums(artistId string, uniqueAlbums []spotify.ID) []*spotify.FullAlbum {
	count := 0
	var allTracks []*spotify.FullAlbum
	for i := 0; i < len(uniqueAlbums); i += 20 {
		i := i
		go func() {
			limit := i + 20
			if limit >= len(uniqueAlbums) {
				limit = i + (i-len(uniqueAlbums))*-1
			}
			results, err := s.Client.GetAlbums(uniqueAlbums[i:limit]...)
			if err != nil {
				s.Logger.Printf("%s: Could not get albums", artistId)
				s.Logger.Println(err)
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