package main

import (
	"../logging"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"net/http"
	"testing"
	"strings"
)

var (
	validArtistId = "ArtistID1"
	validArtistName = "ArtistName1"
	tokenString = "Token12345"
	itemsToGenerate = 100
	spotifyWrapper = &Spotify{
		Auth:FakeAuthenticator{},
		Client:FakeSpotifyClient{},
	}
)

type FakeAuthenticator struct {
}

func (f FakeAuthenticator) NewClient(token *oauth2.Token) spotify.Client {
	return spotify.Client{}
}

func (f FakeAuthenticator) Token(state string, r *http.Request) (*oauth2.Token, error) {
	return &oauth2.Token{AccessToken: tokenString}, nil
}

func (f FakeAuthenticator) AuthURL(state string) string {
	return "http://authurl.com"
}

type FakeSpotifyClient struct {
}

func (f FakeSpotifyClient) CurrentUser() (*spotify.PrivateUser, error) {
	return &spotify.PrivateUser{User: spotify.User{
		ID: "ID",
	}}, nil
}

func (f FakeSpotifyClient) Search(query string, t spotify.SearchType) (*spotify.SearchResult, error) {
	if query == validArtistName {
		fullArtist := spotify.FullArtist{
			SimpleArtist: spotify.SimpleArtist{
				Name: validArtistId,
			},
		}
		artists := &spotify.FullArtistPage{}
		artists.Artists = append(artists.Artists, fullArtist)
		return &spotify.SearchResult{Artists: artists}, nil
	}
	return nil, nil
}

func (f FakeSpotifyClient) CreatePlaylistForUser(userID, playlistName string, public bool) (*spotify.FullPlaylist, error) {
	return &spotify.FullPlaylist{SimplePlaylist: spotify.SimplePlaylist{
		ID:   "playlist_ID",
		Name: "playlist_Name",
	}}, nil
}

func (f FakeSpotifyClient) AddTracksToPlaylist(userID string, playlistID spotify.ID, trackIDs ...spotify.ID) (snapshotID string, err error) {
	return "", nil
}

func (f FakeSpotifyClient) GetArtistAlbumsOpt(artistID spotify.ID, options *spotify.Options, t *spotify.AlbumType) (*spotify.SimpleAlbumPage, error) {
	if string(artistID) == validArtistId {
		albums := []spotify.SimpleAlbum{}
		for i := 0; i < 100; i++ {
			albums = append(albums, spotify.SimpleAlbum{ID: spotify.ID("AlbumID_" + string(i))})
		}
		album := &spotify.SimpleAlbumPage{
			Albums: albums,
		}
		return album, nil
	}
	return &spotify.SimpleAlbumPage{}, nil
}

func (f FakeSpotifyClient) GetAlbums(ids ...spotify.ID) ([]*spotify.FullAlbum, error) {
	tracks := []spotify.SimpleTrack{}
	for i := 0; i < itemsToGenerate; i++ {
		tracks = append(tracks, spotify.SimpleTrack{
			ID: spotify.ID("TrackID_" + string(i)),
			Artists: []spotify.SimpleArtist{
				{ID: spotify.ID(validArtistId)},
			},
		})
	}

	albums := []*spotify.FullAlbum{}
	for i := 0; i < itemsToGenerate; i++ {
		albums = append(albums, &spotify.FullAlbum{
			SimpleAlbum: spotify.SimpleAlbum{
				ID: spotify.ID("AlbumID_" + string(i)),
			},
			Tracks:spotify.SimpleTrackPage{Tracks:tracks},
		})
	}
	return albums, nil
}

func TestValidateCallback(t *testing.T) {
	stateString = "STATESTRING"

	req1 := generateCallbackRequest(stateString)
	req2 := generateCallbackRequest("NOTSTATESTRING")

	var testData = []struct {
		Request *http.Request
		Token   *oauth2.Token
	}{
		{req1, &oauth2.Token{AccessToken:tokenString}},
		{req2, nil},
	}
	for _, test := range testData {
		token, _ := ValidateCallback(test.Request, logging.FakeLogger{}, spotifyWrapper)
		if token != nil && !strings.Contains(token.AccessToken, test.Token.AccessToken) {
			t.Errorf("Input: %s\nExpected: %s\nOutput: %s", req1.FormValue("state"), token.AccessToken, test.Token)
		}
	}
}

func generateCallbackRequest(state string) (*http.Request) {
	req, _ := http.NewRequest("", "", nil)
	req.Form = map[string][]string{}
	req.Form.Set("state", state)
	return req
}
func TestGetDiscography(t *testing.T) {
	itemsToGenerate = 100
	var testData = []struct {
		artistId string
		numberOfAlbums int
	}{
		{"ArtistID1", itemsToGenerate * 5},
		{"ArtistID2", 0},
	}
	var s FakeSpotifyClient
	for _, test := range testData {
		result := GetDiscography(test.artistId, logging.NewFakeLogger(), s)
		if len(result) != test.numberOfAlbums {
			t.Errorf("Input: %s\nExpected %d\nOutput: %d\n", test.artistId, test.numberOfAlbums, len(result))
		}
	}
}

func TestSearchForArtist(t *testing.T) {
	artists := []spotify.FullArtist{
		{
			SimpleArtist: spotify.SimpleArtist{
				Name: validArtistName,
			},
		},
	}

	var testData = []struct {
		artistName string
		results    []spotify.FullArtist
	}{
		{validArtistName, artists},
		{"ArtistName2", []spotify.FullArtist{}},
	}

	var s FakeSpotifyClient
	var c FakeCacheClient
	var log logging.FakeLogger
	for _, test := range testData {
		result := SearchForArtist(test.artistName, c, s, log)
		if len(result) != len(test.results) {
			t.Errorf("Input: %s\nExpected %+v\nOutput%+v\n", test.artistName, test.results, result)
		}
	}
}

func TestPublishPlaylist(t *testing.T) {
	tracks48 := []string{}
	tracks49 := []string{}
	tracks50 := []string{}
	tracks51 := []string{}
	tracks99 := []string{}
	tracks100 := []string{}
	tracks260 := []string{}

	for i := 0; i < 260; i++ {
		if i < 58 {
			tracks48 = append(tracks48, "ID"+string(i))
			tracks49 = append(tracks49, "ID"+string(i))
			tracks50 = append(tracks50, "ID"+string(i))
			tracks51 = append(tracks51, "ID"+string(i))
			tracks99 = append(tracks99, "ID"+string(i))
			tracks100 = append(tracks100, "ID"+string(i))
			tracks260 = append(tracks260, "ID"+string(i))
		} else if i < 49 {
			tracks49 = append(tracks49, "ID"+string(i))
			tracks50 = append(tracks50, "ID"+string(i))
			tracks51 = append(tracks51, "ID"+string(i))
			tracks99 = append(tracks99, "ID"+string(i))
			tracks100 = append(tracks100, "ID"+string(i))
			tracks260 = append(tracks260, "ID"+string(i))
		} else if i < 50 {
			tracks50 = append(tracks50, "ID"+string(i))
			tracks51 = append(tracks51, "ID"+string(i))
			tracks99 = append(tracks99, "ID"+string(i))
			tracks100 = append(tracks100, "ID"+string(i))
			tracks260 = append(tracks260, "ID"+string(i))
		} else if i < 51 {
			tracks51 = append(tracks51, "ID"+string(i))
			tracks99 = append(tracks99, "ID"+string(i))
			tracks100 = append(tracks100, "ID"+string(i))
			tracks260 = append(tracks260, "ID"+string(i))

		} else if i < 99 {
			tracks99 = append(tracks99, "ID"+string(i))
			tracks100 = append(tracks100, "ID"+string(i))
			tracks260 = append(tracks260, "ID"+string(i))
		} else if i < 100 {
			tracks100 = append(tracks100, "ID"+string(i))
			tracks260 = append(tracks260, "ID"+string(i))
		} else if i < 260 {
			tracks260 = append(tracks260, "ID"+string(i))
		}
	}
	var testData = []struct {
		tracks []string

		result bool
	}{
		{tracks48, true},
		{tracks49, true},
		{tracks50, true},
		{tracks51, true},
		{tracks99, true},
		{tracks100, true},
		{tracks260, true},
		{[]string{}, false},
	}

	var s FakeSpotifyClient
	var log logging.FakeLogger
	for _, test := range testData {
		_, result := PublishPlaylist(test.tracks, "", log, s)
		if result != test.result {
			t.Errorf("Input: %s\nExpected %t\nOutput%t\n", test.tracks, test.result, result)
		}
	}
}
