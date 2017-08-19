Vue.component('v-header', {
    template: '<div><h1>Discogrify</h1></div>'
});
Vue.component('v-footer', {
    template: '<div style="height:10%; background-color: #95d6ff">Brian Briscoe 2016</div>'
});
new Vue({
    el: 'body',

    data: {
        loginToken:                 "",
        checkedTracks:              [],
        markets:                    [],
        artistSearchResults:        [],
        allTracks:                  [],
        loginUrl:                   "",
        artistId:                   {},
        sortOption:                 ""
    },

    computed :{
        totalTracks: function () {
            var count = 0;
            this.allTracks.forEach(function(album){
                count = count + album.tracks.items.length
            });
            return count
        }
    },

    methods: {
        sortAlbums: function() {
            switch(this.sortOption) {
                case "Alphabetical (A-Z)":
                    this.allTracks.sort((albumA, albumB) => albumA.name.localeCompare(albumB.name))
                    break;
                case "Alphabetical (Z-A)":
                    this.allTracks.sort((albumA, albumB) => albumB.name.localeCompare(albumA.name))
                    break;
                case "Popularity (most first)":
                    this.allTracks.sort((albumA, albumB) => albumA.popularity.localeCompare(albumB.popularity))
                    break;
                case "Popularity (least first)":
                    this.allTracks.sort((albumA, albumB) => albumB.popularity.localeCompare(albumA.popularity))
                    break;
                case "Release date (oldest first)":
                    this.allTracks.sort((albumA, albumB) => albumA.release_date.localeCompare(albumB.release_date))
                    break;
                case "Release date (recent first)":
                    this.allTracks.sort((albumA, albumB) => albumB.release_date.localeCompare(albumA.release_date))
                    break;
                case "Number of tracks (most first)":
                    this.allTracks.sort((albumA, albumB) => albumB.tracks.items.length - albumA.tracks.items.length)
                    break;
                case "Number of tracks (least first)":
                    this.allTracks.sort((albumA, albumB) => albumA.tracks.items.length - albumB.tracks.items.length)
                    break;
            }
        },
        checkAlbum: function (album, checkedTracks) {
            album.isChecked = !album.isChecked;
            album.tracks.items.forEach(function (track) {
                this.checkedTracks = checkedTracks;
                var index = this.checkedTracks.indexOf(track.id);
                console.log(album.isChecked)
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
                this.allTracks = response;
                this.allTracks.forEach(function(album){
                    album.isVisible = false;
                    album.isChecked = true;
                });
            }).error(function(error) {
                console.log(error)
            })
        }
    }
});
