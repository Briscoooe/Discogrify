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
func indexHandlerFunc(logger logging.Logger) http.Handler {
	return nil
}

func loginToSpotifyHandlerFunc(logger logging.Logger, spotify Spotify) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := spotify.GenerateLoginUrl()
		w.Write([]byte(url))
	})
}

func getTracksHandler(cacheClient caching.Client, logger logging.Logger, spotify Spotify) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if token := r.Context().Value("AuthToken"); token != nil {
			artistId := mux.Vars(r)["artistId"]
			match, _ := regexp.MatchString("^[a-zA-Z0-9]{22,}$", artistId)
			if !match {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Invalid artist ID"))
			} else {
				logger.Printf("%s: Checking cache for artist ID", artistId)
				artistTracks := GetTracksFromCache("artist:" + artistId + ":tracks", cacheClient)
				if artistTracks == nil {
					logger.Printf("%s: Artist ID not found", artistId)
					artistTracks = spotify.GetDiscography(artistId)
					tracksJson, _ := json.Marshal(artistTracks)
					if AddToCache("artist:" + artistId + ":tracks", string(tracksJson), time.Hour * 168, cacheClient) {
						IncrementKeyInCache("artist:" + artistId + ":searched", cacheClient)
						logger.Printf("%s: Successfully added artist to cache", artistId)
					} else {
						logger.Printf("%s: Could not add artist to cache", artistId)
					}
				} else {
					IncrementKeyInCache("artist:" + artistId + ":searched", cacheClient)
					logger.Printf("%s: Artist ID found in cache", artistId)
				}

				logger.Printf("%s: Returning tracks", artistId)
				if err := json.NewEncoder(w).Encode(artistTracks); err != nil {
					panic(err)
				}
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Not logged in"))
		}
	})
}

func callbackHandler(logger logging.Logger, spotify Spotify) http.Handler {
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

func userInfoHandler(spotify Spotify) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if token := r.Context().Value("AuthToken"); token != nil {
			user, err := spotify.GetUserInfo(&oauth2.Token{AccessToken: fmt.Sprintf("%v", token)})
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
func searchArtistHandler(cacheClient caching.Client, logger logging.Logger, spotify Spotify) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if token := r.Context().Value("AuthToken"); token != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			query := mux.Vars(r)["name"]
			logger.Printf("%s: Checking cache for search query ", query)
			results := GetSearchResultsFromCache("artist:search:"+query, cacheClient)
			if len(results) == 0 {
				w.WriteHeader(http.StatusNotFound)
				logger.Printf("%s: Query not found", query)
				results = spotify.SearchForArtist(query, cacheClient)
			} else {
				w.WriteHeader(http.StatusFound)
				logger.Printf("%s: Query results found ", query)
			}
			if err := json.NewEncoder(w).Encode(results); err != nil {
				panic(err)
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Not logged in"))
		}
	})
}


func publishPlaylistHandler(logger logging.Logger, s Spotify) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var message string
		if token := r.Context().Value("AuthToken"); token != nil {
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
