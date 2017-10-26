<template>
  <div id="content" class="row align-center margin2">
    <loader v-if="searching"></loader>
    <div v-else-if="resultsPresent" class="col col-6">
      <div id="results-header" class="row gutters">
        <p class="large-text col col-6 margin2"> {{ artistName }}</p>
        <sort :albums="albums" v-on:sort="updateSort" class="col col-6 margin2"></sort>
        <p class="large-text col col-6 margin2"> {{ checkedTracks.length }} tracks selected</p>
        <button class="col col-6 margin2" v-on:click="addMore">Add more artists</button>
        <button class="col col-6 margin2" v-on:click="clear">Clear all</button>
        <button v-if="!publishing" class="col col-6 margin2" v-on:click="publishPlaylist">Publish playlist</button>
        <spinner class="col col-6 margin2" v-else ></spinner>
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
  import Loader from './Loader'

  export default {
    components: {
      'album': Album,
      'sort': Sort,
      'modal': Modal,
      'spinner': Spinner,
      'loader': Loader
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
        searching: false
      }
    },
    mounted () {
      EventBus.$on('albums', this.initialiseAlbums)
      EventBus.$on('artist', this.initialiseArtist)
      EventBus.$on('searching', this.toggleSearching)
    },
    computed: {
      resultsPresent: function () {
        return this.albums.length > 0
      }
    },
    methods: {
      initialiseAlbums: function (allAlbums) {
        let self = this
        allAlbums.forEach(function (album) {
          if (album.tracks.items !== null) {
            if (!self.albums.includes(album)) {
              self.albums.push(album)
            }
            if (!self.checkedAlbums.includes(album.id)) {
              self.checkedAlbums.push(album.id)
            }
            album.tracks.items.forEach(function (track) {
              if (!self.checkedTracks.includes(track.id)) {
                self.checkedTracks.push(track.id)
              }
            })
          }
        })
        self.unfilteredAlbums = self.albums
        self.searching = false
      },
      initialiseArtist: function (artistName) {
        if (!this.artists.includes(artistName)) {
          this.artists.push(artistName)
        }
        if (!this.artistName.includes(artistName)) {
          this.artistName = this.artistName.length === 0 ? artistName : this.artistName += ', ' + artistName
        }
      },
      toggleSearching: function () {
        this.searching = true
      },
      clear: function () {
        this.unfilteredAlbums = []
        this.albums = []
        this.checkedAlbums = []
        this.checkedTracks = []
        this.artists = []
        this.artistName = ''
        EventBus.$emit('clear')
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
      getPlaylistName: function () {
        console.log(this.artists)
        let returnStr = this.artists[0]
        const suffix = ' - By Discogrify'
        this.artists.slice(1, this.artists.length).forEach(function (artist) {
          let tempStr = returnStr + ', ' + artist + suffix
          if (tempStr.length <= 100) {
            returnStr += ', ' + artist
          } else {
            return returnStr + ' and more'
          }
        })
        return returnStr
      },
      publishPlaylist: function () {
        this.published = false
        let playlist = {}
        playlist.tracks = this.checkedTracks
        playlist.name = this.getPlaylistName()
        this.$http.post('/publish', playlist, {
          before: function () {
            this.publishing = true
          }
        }).then(function (response) {
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

.p {
  line-height: 2em;
}
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
  max-width: inherit;
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
