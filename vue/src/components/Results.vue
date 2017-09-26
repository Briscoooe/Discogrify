# [Discogrify](https://discogrify.com/)
**TL;DR**: A website for creating discographies on Spotify. Login via Spotify, search for an artist, and discography is magically created.

#Description
On Spotify there's no way to add all of an artist's albums, singles and features into a single playlist other than doing it manually. Spotify does have an "Add to playlist" option beside albums but this can take time for artists with many albums and singles.

This is where [Discogrify](https://discogrify.com/) steps in. Using the site, all you have to do is login via Spotify and search for an artist and you'll be presented with an artist's entire Spotify discography. All tracks from all albums, singles, features and compilations. You can optionally edit the selection of tracks to exclude individual tracks or entire albums. Once you're happy with the selection, simply publish the playlist and it will be in your library.

#Usage
First things first, UI design is not my forte so I'm hoping for a pass on that one. Any suggestions are more than welcome. Also the site is not optimised for mobile but will still function the same.
Secondly, ***If you notice anything wrong with the site, please leave a comment so I can fix the issue*** :)
<template>
  <div id="content" class="row align-center">
    <div v-if="resultsPresent" class="col col-6">
      <div id="options" class="margin2">
        <input type="checkbox" class="option" v-model="multipleArtists">&nbsp;Multiple artists
        <input type="checkbox" class="option" v-on:change="filterResults('commentary')"/>&nbsp;Exclude commentary albums
        <input type="checkbox" class="option" v-on:change="filterResults('instrumental')"/>&nbsp;Exclude instrumentals
      </div>
      <div id="results-header" class="margin2">
        <div class="row">
          <p class="large-text col col-6 margin2"> {{ artistName }}</p>
          <sort :albums="albums" v-on:sort="updateSort" class="col col-6 margin2"></sort>
        </div>
        <div class="row">
          <p class="large-text col col-6 margin2"> {{ checkedTracks.length }} tracks selected</p>
          <button class="col col-6 margin2" v-on:click="addMore">Add more tracks</button>
          <button v-if="!publishing" class="col col-6 margin2" v-on:click="publishPlaylist">Publish playlist</button>
          <spinner class="col col-6 margin2" v-else ></spinner>
        </div>
      </div>
      <table class="bordered striped" id="table">
        <thead>
        <tr>
          <th></th>
          <th>Album</th>
          <th>Tracks</th>
          <th>Album artist(s)</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="album in albums" :key="album.id">
          <td>
            <input type="checkbox"
                   :checked="checkedAlbums.indexOf(album.id) > -1"
                   :value="album.id"
                   v-on:change="toggleAlbum(album.id)">
          </td>
          <td v-if="album.tracks.items.length > 0">
            <album v-model="checkedTracks"
                   v-on:update-track="updateTrack"
                   v-on:update-album="updateAlbum"
                   :album="album"></album>
          </td>
          <td>
            {{ album.tracks.items.length }}
          </td>
          <td>
            <div v-if="album.isVisible">
              <ul>
                <li v-for="artist in album.artists">
                  {{ artist.name }}
                </li>
              </ul>
            </div>
            <div v-else>
              {{ album.artists[0].name }}<span v-if="album.artists.length > 1">, {{ album.artists.length -1  }} more</span>
            </div>
          </td>
        </tr>
        </tbody>
      </table>
    </div>
    <modal :show="published" v-if="showModal" @close="showModal = false">
      <span slot="header">
        {{ createPlaylistMessage }}
      </span>
      <span slot="header2">
        <a :href="playlistUrl">View on Spotify</a>
      </span>
    </modal>
  </div>
</template>

<script>
  import Album from './Album'
  import EventBus from '../event-bus'
  import Sort from './Sort'
  import Modal from './Modal'
  import Spinner from './Spinner'
  export default {
    components: {
      'album': Album,
      'sort': Sort,
      'modal': Modal,
      'spinner': Spinner
    },
    props: {
      results: []
    },
    data () {
      return {
        unfilteredAlbums: [],
        albums: [],
        checkedTracks: [],
        checkedAlbums: [],
        artistName: '',
        artists: [],
        createPlaylistMessage: '',
        showModal: false,
        published: false,
        playlistUrl: '',
        publishing: false,
        multipleArtists: false,
        filters: [
          {key: 'commentary', filtered: false, words: ['commentary']},
          {key: 'instrumental', filtered: false, words: ['instrumental', 'instrumentals']}
        ]
      }
    },
    mounted () {
      EventBus.$on('albums', this.initialiseAlbums)
      EventBus.$on('artist', this.initialiseArtist)
    },
    computed: {
      resultsPresent: function () {
        return this.albums.length > 0
      }
    },
    methods: {
      initialiseAlbums: function (allAlbums) {
        let self = this
        if (!self.multipleArtists) {
          self.unfilteredAlbums = []
          self.albums = []
          self.checkedAlbums = []
          self.checkedTracks = []
        }
        allAlbums.forEach(function (album) {
          if (album.tracks.items !== null) {
            self.albums.push(album)
            self.checkedAlbums.push(album.id)
            album.tracks.items.forEach(function (track) {
              self.checkedTracks.push(track.id)
            })
          }
        })
        self.unfilteredAlbums = self.albums
      },
      initialiseArtist: function (artistName) {
        if (this.multipleArtists) {
          this.artists.push(artistName)
          this.artistName += ', ' + artistName
        } else {
          this.artistName = artistName
        }
      },
      toggleAlbum: function (albumId) {
        EventBus.$emit('toggle-album', albumId, this.checkedAlbums.indexOf(albumId) > -1)
      },
      updateAlbum: function (albumId) {
        this.updateArray(this.checkedAlbums, albumId)
      },
      updateTrack: function (trackId) {
        this.updateArray(this.checkedTracks, trackId)
      },
      updateArray: function (array, element) {
        let index = array.indexOf(element)
        if (index > -1) {
          array.splice(index, 1)
        } else {
          array.push(element)
        }
      },
      updateSort: function (albums) {
        this.albums = albums
      },
      filterResults: function (filterKey) {
        let filtered = false
        let words = []
        console.log(filterKey)
        this.filters.forEach(function (filter) {
          if (filter.key === filterKey) {
            filter.filtered = !filter.filtered
            filtered = filter.filtered
            words = filter.words
          }
        })
        if (filtered) {
          let filteredList = this.albums.filter(function (album) {
            words.forEach(function (word) {
              console.log(album)
              console.log(word)
              if (album.name.indexOf(word) !== -1) {
                album.tracks.items.forEach(function (track) {
                  this.checkedTracks.splice(this.checkedTracks.indexOf(track.id), 1)
                })
                return false
              } else {
                album.tracks.items.forEach(function (track) {
                  if (track.name.indexOf(word) !== -1) {
                    let index = this.checkedTracks.indexOf(track.id)
                    if (index !== -1) {
                      this.checkedTracks.splice(this.checkedTracks.indexOf(track.id), 1)
                      return false
                    }
                  }
                })
              }
              return true
            })
          })
          console.log(filteredList)
          this.albums = filteredList
        }
      },
      /* removeInstrumentals: function () {
        var self = this;
        if(self.instrumentals) {
          var updatedList = this.allAlbums.filter(function (album) {
            if (album.name.indexOf("Instrumentals") !== -1) {
              album.tracks.items.forEach(function (track) {
                self.checkedTracks.splice(self.checkedTracks.indexOf(track.id), 1)
              });
              return false
            }
            else {
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
            }
            else {
              album.tracks.items.forEach(function (track) {
                if(self.checkedTracks.indexOf(track.id) === -1) {
                  self.checkedTracks.push(track.id)
                }
              });
            }
          });
          this.allAlbums = this.originalAlbumList
        }

      }, */
      publishPlaylist: function () {
        this.published = false
        let playlist = {}
        playlist.tracks = this.checkedTracks
        playlist.name = this.artistName
        this.$http.post('/publish', playlist, {
          before: function () {
            this.publishing = true
          }
        }).then(function (response) {
          console.log(response)
          this.published = true
          this.playlistUrl = response.data
          this.createPlaylistMessage = 'Playlist created successfully'
        }).catch(function (error) {
          console.log(error)
          switch (error.status) {
            case 304:
              this.published = false
              this.createPlaylistMessage = 'Playlist could not be created. Please try again'
              break
            case 400:
              this.published = true
              this.playlistUrl = error.data
              this.createPlaylistMessage = 'Not all tracks could be added to the playlist'
              break
            case 401:
              this.published = false
              this.createPlaylistMessage = 'Error in authorization. Please log in'
              break
            case 404:
              this.published = false
              this.createPlaylistMessage = 'No tracks present. Playlist not created'
              break
          }
        }).then(function () {
          this.publishing = false
          this.showModal = true
        })
      },
      addMore: function () {
        this.$emit('scroll')
      }
    }
  }
</script>

<style scoped>
  #content {
    width:100%;
  }
  .large-text{
    font-size: var(--font-size-heading);
    text-align: left;
  }

  #results-header {
    background-color: var(--primary-sand);
    padding: 2%;
    border-radius: 3px;
  }
  #table {
    text-align: left;
    font-size: var(--font-size-data);
    width:100%;
  }

  /* Icon Fade */
  .hvr-icon-fade {
    display: inline-block;
    vertical-align: middle;
    -webkit-transform: perspective(1px) translateZ(0);
    transform: perspective(1px) translateZ(0);
    box-shadow: 0 0 1px transparent;
    position: relative;
    padding-right: 2.2em;
  }
  .hvr-icon-fade:before {
    content: "\f00c";
    position: absolute;
    right: 1em;
    padding: 0 1px;
    font-family: FontAwesome;
    -webkit-transform: translateZ(0);
    transform: translateZ(0);
    -webkit-transition-duration: 0.5s;
    transition-duration: 0.5s;
    -webkit-transition-property: color;
    transition-property: color;
  }
  .hvr-icon-fade:hover:before, .hvr-icon-fade:focus:before, .hvr-icon-fade:active:before {
    color: #0F9E5E;
  }
</style>
