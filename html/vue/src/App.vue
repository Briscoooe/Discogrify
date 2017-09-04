<template>
  <div id="app">
    <custom-header></custom-header>
    <login></login>
    <search id="search" v-on:scroll="scroll('results')"></search>
    <sort></sort>
    <results id="results" v-on:scroll="scroll('search')"></results>
  </div>
</template>

<script>
import Header from './components/Header'
import Login from './components/Login'
import Search from './components/Search'
import Sort from './components/Sort'
import Results from './components/Results'
import Jump from '../node_modules/jump.js'

export default {
  name: 'app',
  components: {
    'custom-header': Header,
    'login': Login,
    'search': Search,
    'sort': Sort,
    'results': Results
  },
  methods: {
    scroll: function (element) {
      Jump('#' + element, {
        duration: 400
      })
    }
  }
}
</script>

<style>
  @import url('https://fonts.googleapis.com/css?family=Montserrat');
  @import url('https://cdnjs.cloudflare.com/ajax/libs/kube/6.5.2/css/kube.min.css');
  @import url('https://maxcdn.bootstrapcdn.com/font-awesome/4.7.0/css/font-awesome.min.css');
</style>

<style>
:root {
  --primary-green: #1ED763;
  --primary-grey: #828282;
  --primary-sand: #ECEBE8;
  --primary-black: #0d0d0e;
  --secondary-green: #009C3A;
  --font: 'Montserrat', sans-serif;
  --font-size-title: 2.5em;
  --font-size-control: 1.25em;
  --font-size-data: 1em;
}
#app {
  font-family: var(--font);
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: var(--primary-black);
}
a:hover {
  cursor: pointer;
  cursor: hand;
}
button {
  background-color: var(--secondary-green);
  font-size: var(--font-size-control);
  font-family: var(--font);
}

button:hover {
  background-color: var(--primary-green);
}
li {
  list-style-type: none;
}
</style>

<style scoped>
  #search{
    margin-top: 5%;
  }
</style>
