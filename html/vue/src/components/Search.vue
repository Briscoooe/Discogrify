<template>
  <div id="content" class="row align-center font-override">
    <div class="col col-6">
      <div class="append">
        <input id="search-box" class="search font-override" placeholder="Search artist..." v-on:keyup.enter="searchArtist" v-model="artist.name">
        <button type="button" class="font-override" v-on:click="searchArtist">Search</button>
      </div>
      <transition name="fade">
        <ul v-if="searchResultsPresent">
          <li v-for="artist in artistSearchResults">
            <div id="result-line" v-on:click="getTracks(artist)">
              {{ artist.name }}
            </div>
          </li>
        </ul>
      </transition>
      <!--<button type="button" @click="showModal = true" >press me</button>
      <modal v-if="showModal" @close="showModal = false">
      </modal>-->
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
            let albums = []
            response.data.forEach(function (album) {
              albums.push(album)
            })
            EventBus.$emit('albums', albums)
            EventBus.$emit('artist', this.artist)
            this.$emit('scroll')
          }
        }).catch(function (error) {
          console.log(error)
        })
      }
    }
  }

</script>

<style scoped>
.font-override {
  font-family: var(--font);
  font-size: var(--font-size-control);
}

#search-box:focus{
  outline: none;
  box-shadow: 0 0 5px var(--primary-green);
  border:1px solid var(--secondary-green);
}

#search-box:hover{
  border: 1px solid var(--primary-grey);
  border-radius: 5px;
}

#search-box:focus:hover{
  outline: none;
  box-shadow: 0 0 5px var(--primary-green);
  border:1px solid var(--secondary-green);
  border-radius:0;

}

#result-line{
  text-align: left;
  padding: 1% 2%;
  cursor: pointer;
  border: 1px solid;
}

#result-line:hover{
  background-color: var(--primary-sand);
}
</style>
