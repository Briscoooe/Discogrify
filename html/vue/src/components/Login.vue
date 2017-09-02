<template>
  <div v-if="!isLoggedIn">
    <div>
      <div>
        <button id="login-button" @click="login">
          Please log in to Spotify</button>
      </div>
    </div>
  </div>
</template>

<script>
  export default {
    data () {
      return {
        token: '',
        cookieName: 'auth_token'
      }
    },
    computed: {
      isLoggedIn: function () {
        return this.token !== null
      }
    },
    methods: {
      login: function () {
        this.$http.get('/login').then(function (response) {
          window.location.href = response.body
        }).catch(function (error) {
          console.log(error)
        })
      }
    },
    created: function () {
      this.token = document.cookie.match('(^|;)\\s*' + this.cookieName + '\\s*=\\s*([^;]+)')
    }
  }

</script>

<style scoped>
  #content {
    background-color: #1ed760;
    height:10%;
  }

  #header-text {
    font-family: 'Montserrat', sans-serif;
    font-size: 40px;
    color: #000000;
  }

</style>
