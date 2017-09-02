<template>
  <div>
    <p> {{ artist.name }}</p>
    <div>
      <button id="create-button" v-on:click="publishPlaylist">Create playlist</button>
    </div>
    <table>
      <thead>
      <tr>
        <th><input type="checkbox" checked="true"></th>
        <th>Album</th>
        <th>Album artist(s)</th>
      </tr>
      </thead>
      <tbody>
      <tr  v-for="album in allAlbums" :key="album.id">
        <td>
          <input type="checkbox">
        </td>
        <td  v-if="album.tracks.items.length > 0">
          <album v-model="checkedTracks" v-on:update="updateTrack" :album="album"></album>
        </td>
        <td>
          <div >
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
          </div>
        </td>
      </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
  import Album from 'Album'

  export default {
    components: {
      'album': Album
    },
    props: {
      results: []
    },
    data () {
      return {
        artistSearchResults: [],
        allAlbums: [],
        checkedTracks: [],
        originalAlbumList: [],
        artist: {}
      }
    },
    created: function () {
      let checkedTracks = []
      let albums = []
      response.data.forEach(function (album) {
        if (album.tracks.items !== null) {
          albums.push(album)
          album.tracks.items.forEach(function (track) {
            checkedTracks.push(track.id)
          })
        }
      })
      this.checkedTracks = checkedTracks
      this.allAlbums = albums
    },
    computed: {

    },
    methods: {
      updateTrack: function (trackId) {
        let index = this.checkedTracks.indexOf(trackId)
        if (index < 0) {
          this.checkedTracks.push(trackId)
        }
        else  {
          this.checkedTracks.splice(index, 1)
        }
      },
      publishPlaylist: function () {
        let playlist = {}
        playlist.tracks = this.checkedTracks
        playlist.name = this.artist.name

        this.$http.post('/publish', playlist).then(function (response) {
          alert("Created")
        }).catch(function (error) {
          console.log(error)
        })
      },
    }
  }

</script>

<style scoped>

</style>
