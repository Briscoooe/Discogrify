<template>
  <div id="content">
    <div class="row align-center" id="info-message">
      Discogrify is a simple service that allows you to catalog an artist's entire Spotify discography into a single playlist<br>
      There are only three simple steps:
    </div>
    <div class="row align-center" id="list">
      <ul>
        <li>
          Login
        </li>
        <li>
          Search for an artist
        </li>
        <li>
          Select an artist
        </li>
        <li>
          Publish!
        </li>
      </ul>
      There are also some optional steps you can take before publishing
      <ul>
        <li>
          Add or remove albums/tracks from the playlist
        </li>
        <li>
          Add multiple artists to the playlist
        </li>
      </ul>
    </div>
    <div class="row align-center">
      <button id="login-button" @click="login">
        <i aria-hidden="false" class="fa fa-spotify hvr"></i>
        Log in to Spotify</button>
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

#list {
  float:left;
}
</style>
