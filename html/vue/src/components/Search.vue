<template>
  <div>
    <div>
      <input placeholder="Search artist..." v-on:keyup.enter="searchArtist" v-model="artist.name">
      <button type="button" v-on:click="searchArtist">Search</button>
    </div>

    <transition name="fade">
      <ul v-if="searchResultsPresent">
        <li v-for="artist in artistSearchResults">
          <div v-on:click="getTracks(artist)">
            {{ artist.name }}
          </div>
        </li>
      </ul>
    </transition>
  </div>
</template>

<script>
  import EventBus from '../event-bus'

  export default {
    data () {
      return {
        artistSearchResults: [],
        artist: {}
      }
    },

    computed: {
      searchResultsPresent: function () {
        return this.artistSearchResults !== null && this.artistSearchResults.length > 0
      }
    },
    methods: {
      searchArtist: function () {
        this.$http.get('/search/' + encodeURIComponent(this.artist.name)).then(function (response) {
          this.artistSearchResults = response.data
        }).catch(function (error) {
          console.log(error)
        })
      },
      getTracks: function (artist) {
        this.artist.name = artist.name
        this.$http.get('/tracks/' + artist.id).then(function (response) {
          if (response.data !== '') {
            var albums = []
            response.data.forEach(function (album) {
              albums.push(album)
            })
            EventBus.$emit('albums', albums)
            EventBus.$emit('artist', this.artist)
          }
        }).catch(function (error) {
          console.log(error)
        })
      }
    }
  }

</script>

<style scoped>


</style>
