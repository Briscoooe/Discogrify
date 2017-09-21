<template>
  <div id="content">
    <div id="info-message">
      <!--Discogrify is a simple service that allows you to catalog an artist's entire Spotify discography in a few simple
      clicks. Once logged in, simply search for the artist-->
    </div>
    <div>
      <button id="login-button" @click="login">
        <i aria-hidden="false" class="fa fa-spotify hvr"></i>
        Please log in to Spotify</button>
    </div>
  </div>
</template>

<script>
  export default {
    methods: {
      login: function () {
        this.$http.get('/login').then(function (response) {
          window.location.href = response.body
        }).catch(function (error) {
          console.log(error)
        })
      }
    }
  }

</script>

<style scoped>
#content {
  padding-top:5%;
  font-size: var(--font-size-control);
}
#info-message {

}
</style>
