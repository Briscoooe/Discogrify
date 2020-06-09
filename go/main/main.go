package main

import (
	"github.com/Briscoooe/Discogrify/go/caching"
	"github.com/Briscoooe/Discogrify/go/logging"
	"net/http"
	"sync"
	"github.com/rs/cors"
)

var (
	done  = make(chan bool)
	mutex = &sync.Mutex{}
	spotifyAuthStateCookie string
	authTokenCookie string
)

func main() {
	config := LoadConfiguration("./config.json")

	l := logging.NewRollingLogger(
		config.Logger.Filename,
		config.Logger.MaxSize,
		config.Logger.MaxBackups,
		config.Logger.MaxAge)
	l.Log("Starting application...")

	c := caching.NewRedisClient(
		*l,
		config.Redis.Host,
		config.Redis.Port,
		config.Redis.Password,
		config.Redis.Db,
		config.Redis.HoursExpiration)

	s := InitSpotifyClient(config.Spotify.RedirectURI)

	spotifyAuthStateCookie = config.Cookie.SpotifyAuthState
	authTokenCookie = config.Cookie.AuthToken
	router := SetupRouter(c, l, s, config.Cookie.Expiration, config.FilePath)
	contextedRouter := AddContext(router)

	handler := cors.Default().Handler(contextedRouter)
	l.Fatal(http.ListenAndServe(":8080", handler))
}
