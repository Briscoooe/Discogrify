<template>
  <div id="content" class="row align-center font-override margin2">
    <div class="col col-6">
      <div class="append">
        <input class="search font-override green-control"
               placeholder="Search artist..."
               v-on:keyup.enter="searchArtist"
                v-model="artist.name">
        <button type="button" class="font-override" v-on:click="searchArtist">Search</button>
      </div>
      <div v-if="loggedIn && artistSearchResults.length == 0" id="no-results" class="col col-6 margin2">
        <div> {{ noResultsMessage }} </div>
      </div>
      <transition name="fade">
        <ul id="list" v-if="show">
          <li id="list-item" v-for="artist in artistSearchResults">
            <div id="result-line" class="hvr-underline-from-left" v-on:click="getTracks(artist)">
              {{ artist.name }}<span class="hvr-icon-forward"></span>
            </div>
          </li>
        </ul>
      </transition>
      <modal v-if="showModal" @close="showModal = false">
        <span slot="header">You must login through Spotify to continue</span>
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
        cookieName: 'auth_token',
        artistSearchResults: [],
        artist: {},
        showModal: false,
        show: false,
        noResultsMessage: 'Results will appear here when you search'
      }
    },
    computed: {
      loggedIn: function () {
        return document.cookie.match('(^|;)\\s*' + this.cookieName + '\\s*=\\s*([^;]+)')
      }
    },
    methods: {
      searchArtist: function () {
        if (!this.loggedIn) {
          this.showModal = true
          return
        } else {
          this.showModal = false
        }
        if (!this.artist.name) {
          return
        }
        this.$http.get('/search/' + encodeURIComponent(this.artist.name)).then(function (response) {
          if (response.data) {
            this.artistSearchResults = response.data
            this.show = true
          } else {
            this.noResultsMessage = 'No results for "' + this.artist.name + '"'
          }
        }).catch(function (error) {
          console.log(error)
        })
      },
      getTracks: function (artist) {
        this.artist.name = artist.name
        this.$http.get('/tracks/' + artist.id).then(function (response) {
          if (response.data) {
            console.log(response)
            let albums = []
            response.data.forEach(function (album) {
              albums.push(album)
            })
            EventBus.$emit('albums', albums)
            EventBus.$emit('artist', this.artist.name)
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
#content {
  width:100%;
}
.font-override {
  font-family: var(--font);
  font-size: var(--font-size-control);
}

#list {
  margin: 0;
  width: 100%;
}

#list-item {
  margin-top: 10px;
}
#result-line{
  text-align: left;
  padding: 1% 2%;
  cursor: pointer;
  border: 1px solid;
  width: 100%;
  border-radius: 3px;
}

#result-line:hover{
  background-color: var(--primary-sand);
}

#no-results {
  font-size: var(--font-size-control);
  padding: 2%;
  width: 100%;
  background-color: var(--primary-sand);
}
/* Icon Forward */
.hvr-icon-forward {
  float:right;
  display: inline-block;
  vertical-align: middle;
  -webkit-transform: perspective(1px) translateZ(0);
  transform: perspective(1px) translateZ(0);
  box-shadow: 0 0 1px transparent;
  position: relative;
  padding-right: 2.2em;
  -webkit-transition-duration: 0.1s;
  transition-duration: 0.1s;
}
.hvr-icon-forward:before {
  content: "\f138";
  position: absolute;
  right: 1em;
  padding: 0 1px;
  font-family: FontAwesome;
  -webkit-transform: translateZ(0);
  transform: translateZ(0);
  -webkit-transition-duration: 0.1s;
  transition-duration: 0.1s;
  -webkit-transition-property: transform;
  transition-property: transform;
  -webkit-transition-timing-function: ease-out;
  transition-timing-function: ease-out;
}
.hvr-icon-forward:hover:before, .hvr-icon-forward:focus:before, .hvr-icon-forward:active:before {
  -webkit-transform: translateX(4px);
  transform: translateX(4px);
}

/* Underline From Left */
.hvr-underline-from-left {
  display: inline-block;
  vertical-align: middle;
  -webkit-transform: perspective(1px) translateZ(0);
  transform: perspective(1px) translateZ(0);
  box-shadow: 0 0 1px transparent;
  position: relative;
  overflow: hidden;
}
.hvr-underline-from-left:before {
  content: "";
  position: absolute;
  z-index: -1;
  left: 0;
  right: 100%;
  bottom: 0;
  background: var(--primary-green);
  height: 4px;
  -webkit-transition-property: right;
  transition-property: right;
  -webkit-transition-duration: 0.3s;
  transition-duration: 0.3s;
  -webkit-transition-timing-function: ease-out;
  transition-timing-function: ease-out;
}
.hvr-underline-from-left:hover:before, .hvr-underline-from-left:focus:before, .hvr-underline-from-left:active:before {
  right: 0;
}
</style>
