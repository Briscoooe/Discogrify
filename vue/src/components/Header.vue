<template>
  <div class="row align-center" id="content">
    <div id="inner" class="col col-6">
      <div id="header-div" class="margin2">
        <img id="logo-img" src="../assets/discogrify_full_green.png">
        <div id="logo-text">Discogrify</div>
      </div>
      <div id="account" class="margin2" v-show="isLoggedIn">
        <p>
          Logged in as {{ user.display_name === "" ? user.id : user.display_name }}
          <button class="logout" v-on:click="logout">Logout</button> </p>
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
  height:15%;
  display: flex;
  align-items: center
}

#inner {
  border-bottom: 2px solid var(--primary-sand);
}
#header-div {
  font-family: var(--font);
  font-size: var(--font-size-title);
  color: var(--primary-black);
  float:left;
  display: flex;
  align-items: center
}

#logo-img {
  height: 1.25em;
  width: 1.25em;
  float:left;
}

#logo-text {
  float:right;
  vertical-align: middle;
  margin: 0;
}

#account {
  float:right;
  vertical-align: middle;
}

.logout {
  background-color: var(--secondary-green);
  color: #fff;
}
.logout:hover, .logout:focus, .logout:active {
  background-color: var(--primary-green);
}

</style>
