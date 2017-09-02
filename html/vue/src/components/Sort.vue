<template>
  <select v-model='sortOption' v-on:change='sortAlbums'>
    <option disabled value=''>Sort by...</option>
    <option v-for='option in sortOptions' :value='option.id'>
      {{ option.text }}
    </option>
  </select>
</template>

<script>
  export default {
    data () {
      return {
        sortOptions: [
          {text: 'Alphabetical (A-Z)', id: 'alphaAToZ'},
          {text: 'Alphabetical (Z-A)', id: 'alphaZToA'},
          {text: 'Popularity (most first)', id: 'popularMost'},
          {text: 'Popularity (least first)', id: 'popularLeast'},
          {text: 'Release date (oldest first)', id: 'dateOldest'},
          {text: 'Release date (recent first)', id: 'dateRecent'},
          {text: 'Number of tracks (most first)', id: 'tracksMost'},
          {text: 'Number of tracks (least first)', id: 'tracksLeast'}
        ]
      }
    },
    methods: {
      sortAlbums: function () {
        switch (this.sortOption) {
          case 'alphaAToZ':
            this.allAlbums.sort((albumA, albumB) => albumA.name.localeCompare(albumB.name))
            break
          case 'alphaZToA':
            this.allAlbums.sort((albumA, albumB) => albumB.name.localeCompare(albumA.name))
            break
          case 'popularMost':
            this.allAlbums.sort((albumA, albumB) => albumB.popularity - albumA.popularity)
            break
          case 'popularLeast':
            this.allAlbums.sort((albumA, albumB) => albumA.popularity - albumB.popularity)
            break
          case 'dateOldest':
            this.allAlbums.sort((albumA, albumB) => albumA.release_date.localeCompare(albumB.release_date))
            break
          case 'dateRecent':
            this.allAlbums.sort((albumA, albumB) => albumB.release_date.localeCompare(albumA.release_date))
            break
          case 'tracksMost':
            this.allAlbums.sort((albumA, albumB) => albumB.tracks.items.length - albumA.tracks.items.length)
            break
          case 'tracksLeast':
            this.allAlbums.sort((albumA, albumB) => albumA.tracks.items.length - albumB.tracks.items.length)
            break
        }
      }
    }
  }

</script>

<style scoped>


</style>
