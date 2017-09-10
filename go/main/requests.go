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

var (
	stateString string
)

func AddContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, _ := r.Cookie("auth_token")
		if cookie != nil {
			ctx := context.WithValue(r.Context(), "AuthToken", cookie.Value)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
func IndexHandlerFunc(l logging.Logger) http.Handler {
	return nil
}

func LoginToSpotifyHandlerFunc(s *Spotify) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := GenerateLoginUrl(s)
		w.Write([]byte(url))
	})
}

func GetTracksHandler(c caching.Client, log logging.Logger, s *Spotify) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tok := r.Context().Value("AuthToken"); tok != nil {
			id := mux.Vars(r)["artistId"]
			match, _ := regexp.MatchString("^[a-zA-Z0-9]{22,}$", id)
			if !match {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Invalid artist ID"))
			} else {
				IncrementKeyInCache(id, c)
				log.Printf("%s: Checking cache for artist ID", id)
				tracks := GetTracksFromCache(id, c, log)
				if tracks == nil {
					tracks = GetDiscography(id, log, s.NewClient(tok.(string)))
					tracksJson, _ := json.Marshal(tracks)
					AddToCache(id, string(tracksJson), c, log, formatArtistTracks)
				}
				log.Printf("%s: Returning tracks", id)
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(tracks); err != nil {
					panic(err)
				}
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Not logged in"))
		}
	})
}

func CallbackHandler(log logging.Logger, s *Spotify, cookieName string, expiration int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tok, err := ValidateCallback(r, log, s)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
		} else {
			cookie := http.Cookie{Name: cookieName, Value: tok.AccessToken, Expires: time.Now().Add(time.Duration(expiration) * time.Hour)}
			http.SetCookie(w, &cookie)
			http.Redirect(w, r, "/", http.StatusFound)
		}
	})
}

func SearchArtistHandler(c caching.Client, log logging.Logger, s *Spotify) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tok := r.Context().Value("AuthToken"); tok != nil {
			query := mux.Vars(r)["name"]
			log.Printf("%s: Checking cache for search query ", query)
			results := GetSearchResultsFromCache(query, c, log)

			if len(results) == 0 {
				results = SearchForArtist(query, c, s.NewClient(tok.(string)), log)
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(results); err != nil {
				panic(err)
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
				log.Println(err)
			}

			id, result := PublishPlaylist(newPlaylist.Tracks, newPlaylist.Title, log, s.NewClient(tok.(string)))
			if result {
				w.WriteHeader(http.StatusCreated)
				message = id
			} else {
				w.WriteHeader(http.StatusBadRequest)
				message = "Could not create playlist"
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			message = "Not logged in"
		}
		w.Write([]byte(message))
	})
}
