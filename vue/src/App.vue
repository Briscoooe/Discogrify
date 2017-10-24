<template>
  <div id="app">
    <custom-header></custom-header>
    <div style="flex: 1">
      <login v-if="!loggedIn"></login>
      <search v-if="loggedIn" id="search" v-on:scroll="scroll('results')"></search>
      <results v-if="loggedIn" id="results" v-on:scroll="scroll('search')"></results>
    </div>
    <custom-footer></custom-footer>
  </div>
</template>

<script>
import Header from './components/Header'
import Login from './components/Login'
import Search from './components/Search'
import Results from './components/Results'
import Footer from './components/Footer'
import Jump from '../node_modules/jump.js'

export default {
  name: 'app',
  components: {
    'custom-header': Header,
    'login': Login,
    'search': Search,
    'results': Results,
    'custom-footer': Footer
  },
  data () {
    return {
      cookieName: 'auth_token'
    }
  },
  computed: {
    loggedIn: function () {
      let token = document.cookie.match('(^|;)\\s*' + this.cookieName + '\\s*=\\s*([^;]+)')
      return token !== null
    }
  },
  methods: {
    scroll: function (element) {
      Jump('#' + element, {
        duration: 1000
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
  --font-size-heading: 2em;
  --font-size-control: 1.25em;
  --font-size-data: 1em;
}
* {
  margin: 0;
  padding: 0;
}
#app {
  font-family: var(--font);
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: var(--primary-black);
  display: flex;
  min-height: 100vh;
  flex-direction: column;
  padding-left:2%;
  padding-right:2%;
}

.margin1 {
  margin-top: 1%;
  margin-bottom: 1%;
}
.margin2 {
  margin-top: 2%;
  margin-bottom: 2%;
}

a:hover {
  cursor: pointer;
}

/* Fade */
button{
  background-color: var(--secondary-green);
  font-size: var(--font-size-control);
  font-family: var(--font);
  display: inline-block;
  vertical-align: middle;
  -webkit-transform: perspective(1px) translateZ(0);
  transform: perspective(1px) translateZ(0);
  box-shadow: 0 0 1px transparent;
  overflow: hidden;
  -webkit-transition-duration: 0.3s;
  transition-duration: 0.3s;
  -webkit-transition-property: color, background-color;
  transition-property: color, background-color;
}
button:hover, button:focus, button:active {
  background-color: var(--primary-green);
  color: var(--primary-sand);
}
li {
  list-style-type: none;
}
input[type="checkbox"]:hover {
  cursor: pointer;
}
.green-control:focus{
  outline: none;
  box-shadow: 0 0 5px var(--primary-green);
  border:1px solid var(--secondary-green);
}

.green-control:hover{
  border: 1px solid var(--primary-grey);
  border-radius: 5px;
}

.green-control:focus:hover{
  outline: none;
  box-shadow: 0 0 5px var(--primary-green);
  border:1px solid var(--secondary-green);
  border-radius:0;

}
.fade-enter-active, .fade-leave-active {
  transition: opacity .3s
}
.fade-enter, .fade-leave-to {
  opacity: 0
}
</style>
