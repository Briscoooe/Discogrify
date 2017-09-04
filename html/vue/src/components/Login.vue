<template>
  <div id="content" v-if="!isLoggedIn">
    <div>
      <div>
        <button id="login-button" @click="login">
          <i aria-hidden="false" class="fa fa-spotify"></i>
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
  padding-top:5%;
  font-size: var(--font-size-title);
}
</style>
