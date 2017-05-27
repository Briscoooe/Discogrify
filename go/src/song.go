package main

type Song struct {
	ID      string   `json:"id"`
	Artists []Artist `json:"artists"`
	Title   string   `json:"title"`
}
type Songs []Song
