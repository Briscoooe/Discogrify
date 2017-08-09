package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/Briscooe/Discogrify/go/logging"
	"github.com/Briscooe/Discogrify/go/caching"
)

func setupRouter(cacheClient caching.Client, logger logging.Logger, spotify Spotify) *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	router.Path("/").Handler(http.FileServer(http.Dir("/home/r00t/go/src/github.com/Briscooe/Discogrify/html")))
	router.Handle("/login", loginToSpotifyHandlerFunc(logger, spotify)).Methods("GET")
	router.Handle("/index", indexHandlerFunc(logger)).Methods("GET")
	router.Handle("/callback", callbackHandler(cacheClient, logger, spotify)).Methods("GET")
	router.Handle("/tracks/{artistId}", getTracksHandler(cacheClient, logger, spotify)).Methods("GET")
	router.Handle("/search/{name}", searchArtistHandler(cacheClient, logger, spotify)).Methods("GET")
	router.Handle("/publish", publishPlaylistHandler(cacheClient, spotify)).Methods("POST")

	return router
}