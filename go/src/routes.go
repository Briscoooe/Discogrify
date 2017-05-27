package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"GetSongsByArtist",
		"GET",
		"/artist/{id}",
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
		"Authenticate",
		"POST",
		"/user/{id}",
		AuthenticateUser,
	},
	Route{
		"FollowPlaylist",
		"POST",
		"/playlist/{id}",
		FollowPlaylist,
	},
}
