<template>
  <div class="row align-center" id="content">
    <div class="col col-6">
      <p id="header-text">Discogrify</p>
      <div id="account" v-show="isLoggedIn">
        <p>
          Logged in as {{ user.display_name === "" ? user.id : user.display_name }}
          <a id="logout" v-on:click="logout">Logout</a> </p>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data () {
    return {
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
        this.user = response.data
        console.log(response)
      }).catch(function (error) {
        console.log(error)
      })
    }
  }
}
</script>

<style scoped>
#content {
  background-color: var(--primary-green);
  height:15%;
  padding:1%;
}

#header-text {
  font-family: var(--font);
  font-size: var(--font-size-title);
  color: var(--primary-black);
  float:left;
}

#account {
  float:right;
}

#logout {
  color: var(--primary-sand);
}

</style>
