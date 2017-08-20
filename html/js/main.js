var eventHub = new Vue();

Vue.component('album', {
    template: '#album-col',
    props : {
        album: Object
    },
    data: function () {
        return {
            checkedTracks: [],
            show: false,
            checked: true
        }
    },
    created: function () {

    },
    methods: {
        updateAlbum: function () {
            var self = this;
            self.checked = !self.checked;
            self.album.tracks.items.forEach(function (track) {
                var index = self.checkedTracks.indexOf(track.id);
                if(self.checked) {
                    if(index < 0) {
                        self.checkedTracks.push(track.id)
                    }
                }
                else {
                    if(index > -1) {
                        self.checkedTracks.splice(index, 1)
                    }
                }
                self.updateTrack(track.id)
            });
        },
        updateTrack: function(trackId) {
            this.$emit('update', trackId);
        }
    }

});

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
        loginToken:                 "",
        checkedTracks:              [],
        markets:                    [],
        artistSearchResults:        [],
        allAlbums:                  [],
        loginUrl:                   "",
        artist:                     {},
        sortOption:                 ""
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
        toggleAllAlbums: function () {
            //eventHub.$emit('toggle-album')
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
            console.log(artistId)
            this.$http.get('/tracks/' + artistId).then(function(response) {
                console.log(response)
                if(response.data !== ""){
                    console.log(response)
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
                    this.sortAlbums()
                }
            }).catch(function(error) {
                console.log(error)
            })
        }
    }
});
