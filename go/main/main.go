package main

import (
	"net/http"
	"context"
	"sync"
	"github.com/Briscooe/Discogrify/go/caching"
	"github.com/Briscooe/Discogrify/go/logging"
	"golang.org/x/oauth2/clientcredentials"
	"os"
	"github.com/zmb3/spotify"
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
		config.Redis.Db,
		config.Redis.HoursExpiration)

	spotifyClient := NewSpotifyClient(*logger)

	if config.DevMode {
		conf := &clientcredentials.Config{
			ClientID:     os.Getenv("SPOTIFY_ID"),
			ClientSecret: os.Getenv("SPOTIFY_SECRET"),
			TokenURL:     spotify.TokenURL,
		}
		token,_ := conf.Token(context.Background())
		spotifyClient.Client = spotify.Authenticator{}.NewClient(token)
		_,_ = spotifyClient.Client.CurrentUser()
	}

	router := setupRouter(cacheClient, logger, spotifyClient)
	contextedRouter := AddContext(router)
	logger.Fatal(http.ListenAndServe(":8080", contextedRouter))
}
