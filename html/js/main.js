Vue.component('album', {
    template: '#album-col',
    props : {
        album: Object
    },
    data: function () {
        return {
            checkedTracks: [],
            show: false
        }
    },
    methods: {
        checkAlbum: function (album, checkedTracks) {
            album.isChecked = !album.isChecked;
            album.tracks.items.forEach(function (track) {
                this.checkedTracks = checkedTracks;
                var index = this.checkedTracks.indexOf(track.id);
                if(album.isChecked) {
                    if(index < 0) {
                        this.checkedTracks.push(track.id)
                    }
                }
                else {
                    if(index > -1) {
                        this.checkedTracks.splice(index, 1)
                    }
                }
            });
        },
        addTrack : function (trackId) {
            this.$emit('add', trackId)
        }
    }
});
new Vue({
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
        artistId:                   {},
        sortOption:                 ""
    },

    computed :{
        totalTracks: function () {
            var count = 0;
            this.allAlbums.forEach(function(album){
                count = count + album.tracks.items.length
            });
            return count
        }
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
                    this.allAlbums.sort((albumA, albumB) => albumA.popularity.localeCompare(albumB.popularity))
                    break;
                case "popularLeast":
                    this.allAlbums.sort((albumA, albumB) => albumB.popularity.localeCompare(albumA.popularity))
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
        addTrack : function (trackId) {
            this.checkedTracks.push(trackId)
        },
        login: function() {
            this.$http.get('/login').then(function(response) {
                window.location.href = response.data.url;
                console.log("Response: " + response.data.url)
            }).error(function(error) {
                console.log(error)
            })
        },
        publishPlaylist: function() {
            this.$http.post('/publish', this.checkedTracks, "teststring").then(function(response) {
                alert("Created")
            }).error(function(error) {
                console.log(error)
            })
        },
        searchArtist: function() {
            if (!$.trim(this.artist.name)) {
                this.artist = {};
                return
            }

            this.$http.get('/search/' + encodeURIComponent(this.artist.name)).success(function(response) {
                this.artistSearchResults = response
            }).error(function(error) {
                console.log(error)
            })
        },

        getTracks: function(artistId) {
            this.$http.get('/tracks/' + artistId).success(function(response) {
                this.checkedTracks = [];
                this.allAlbums = response;
                this.allAlbums.forEach(function(album){
                    album.isChecked = true;
                });
            }).error(function(error) {
                console.log(error)
            })
        }
    }
});
