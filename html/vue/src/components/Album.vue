<template >
  <div>
    <div v-on:click='show = !show'>{{ album.name }}
      <b>({{ album.tracks.items.length }} {{ album.tracks.items.length == 1 ? 'track' : 'tracks' }})</b>
    </div>
    <transition name='fade'>
      <div v-show='show'>
        <ul>
          <li v-for='track in album.tracks.items'>
            <input type='checkbox'
                   checked='checkedTracks.contains(track.id)'
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
        show: false,
        checked: true
      }
    },
    computed: {
      allTracksChecked: function () {
        return this.checkedTracks.length === this.album.tracks.items.length
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
      updateAlbum: function (albumId) {
        let self = this
        if (self.album.id === albumId) {
          if (this.allTracksChecked) {
            self.checked = true
          }
          self.album.tracks.items.forEach(function (track) {
            let index = self.checkedTracks.indexOf(track.id)
            if (self.checked) {
              if (index < 0) {
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
        if (this.allTracksChecked) {
          this.$emit('update-album', this.album.id, true)
        } else {
          this.$emit('update-album', this.album.id, false)
        }
      }
    }
  }
</script>

<style></style>
