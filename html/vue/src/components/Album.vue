<template id="album-col">
  <div>
    <p v-on:click="show = !show">X</p>
    <div>{{ album.name }}
      <b>({{ album.tracks.items.length }} {{ album.tracks.items.length == 1 ? 'track' : 'tracks' }})</b>
    </div>
    <transition name="fade">
      <div v-show="show">
        <ul>
          <li v-for="track in album.tracks.items">
            <input type="checkbox"
                   checked="checkedTracks.contains(track.id)"
                   :value="track.id"
                   v-model="checkedTracks"
                   v-on:change="updateTrack($event.target.value)">
            {{ track.name }}
          </li>
        </ul>
      </div>
    </transition>
  </div>
</template>

<script>
  export default {
    props : {
      album: Object
    },
    data: function () {
      return {
        checkedTracks: [],
        show: false,
        checked: true
      }
    },
    computed:  {
      allTracksChecked: function () {
        return this.checkedTracks.length === this.album.tracks.items.length
      }
    },
    created: function() {
      let self = this;
      self.album.tracks.items.forEach(function (track) {
          self.checkedTracks.push(track.id)
      })
    },
    methods: {
      updateAlbum: function () {
        let self = this;
        if(self.allTracksChecked) {
          self.checked = true;
        }
        self.album.tracks.items.forEach(function (track) {
          let index = self.checkedTracks.indexOf(track.id);
          if(self.checked) {
            if(index < 0) {
              self.checkedTracks.push(track.id);
              self.updateTrack(track.id)
            }
          }
          else {
            if(index > -1) {
              self.checkedTracks.splice(index, 1);
              self.updateTrack(track.id)
            }
          }
        });
      },
      updateTrack: function(trackId) {
        this.$emit('update', trackId);
      }
    }
  }
</script>
