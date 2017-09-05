package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/Briscooe/Discogrify/go/logging"
	"github.com/Briscooe/Discogrify/go/caching"
)

func setupRouter(cacheClient caching.Client, logger logging.Logger, spotify Spotify) *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	router.Handle("/login", loginToSpotifyHandlerFunc(logger, spotify)).Methods("GET")
	router.Handle("/index", indexHandlerFunc(logger)).Methods("GET")
	router.Handle("/callback", callbackHandler(cacheClient, logger, spotify)).Methods("GET")
	router.Handle("/tracks/{artistId}", AddContext(getTracksHandler(cacheClient, logger, spotify))).Methods("GET")
	router.Handle("/search/{name}", AddContext(searchArtistHandler(cacheClient, logger, spotify))).Methods("GET")
	router.Handle("/publish", AddContext(publishPlaylistHandler(cacheClient, spotify))).Methods("POST")
	router.Handle("/user", AddContext(userInfoHandler(cacheClient, spotify))).Methods("GET")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("/home/r00t/go/src/github.com/Briscooe/Discogrify/vue/dist")))

	return router
}