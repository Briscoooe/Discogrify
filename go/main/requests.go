package main

import (
	"encoding/json"
	"net/http"
	"github.com/Briscooe/Discogrify/go/logging"
	"github.com/gorilla/mux"
	"time"
	"github.com/Briscooe/Discogrify/go/caching"
	"io/ioutil"
	"regexp"
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"strings"
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
func IndexHandlerFunc(logger logging.Logger) http.Handler {
	return nil
}

func LoginToSpotifyHandlerFunc(logger logging.Logger, spotify Spotify) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := spotify.GenerateLoginUrl()
		w.Write([]byte(url))
	})
}

func GetTracksHandler(cacheClient caching.Client, logger logging.Logger, spotify Spotify) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tok := r.Context().Value("AuthToken"); tok != nil {
			id := mux.Vars(r)["artistId"]
			match, _ := regexp.MatchString("^[a-zA-Z0-9]{22,}$", id)
			if !match {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Invalid artist ID"))
			} else {
				logger.Printf("%s: Checking cache for artist ID", id)
				tracks := GetTracksFromCache(id, cacheClient, logger)
				if tracks == nil {
					tracks = spotify.GetDiscography(id)
					tracksJson, _ := json.Marshal(tracks)
					AddToCache(id, string(tracksJson), cacheClient, logger)
				}
				logger.Printf("%s: Returning tracks", id)
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

func CallbackHandler(logger logging.Logger, spotify Spotify) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tok, err, msg := spotify.ValidateCallback(r)
		if err != nil {
			logger.Println(msg)
			logger.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(msg))
		} else {
			cookie := http.Cookie{Name:"auth_token", Value: tok.AccessToken, Expires: time.Now().Add(time.Hour)}
			http.SetCookie(w,&cookie)
			http.Redirect(w, r, "/", http.StatusFound)
		}
	})
}

func UserInfoHandler(spotify Spotify) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tok := r.Context().Value("AuthToken"); tok != nil {
			user, err := spotify.GetUserInfo(&oauth2.Token{AccessToken: fmt.Sprintf("%v", tok)})
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Invalid token"))
			} else {
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(http.StatusOK)
				if err := json.NewEncoder(w).Encode(user); err != nil {
					panic(err)
				}
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Not logged in"))
		}
	})
}
func SearchArtistHandler(cacheClient caching.Client, logger logging.Logger, spotify Spotify) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tok := r.Context().Value("AuthToken"); tok != nil {
			query := strings.ToLower(mux.Vars(r)["name"])
			logger.Printf("%s: Checking cache for search query ", query)
			results := GetSearchResultsFromCache(query, cacheClient, logger)

			if len(results) == 0 {
				w.WriteHeader(http.StatusNotFound)
				results = spotify.SearchForArtist(query, cacheClient)
			} else {
				w.WriteHeader(http.StatusFound)
			}

			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			if err := json.NewEncoder(w).Encode(results); err != nil {
				panic(err)
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Not logged in"))
		}
	})
}


func PublishPlaylistHandler(logger logging.Logger, s Spotify) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var message string
		if tok := r.Context().Value("AuthToken"); tok != nil {
			body, _ := ioutil.ReadAll(r.Body)
			type playlist struct {
				Tracks []string
				Title string `json:"name"`
			}

			var newPlaylist playlist
			err := json.Unmarshal(body, &newPlaylist)
			if err != nil{
				logger.Println(err)
			}

			result, msg := s.PublishPlaylist(newPlaylist.Tracks, newPlaylist.Title)
			if result {
				w.WriteHeader(http.StatusCreated)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
			message = msg
		}else {
			w.WriteHeader(http.StatusUnauthorized)
			message = "Not logged in"
		}
		w.Write([]byte(message))
	})
}
