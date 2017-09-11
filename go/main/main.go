package main

import (
	"github.com/Briscooe/Discogrify/go/caching"
	"github.com/Briscooe/Discogrify/go/logging"
	"net/http"
	"sync"
	"github.com/rs/cors"
)

var (
	done  = make(chan bool)
	mutex = &sync.Mutex{}
)

func main() {
	config := LoadConfiguration("./config.json")

	logger := logging.NewRollingLogger(
		config.Logger.Filename,
		config.Logger.MaxSize,
		config.Logger.MaxBackups,
		config.Logger.MaxAge)
	logger.Println("Starting application...")

	cacheClient := caching.NewRedisClient(
		*logger,
		config.Redis.Host,
		config.Redis.Port,
		config.Redis.Password,
		config.Redis.Db,
		config.Redis.HoursExpiration)

	spotifyClient := InitSpotifyClient(config.Spotify.RedirectURI)

	router := SetupRouter(cacheClient, logger, spotifyClient, config.Cookie.CookieName, config.Cookie.Expiration, config.FilePath)
	contextedRouter := AddContext(router)

	handler := cors.Default().Handler(contextedRouter)
	logger.Fatal(http.ListenAndServe(":8080", handler))
}
