package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/Briscooe/Discogrify/go/logging"
	"github.com/Briscooe/Discogrify/go/caching"
)

func setupRouter(cacheClient caching.Client, logger logging.Logger, spotify Spotify) *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	router.Handle("/login", LoginToSpotifyHandlerFunc(logger, spotify)).Methods("GET")
	router.Handle("/index", IndexHandlerFunc(logger)).Methods("GET")
	router.Handle("/callback", CallbackHandler(logger, spotify)).Methods("GET")
	router.Handle("/tracks/{artistId}", AddContext(GetTracksHandler(cacheClient, logger, spotify))).Methods("GET")
	router.Handle("/search/{name}", AddContext(SearchArtistHandler(cacheClient, logger, spotify))).Methods("GET")
	router.Handle("/publish", AddContext(PublishPlaylistHandler(logger, spotify))).Methods("POST")
	router.Handle("/user", AddContext(UserInfoHandler(spotify))).Methods("GET")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("/home/r00t/go/src/github.com/Briscooe/Discogrify/vue/dist")))

	return router
}