<template>
  <div id="app">
    <custom-header></custom-header>
    <login></login>
    <search v-on:albums="albums"></search>
    <sort></sort>
  </div>
</template>

<script>
import Header from './components/Header'
import Login from './components/Login'
import Search from './components/Search'
import Sort from './components/Sort'

export default {
  name: 'app',
  components: {
    'custom-header': Header,
    'login': Login,
    'search': Search,
    'sort': Sort
  },
  data: function () {
    return {
      albums: []
    }
  },
  methods: {
    albums: function () {
      console.log(this.albums)
    }
  }
}
</script>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}
</style>
