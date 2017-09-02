<template>
  <div id="content">
    <p id="header-text">{{ title }}</p>
    <div v-show="isLoggedIn">
      <p>
        Logged in as {{ user.display_name === "" ? user.id : user.display_name }}
        <a v-on:click="logout">Logout</a> </p>
    </div>
  </div>
</template>

<script>
export default {
  data () {
    return {
      title: 'Discogrify',
      token: '',
      cookieName: 'auth_token',
      user: {}
    }
  },
  computed: {
    isLoggedIn: function () {
      return this.token !== null
    }
  },
  created: function () {
    this.token = document.cookie.match('(^|;)\\s*' + this.cookieName + '\\s*=\\s*([^;]+)')
    if (this.token !== null) {
      this.getUserInfo()
    }
  },
  methods: {
    logout: function () {
      document.cookie = this.cookieName + '=;expires=Thu, 01 Jan 1970 00:00:01 GMT;'
      this.token = null
      window.location.reload(true)
    },
    getUserInfo: function () {
      this.$http.get('/user').then(function (response) {
        console.log(response)
        this.user = response.data
      }).catch(function (error) {
        console.log(error)
      })
    }
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
