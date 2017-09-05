<template>
  <div id="content" class="row align-center">
    <div class="col col-6">
      <div class="row">
        <p id="artist-name" class="margin large-text col col-6"> {{ artist.name }}</p>
        <sort v-on:sort="updateSort" class="margin col col-6"></sort>
      </div>
      <div class="row">
        <p class="margin large-text col col-6"> {{ numberOfTracks }} tracks selected</p>
        <button class="margin col col-6" v-on:click="publishPlaylist">Create playlist</button>
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
  </div>
</template>

<script>
  import Album from './Album'
  import EventBus from '../event-bus'
  import Sort from './Sort'
  export default {
    components: {
      'album': Album,
      'sort': Sort
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
    created: function () {

    },
    mounted () {
      EventBus.$on('albums', this.initialiseAlbums)
      EventBus.$on('artist', this.initialiseArtist)
    },
    computed: {
      numberOfTracks: function () {
        return this.checkedTracks.length
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
  font-size: var(--font-size-title);
}
.margin {
  margin-top: 1%;
  margin-bottom: 1%;
}
#artist-name {
  text-align:left;
}
#table {
  text-align: left;
  font-size: var(--font-size-data);
}
</style>
