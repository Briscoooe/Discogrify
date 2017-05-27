package main

import (
	"time"
)

type Playlist struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Followers   int       `json:"followers"`
	DateCreated time.Time `json:"dateCreated"`
	CreatedBy   string    `json:"userId"`
	Songs       []Songs   `json:"songs"`
}
type Playlists []Playlist
