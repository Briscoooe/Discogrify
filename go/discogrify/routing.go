package discogrify

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/Briscooe/Discogrify/go/logging"
	"github.com/Briscooe/Discogrify/go/caching"
)

func SetupRouter(c caching.Client, log logging.Logger, s *Spotify, cookieName string, expiration int) *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	router.Handle("/login", LoginToSpotifyHandlerFunc(s)).Methods("GET")
	router.Handle("/index", IndexHandlerFunc(log)).Methods("GET")
	router.Handle("/callback", CallbackHandler(log, s, cookieName, expiration)).Methods("GET")
	router.Handle("/tracks/{artistId}", AddContext(GetTracksHandler(c, log, s))).Methods("GET")
	router.Handle("/search/{name}", AddContext(SearchArtistHandler(c, log, s))).Methods("GET")
	router.Handle("/publish", AddContext(PublishPlaylistHandler(log, s))).Methods("POST")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("/home/r00t/go/src/github.com/Briscooe/Discogrify/vue/dist")))

	return router
}