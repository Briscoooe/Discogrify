<template>
  <div id="content" class="row align-center">
    <div v-if="resultsPresent" class="col col-6">
      <div id="results-header">
        <div class="row">
          <p class="margin large-text col col-6"> {{ artist.name }}</p>
          <sort :albums="allAlbums" v-on:sort="updateSort" class="margin col col-6"></sort>
        </div>
        <div class="row">
          <p class="margin large-text col col-6"> {{ checkedTracks.length }} tracks selected</p>
          <button class="margin col col-6" v-on:click="publishPlaylist">Publish playlist</button>
        </div>
      </div>
      <table class="striped" id="table">
        <thead>
        <tr>
          <th></th>
          <th>Album</th>
          <th>Tracks</th>
          <th>Album artist(s)</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="album in allAlbums" :key="album.id">
          <td>
            <input type="checkbox"
                   :value="album.id"
                   :checked="checkedAlbums.indexOf(album.id) > -1"
                   v-model="checkedAlbums"
                   @click="toggleSingleAlbum(album.id)">
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
    <div id="no-results" class="col col-6 margin" v-else>
      <div>Results will appear here when you perform a search</div>
    </div>
  </div>
</template>

<script>
  import Album from './Album'
  import EventBus from '../event-bus'
  import Sort from './Sort'
  import Modal from './Modal'
  export default {
    components: {
      'album': Album,
      'sort': Sort,
      'modal': Modal
    },
    props: {
      results: []
    },
    data () {
      return {
        allAlbums: [],
        checkedTracks: [],
        artist: {},
        checkedAlbums: []
      }
    },
    mounted () {
      EventBus.$on('albums', this.initialiseAlbums)
      EventBus.$on('artist', this.initialiseArtist)
    },
    computed: {
      resultsPresent: function () {
        return this.checkedTracks.length > 0
      }
    },
    methods: {
      toggleSingleAlbum: function (albumId) {
        EventBus.$emit('toggle-album', albumId)
        this.updateAlbum(albumId)
      },
      updateAlbum: function (albumId, val) {
        console.log(val)
        let index = this.checkedAlbums.indexOf(albumId)
        console.log('index: ' + index)
        if (index > -1) {
          this.checkedAlbums.splice(index, 1)
        } else {
          this.checkedAlbums.push(albumId)
        }
      },
      initialiseAlbums: function (allAlbums) {
        let self = this
        self.allAlbums = []
        self.checkedAlbums = []
        self.checkedTracks = []
        allAlbums.forEach(function (album) {
          if (album.tracks.items !== null) {
            self.allAlbums.push(album)
            self.checkedAlbums.push(album.id)
            album.tracks.items.forEach(function (track) {
              self.checkedTracks.push(track.id)
            })
          }
        })
      },
      initialiseArtist: function (artist) {
        this.artist = artist
      },
      updateSort: function (albums) {
        this.allAlbums = albums
      },
      updateTrack: function (trackId) {
        let index = this.checkedTracks.indexOf(trackId)
        if (index < 0) {
          this.checkedTracks.push(trackId)
        } else {
          this.checkedTracks.splice(index, 1)
        }
      },
      publishPlaylist: function () {
        let playlist = {}
        playlist.tracks = this.checkedTracks
        playlist.name = this.artist.name
        this.$http.post('/publish', playlist).then(function (response) {
          alert('Created')
        }).catch(function (error) {
          console.log(error)
        })
      }
    }
  }

</script>

<style scoped>
.large-text{
  font-size: var(--font-size-heading);
  text-align: left;
}
.margin {
  margin-top: 2%;
  margin-bottom: 2%;
}
#results-header {
  background-color: var(--primary-sand);
  padding: 2%;
  margin-top: 2%;
  margin-bottom: 2%;
  border-radius: 3px;
}
#table {
  text-align: left;
  font-size: var(--font-size-data);
}


#no-results {
  font-size: var(--font-size-control);
  padding: 2%;
  height: 100%;
  background-color: var(--primary-sand);
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
