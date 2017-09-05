<template>
  <select class="green-control" v-model='sortOption' v-on:change='sortAlbums'>
    <option v-for='option in sortOptions' :value='option.id'>
      {{ option.text }}
    </option>
  </select>
</template>

<script>
  import EventBus from '../event-bus'

  export default {
    data () {
      return {
        sortOptions: [
          {text: 'Release date (recent first)', id: 'dateRecent'},
          {text: 'Release date (oldest first)', id: 'dateOldest'},
          {text: 'Alphabetical (A-Z)', id: 'alphaAToZ'},
          {text: 'Alphabetical (Z-A)', id: 'alphaZToA'},
          {text: 'Popularity (most first)', id: 'popularMost'},
          {text: 'Popularity (least first)', id: 'popularLeast'},
          {text: 'Number of tracks (most first)', id: 'tracksMost'},
          {text: 'Number of tracks (least first)', id: 'tracksLeast'}
        ],
        albums: []
      }
    },
    mounted () {
      EventBus.$on('albums', this.updateAlbums)
      this.sortOption = this.sortOptions[0].id
    },
    methods: {
      updateAlbums: function (albums) {
        console.log('updateAlbums')
        this.albums = albums
        this.sortAlbums()
      },
      sortAlbums: function () {
        console.log(this.albums)
        console.log(this.sortOption)
        switch (this.sortOption) {
          case 'dateRecent':
            this.albums.sort((albumA, albumB) => albumB.release_date.localeCompare(albumA.release_date))
            break
          case 'dateOldest':
            this.albums.sort((albumA, albumB) => albumA.release_date.localeCompare(albumB.release_date))
            break
          case 'alphaAToZ':
            this.albums.sort((albumA, albumB) => albumA.name.localeCompare(albumB.name))
            break
          case 'alphaZToA':
            this.albums.sort((albumA, albumB) => albumB.name.localeCompare(albumA.name))
            break
          case 'popularMost':
            this.albums.sort((albumA, albumB) => albumB.popularity - albumA.popularity)
            break
          case 'popularLeast':
            this.albums.sort((albumA, albumB) => albumA.popularity - albumB.popularity)
            break
          case 'tracksMost':
            this.albums.sort((albumA, albumB) => albumB.tracks.items.length - albumA.tracks.items.length)
            break
          case 'tracksLeast':
            this.albums.sort((albumA, albumB) => albumA.tracks.items.length - albumB.tracks.items.length)
            break
        }
        this.$emit('sort', this.albums)
      }
    }
  }

</script>

<style scoped>
select {
  font-size: var(--font-size-control);
}

</style>
