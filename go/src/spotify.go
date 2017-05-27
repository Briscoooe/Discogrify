package main

var currentId int

var playlists []Playlist

func init() {
	RepoCreatePlaylist(Playlist{Title: "Write presentation"})
	RepoCreatePlaylist(Playlist{Title: "Host meetup"})
}

func RepoCreatePlaylist(p Playlist) Playlist {
	currentId += 1
	p.Followers = currentId
	playlists = append(playlists, p)
	return p
}
