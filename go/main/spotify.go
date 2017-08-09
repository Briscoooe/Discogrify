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
)

type Spotify interface {
	CheckLoggedIn() bool
	GenerateLoginUrl() map[string]string
	ValidateCallback(r *http.Request) (*spotify.PrivateUser, error, string)
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

func (s *SpotifyClient) CheckLoggedIn() bool {
	_, err := s.Client.CurrentUser()
	if err != nil {
		return false
	}
	return true
}
func (s *SpotifyClient) GenerateLoginUrl() map[string]string {
	s.StateString = s.Authenticator.AuthURL(GenerateStateString())

	url := s.Authenticator.AuthURL(s.StateString)

	urlJson := map[string]string{"url": url}

	return urlJson
}

func (s *SpotifyClient) ValidateCallback(r *http.Request) (user *spotify.PrivateUser, err error, errMsg string) {
	tok, err := s.Authenticator.Token(s.StateString, r)
	if err != nil {
		return user, err, "Could not get token"
		s.Logger.Fatal(err)
	}
	if st := r.FormValue("state"); st != s.StateString {
		return user, err, fmt.Sprintf("State mismatch: %s != %s\n", st, stateString)
	}

	s.Client = s.Authenticator.NewClient(tok)
	s.Client.AutoRetry = true
	user, err = s.Client.CurrentUser()
	if err != nil {
		return user, err, err.Error()
	}

	return user, nil, ""
}

func (s *SpotifyClient) GetDiscography(artistId string) []*spotify.FullAlbum {
	s.Logger.Printf("%s: Getting unique albums...", artistId)
	uniqueAlbums := s.getUniqueAlbums(artistId)
	s.Logger.Printf("%s: Unique albums found: %v\n", artistId, len(uniqueAlbums))
	s.Logger.Printf("%s: Getting unique tracks...", artistId)
	allTracks := s.getAllTracksFromAlbums(artistId, uniqueAlbums)
	s.Logger.Printf("%s: Unique tracks found: %v\n", artistId, len(allTracks))
	return allTracks
}

func (s *SpotifyClient) SearchForArtist(artistName string, cacheClient caching.Client) []spotify.FullArtist {
	s.Logger.Println("Searching for artist:", artistName)
	result, err := s.Client.Search(artistName, spotify.SearchTypeArtist)

	if err != nil {
		s.Logger.Fatal(err)
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
		s.Logger.Fatal("Could not create playlist", err)
	}

	uris := make([]string, len(tracks))
	for i, id := range tracks {
		uris[i] = fmt.Sprintf("spotify:track:%s", id)
	}

	var body = struct {
		Uris []string `json:uris`
	} {
		Uris:uris,
	}

	bodyJSON, err := json.Marshal(body)
	if err != nil {
		s.Logger.Fatal(err)
	}

	spotifyUrl := fmt.Sprintf("https://api.spotify.com/v1/users/%s/playlists/%s/tracks", user.ID, string(playlist.ID))
	req, err := http.NewRequest("POST", spotifyUrl, bytes.NewReader(bodyJSON))

	req.Header.Set("Content-Type", "application/json")
	tok, err := s.Client.Token()
	req.Header.Set("Authorization", "Bearer " + string(tok.AccessToken))
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		s.Logger.Fatal("Could not send re	quest", err)
	}

	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	bodyy, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(bodyy))
	return true
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
				results, _ := s.Client.GetArtistAlbumsOpt(spotify.ID(artistId), options, &albumTypes)
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
				s.Logger.Fatal(err)
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
				album.Tracks = spotify.SimpleTrackPage{
					Tracks: tempTrackList,
				}
				allTracks = append(allTracks, album)
			}
			done <- true
		}()
	}
	for i := 0; i < len(uniqueAlbums); i += 20 {
		<-done
	}
	return allTracks
}