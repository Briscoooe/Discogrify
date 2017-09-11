package main

import (
	"net/http"

	"github.com/Briscooe/Discogrify/go/caching"
	"github.com/Briscooe/Discogrify/go/logging"
	"github.com/gorilla/mux"
)

func SetupRouter(c caching.Client, l logging.Logger, s *Spotify, cookieName string, expiration int, path string) *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	router.Handle("/login", LoginToSpotifyHandlerFunc(s)).Methods("GET")
	router.Handle("/callback", CallbackHandler(l, s, cookieName, expiration)).Methods("GET")
	router.Handle("/tracks/{artistId}", AddContext(GetTracksHandler(c, l, s))).Methods("GET")
	router.Handle("/search/{name}", AddContext(SearchArtistHandler(c, l, s))).Methods("GET")
	router.Handle("/user", AddContext(UserInfoHandler(l,s))).Methods("GET")
	router.Handle("/publish", AddContext(PublishPlaylistHandler(l, s))).Methods("POST")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(path)))

	return router
}
