var app = new Vue({
    el: '#app',

    data: {
        sortOptions:                [
            { text: 'Alphabetical (A-Z)', id: 'alphaAToZ'},
            { text: 'Alphabetical (Z-A)', id: 'alphaZToA'},
            { text: 'Popularity (most first)', id: 'popularMost'},
            { text: 'Popularity (least first)', id: 'popularLeast'} ,
            { text: 'Release date (oldest first)', id: 'dateOldest'} ,
            { text: 'Release date (recent first)', id: 'dateRecent'}  ,
            { text: 'Number of tracks (most first)', id: 'tracksMost'}  ,
            { text: 'Number of tracks (least first)', id: 'tracksLeast'}
        ],
        instrumentalKeys:           [
            "Instrumentals"
        ],
        originalAlbumList:          [],
        loginToken:                 "",
        checkedTracks:              [],
        markets:                    [],
        artistSearchResults:        [],
        allAlbums:                  [],
        loginUrl:                   "",
        artist:                     {},
        sortOption:                 "",
        duplicates:                 false,
        instrumentals:              false

    },

    methods: {
        sortAlbums: function() {
            switch(this.sortOption) {
                case "alphaAToZ":
                    this.allAlbums.sort((albumA, albumB) => albumA.name.localeCompare(albumB.name))
                    break;
                case "alphaZToA":
                    this.allAlbums.sort((albumA, albumB) => albumB.name.localeCompare(albumA.name))
                    break;
                case "popularMost":
                    this.allAlbums.sort((albumA, albumB) => albumB.popularity - albumA.popularity)
                    break;
                case "popularLeast":
                    this.allAlbums.sort((albumA, albumB) => albumA.popularity - albumB.popularity)
                    break;
                case "dateOldest":
                    this.allAlbums.sort((albumA, albumB) => albumA.release_date.localeCompare(albumB.release_date))
                    break;
                case "dateRecent":
                    this.allAlbums.sort((albumA, albumB) => albumB.release_date.localeCompare(albumA.release_date))
                    break;
                case "tracksMost":
                    this.allAlbums.sort((albumA, albumB) => albumB.tracks.items.length - albumA.tracks.items.length)
                    break;
                case "tracksLeast":
                    this.allAlbums.sort((albumA, albumB) => albumA.tracks.items.length - albumB.tracks.items.length)
                    break;
            }
        },
        removeInstrumentals: function () {
            /*
            var self = this;
            if(self.instrumentals) {
                this.originalAlbumList = this.allAlbums;
                var updatedList = this.allAlbums.filter(function (album) {
                    if (album.name.indexOf("Instrumentals") !== -1) {
                        album.tracks.items.forEach(function (track) {
                            self.checkedTracks.splice(self.checkedTracks.indexOf(track.id), 1)
                        });
                        return false
                    } else {
                        album.tracks.items.forEach(function (track) {
                            if(track.name.indexOf("Instrumental") !== -1) {
                                var index = self.checkedTracks.indexOf(track.id);
                                if (index !== -1) {
                                    self.checkedTracks.splice(self.checkedTracks.indexOf(track.id), 1);
                                    return false
                                }
                            }
                        });
                    }
                    return true
                });
                this.allAlbums = updatedList
            } else {
                this.originalAlbumList.forEach(function (album) {
                    if (self.allAlbums.indexOf(album) === -1) {
                        album.tracks.items.forEach(function (track) {
                            self.checkedTracks.push(track.id)
                        })
                    } else {
                        album.tracks.items.forEach(function (track) {
                            if(self.checkedTracks.indexOf(track.id) === -1) {
                                self.checkedTracks.push(track.id)
                            }
                        });
                    }
                });
                this.allAlbums = this.originalAlbumList
            }*/

        },
        removeDuplicates: function () {
            var updatedList = this.allAlbums.filter(function (albumA) {
                return albumA.name.indexOf("Instrumentals")
            });
            this.allAlbums = updatedList
        },
        updateTrack : function (trackId) {
            var index = this.checkedTracks.indexOf(trackId);
            if(index < 0) {
                this.checkedTracks.push(trackId)
            }
            else  {
                this.checkedTracks.splice(index, 1)
            }
        },
        login: function() {
            this.$http.get('/login').then(function(response) {
                window.location.href = response.data.url;
                console.log("Response: " + response.data.url)
            }).catch(function(error) {
                console.log(error)
            })
        },
        publishPlaylist: function() {
            this.$http.post('/publish', this.checkedTracks, "teststring").then(function(response) {
                alert("Created")
            }).catch(function(error) {
                console.log(error)
            })
        },
        searchArtist: function() {
            if (!$.trim(this.artist.name)) {
                this.artist = {};
                return
            }

            this.$http.get('/search/' + encodeURIComponent(this.artist.name)).then(function(response) {
                this.artistSearchResults = response.data;
            }).catch(function(error) {
                console.log(error)
            })
        },

        getTracks: function(artistId) {
            this.$http.get('/tracks/' + artistId).then(function(response) {
                if(response.data !== ""){
                    var checkedTracks = [];
                    var albums = [];
                    response.data.forEach(function(album){
                        if(album.tracks.items !== null) {
                            albums.push(album);
                            album.tracks.items.forEach(function (track) {
                                checkedTracks.push(track.id)
                            });
                        }
                    });
                    this.checkedTracks = checkedTracks;
                    this.allAlbums = albums;
                    this.originalAlbumList = albums;
                    this.sortAlbums()
                }
            }).catch(function(error) {
                console.log(error)
            })
        }
    }
});
