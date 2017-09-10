<template >
  <div id='content'>
    <div class='clickable' v-on:click='show = !show'>{{ album.name }}
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

</style>
