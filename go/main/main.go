package main

import (
	"net/http"
	"sync"
	"github.com/Briscooe/Discogrify/go/caching"
	"github.com/Briscooe/Discogrify/go/logging"
)

const redirectURI = "http://localhost:8080/callback"

var (
	done         = make(chan bool)
	mutex        = &sync.Mutex{}
)

func main() {
	config := loadConfiguration("/home/r00t/go/src/github.com/Briscooe/Discogrify/go/main/config.json")

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
		config.Redis.Db)

	spotifyClient := NewSpotifyClient(*logger)
	router := setupRouter(cacheClient, logger, spotifyClient)

	logger.Fatal(http.ListenAndServe(":8080", router))
}
