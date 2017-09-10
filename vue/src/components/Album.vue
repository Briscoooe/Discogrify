<template >
  <div id='content'>
    <div class='clickable' v-on:click='show = !show'>
      <i aria-hidden="false" class="fa fa-caret-down" ></i>{{ album.name }}
    </div>
    <transition name='fade'>
      <div v-show='show'>
        <ul>
          <li v-for='track in album.tracks.items'>
            <input id='checkbox'
                   type='checkbox'
                   :checked='checkedTracks.indexOf(track.id) > -1'
                   :value='track.id'
                   v-model='checkedTracks'
                   v-on:change='updateTrack($event.target.value)'>
            {{ track.name }}
          </li>
        </ul>
      </div>
    </transition>
  </div>
</template>

<script>
  import EventBus from '../event-bus'

  export default {
    props: {
      album: Object
    },
    data: function () {
      return {
        checkedTracks: [],
        show: false
      }
    },
    mounted () {
      EventBus.$on('toggle-album', this.updateAlbum)
    },
    created: function () {
      let self = this
      self.album.tracks.items.forEach(function (track) {
        self.checkedTracks.push(track.id)
      })
    },
    methods: {
      updateAlbum: function (albumId, isChecked) {
        let self = this
        if (self.album.id === albumId) {
          self.album.tracks.items.forEach(function (track) {
            let index = self.checkedTracks.indexOf(track.id)
            if (!isChecked) {
              if (index === -1) {
                self.checkedTracks.push(track.id)
                self.updateTrack(track.id)
              }
            } else {
              if (index > -1) {
                self.checkedTracks.splice(index, 1)
                self.updateTrack(track.id)
              }
            }
          })
        }
      },
      updateTrack: function (trackId) {
        this.$emit('update-track', trackId)
        if (this.checkedTracks.length === this.album.tracks.items.length ||
            this.checkedTracks.length === 0) {
          this.$emit('update-album', this.album.id)
        }
      }
    }
  }
</script>

<style scoped>
#content {
  text-align: left;
}
.clickable:hover {
  cursor: pointer;
}

/* Icon Sink */
.hvr-icon-sink {
  display: inline-block;
  vertical-align: middle;
  -webkit-transform: perspective(1px) translateZ(0);
  transform: perspective(1px) translateZ(0);
  box-shadow: 0 0 1px transparent;
  position: relative;
  padding-right: 2.2em;
  -webkit-transition-duration: 0.3s;
  transition-duration: 0.3s;
}
.hvr-icon-sink:before {
  content: "\f01a";
  position: absolute;
  right: 1em;
  padding: 0 1px;
  font-family: FontAwesome;
  -webkit-transform: translateZ(0);
  transform: translateZ(0);
  -webkit-transition-duration: 0.3s;
  transition-duration: 0.3s;
  -webkit-transition-property: transform;
  transition-property: transform;
  -webkit-transition-timing-function: ease-out;
  transition-timing-function: ease-out;
}
.hvr-icon-sink:hover:before, .hvr-icon-sink:focus:before, .hvr-icon-sink:active:before {
  -webkit-transform: translateY(4px);
  transform: translateY(4px);
}

/* Icon Float */
.hvr-icon-float {
  display: inline-block;
  vertical-align: middle;
  -webkit-transform: perspective(1px) translateZ(0);
  transform: perspective(1px) translateZ(0);
  box-shadow: 0 0 1px transparent;
  position: relative;
  padding-right: 2.2em;
  -webkit-transition-duration: 0.3s;
  transition-duration: 0.3s;
}
.hvr-icon-float:before {
  content: "\f01b";
  position: absolute;
  right: 1em;
  padding: 0 1px;
  font-family: FontAwesome;
  -webkit-transform: translateZ(0);
  transform: translateZ(0);
  -webkit-transition-duration: 0.3s;
  transition-duration: 0.3s;
  -webkit-transition-property: transform;
  transition-property: transform;
  -webkit-transition-timing-function: ease-out;
  transition-timing-function: ease-out;
}
.hvr-icon-float:hover:before, .hvr-icon-float:focus:before, .hvr-icon-float:active:before {
  -webkit-transform: translateY(-4px);
  transform: translateY(-4px);
}
</style>
