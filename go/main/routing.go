package main

import (
	"net/http"

	"github.com/Briscoooe/Discogrify/go/caching"
	"github.com/Briscoooe/Discogrify/go/logging"
	"github.com/gorilla/mux"
)

func SetupRouter(c caching.Client, l logging.Logger, s *Spotify, expiration int, path string) *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	router.Handle("/login", LoginToSpotifyHandlerFunc(s, expiration)).Methods("GET")
	router.Handle("/callback", CallbackHandler(l, s, expiration)).Methods("GET")
	router.Handle("/tracks/{artistId}", AddContext(GetTracksHandler(c, l, s))).Methods("GET")
	router.Handle("/search/{name}", AddContext(SearchArtistHandler(c, l, s))).Methods("GET")
	router.Handle("/user", AddContext(UserInfoHandler(l,s))).Methods("GET")
	router.Handle("/publish", AddContext(PublishPlaylistHandler(l, s))).Methods("POST")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(path)))

	return router
}
