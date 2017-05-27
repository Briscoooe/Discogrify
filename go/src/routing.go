package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	router.Path("/").Handler(http.FileServer(http.Dir("../../html")))

	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Callback",
		"GET",
		"/callback",
		Index,
	},
	Route{
		"Index",
		"GET",
		"/playlists",
		GetPlaylists,
	},
	Route{
		"GetSongsByArtist",
		"POST",
		"/artist",
		GetSongsByArtist,
	},
	Route{
		"GetPlaylist",
		"GET",
		"/playlist/{id}",
		GetPlaylist,
	},
	Route{
		"PublishPlaylist",
		"POST",
		"/playlist/{id}",
		PublishPlaylist,
	},
	Route{
		"Login",
		"GET",
		"/login",
		LoginUser,
	},

	Route{
		"FollowPlaylist",
		"POST",
		"/playlist/{id}",
		FollowPlaylist,
	},
}
