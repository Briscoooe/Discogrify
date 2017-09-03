<template>
  <div id="content" class="row align-center font-override">
    <div class="col col-6">
      <div class="append">
        <input class="search" placeholder="Search artist..." class="font-override" v-on:keyup.enter="searchArtist" v-model="artist.name">
        <button type="button" class="button outline" class="font-override" v-on:click="searchArtist">Search</button>
      </div>
      <transition name="fade">
        <ul v-if="searchResultsPresent">
          <li v-for="artist in artistSearchResults">
            <div id="search-line" v-on:click="getTracks(artist)">
              {{ artist.name }}
            </div>
          </li>
        </ul>
      </transition>
      <button type="button" @click="showModal = true" >press me</button>
      <modal v-if="showModal" @close="showModal = false">
      </modal>
    </div>
  </div>
</template>

<script>
  import EventBus from '../event-bus'
  import Modal from './Modal'
  export default {
    components: {
      'modal': Modal
    },
    data () {
      return {
        artistSearchResults: [],
        artist: {},
        showModal: false
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
#content {
  margin-top: 5%;
}
.font-override {
  font-size:20px;
}

.fade-enter-active, .fade-leave-active {
  transition: opacity .2s
}
.fade-enter, .fade-leave-to {
  opacity: 0
}

</style>
