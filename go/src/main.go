package main

import (
	"log"
	"net/http"
	"fmt"
)

var (
	config Configuration
)
func main() {
	config = loadConfiguration("config.json")
	setupClient()
	router := NewRouter()

	fmt.Print(config)
	log.Fatal(http.ListenAndServe(":8080", router))
}
