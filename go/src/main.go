package main

import (
	"net/http"
)

var (
	config Configuration
)
func main() {
	setupLogger()
	rollingLog.Println("Starting application...")
	config = loadConfiguration("/home/r00t/go/src/github.com/Briscooe/Discogrify/go/src/config.json")
	setupRedisClient()
	router := setupRouter()
	rollingLog.Fatal(http.ListenAndServe(":8081", router))
}
