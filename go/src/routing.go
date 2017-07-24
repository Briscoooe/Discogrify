package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func setupRouter(client CacheClient) *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	router.Path("/").Handler(http.FileServer(http.Dir("/home/r00t/go/src/github.com/Briscooe/Discogrify/html")))
	router.HandleFunc("/login", loginToSpotifyHandlerFunc).Methods("GET")
	router.HandleFunc("/index", indexHandlerFunc).Methods("GET")

	router.Handle("/callback", callbackHandler(client)).Methods("GET")
	router.Handle("/tracks/{artistId}", getTracksHandler(client)).Methods("GET")
	router.Handle("/search/{name}", searchArtistHandler(client)).Methods("GET")
	router.Handle("/publish", publishPlaylistHandle(client)).Methods("POST")

	return router
}