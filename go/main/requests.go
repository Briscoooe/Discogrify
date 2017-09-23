package main

import (
	"context"
	"encoding/json"
	"github.com/Briscooe/Discogrify/go/caching"
	"github.com/Briscooe/Discogrify/go/logging"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

func AddContext(next http.Handler, l logging.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.Log("New request from IP: " + r.RemoteAddr)
		cookie, _ := r.Cookie(authTokenCookie)
		if cookie != nil {
			ctx := context.WithValue(r.Context(), "AuthToken", cookie.Value)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func LoginToSpotifyHandlerFunc(s *Spotify, expiration int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state := GenerateStateString()
		url := GenerateLoginUrl(s, state)
		cookie := http.Cookie{Name: spotifyAuthStateCookie, Value: state, Expires: time.Now().Add(time.Duration(expiration) * time.Minute)}
		http.SetCookie(w, &cookie)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write([]byte(url))
	})
}

func CallbackHandler(log logging.Logger, s *Spotify, expiration int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var status int
		cookie, _ := r.Cookie(spotifyAuthStateCookie)
		if cookie != nil {
			tok, err := ValidateCallback(r, log, s, cookie.Value)
			if err != nil {
				status = http.StatusBadRequest
				w.Write([]byte(err.Error()))
			} else {
				cookie := http.Cookie{Name: authTokenCookie, Value: tok.AccessToken, Expires: time.Now().Add(time.Duration(expiration) * time.Hour)}
				http.SetCookie(w, &cookie)
				status = http.StatusFound
			}
		} else {
			status = http.StatusBadRequest
		}
		http.Redirect(w, r, "/", status)
	})
}

func UserInfoHandler(l logging.Logger, s *Spotify) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tok := r.Context().Value("AuthToken"); tok != nil {
			user := GetUserInfo(l, s.NewClient(tok.(string)))
			if user != nil {
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(user); err != nil {
					l.LogErr(err)
				}
			} else {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Could not get user"))
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Not logged in"))
		}
	})
}

func GetTracksHandler(c caching.Client, l logging.Logger, s *Spotify) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tok := r.Context().Value("AuthToken"); tok != nil {
			id := mux.Vars(r)["artistId"]
			match, _ := regexp.MatchString("^[a-zA-Z0-9]{22,}$", id)
			if !match {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Invalid artist ID"))
			} else {
				IncrementKeyInCache(id, c)
				l.Logf("%s: Checking cache for artist ID", id)
				tracks := GetTracksFromCache(id, c, l)
				if tracks == nil {
					tracks = GetDiscography(id, l, s.NewClient(tok.(string)))
					tracksJson, _ := json.Marshal(tracks)
					AddToCache(id, string(tracksJson), c, l, formatArtistTracks)
				}
				l.Logf("%s: Returning tracks", id)
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(tracks); err != nil {
					l.LogErr(err)
				}
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Not logged in"))
		}
	})
}

func SearchArtistHandler(c caching.Client, l logging.Logger, s *Spotify) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tok := r.Context().Value("AuthToken"); tok != nil {
			query := mux.Vars(r)["name"]
			l.Logf("%s: Checking cache for search query ", query)
			results := GetSearchResultsFromCache(query, c, l)

			if len(results) == 0 {
				results = SearchForArtist(query, c, s.NewClient(tok.(string)), l)
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(results); err != nil {
				l.LogErr(err)
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Not logged in"))
		}
	})
}

func PublishPlaylistHandler(log logging.Logger, s *Spotify) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var message string
		if tok := r.Context().Value("AuthToken"); tok != nil {
			body, _ := ioutil.ReadAll(r.Body)
			type playlist struct {
				Tracks []string
				Title  string `json:"name"`
			}

			var newPlaylist playlist
			err := json.Unmarshal(body, &newPlaylist)
			if err != nil {
				log.Log(err)
			}

			url, status := PublishPlaylist(newPlaylist.Tracks, newPlaylist.Title, log, s.NewClient(tok.(string)))
			message = url
			w.WriteHeader(status)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			message = "Not logged in"
		}
		w.Write([]byte(message))
	})
}
