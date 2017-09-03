<template>
  <div id="app">
    <custom-header></custom-header>
    <login></login>
    <search></search>
    <sort></sort>
    <results></results>
  </div>
</template>

<script>
import Header from './components/Header'
import Login from './components/Login'
import Search from './components/Search'
import Sort from './components/Sort'
import Results from './components/Results'

export default {
  name: 'app',
  components: {
    'custom-header': Header,
    'login': Login,
    'search': Search,
    'sort': Sort,
    'results': Results
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
}
#app {
  font-family: 'Montserrat', sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}
a:hover {
  cursor: pointer;
  cursor: hand;
}

</style>
