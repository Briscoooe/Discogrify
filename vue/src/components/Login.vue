<template>
  <div id="content" class="row align-center">
    <div class="col col-6">
      <div class="row margin2" id="info-message-div">
          Create discographies on Spotify in a matter of clicks!
      </div>
      <div class="row margin2 info-content" id="site-info-div">
        <div class="col col-4 margin1">
          <b class="info-header">Easy to use</b>
          <br>
          Create discographies in 3 simple steps:
          <ul>
            <li>
              Login into the site using the button below
            </li>
            <li>
              Search for an artist
            </li>
            <li>
              Publish the playlist to your own account
            </li>
            <i><b>Note</b>: For duplicates in playlists, see <a href="https://jmperezperez.com/spotify-dedup/">Spotify dedup</a></i>
          </ul>
        </div>
        <div class="col col-4 margin1">
          <b class="info-header">How?</b>
          <br>
          Using the Spotify API, I was able to retrieve each track containing any given artist and present them in
          a single list which you can alter and publish as a playlist.
          If you are interested in specifics, check out the source code on
          <a href="https://github.com/Briscoooe/Discogrify/">GitHub</a>. Any contributions welcome!
        </div>
        <div class="col col-4 margin1">
          <b class="info-header">Suggestions</b>
          <br>
          If you have any feature requests, notice any bugs or anything of the sort please
          <a href="https://github.com/Briscoooe/Discogrify/issues/new">open an issue on GitHub</a> or
          <a href="https://www.reddit.com/user/briscoooe">send me a message on Reddit</a>
        </div>
      </div>
      <div class="row align-center margin2" id="login-div">
        <button class="col" id="login-button" @click="login">
          <i aria-hidden="false" class="fa fa-spotify hvr"></i>
          Please log in to Spotify</button>
      </div>
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

.col-4 {
  padding: 1%;
}
#content {
}

#info-message-div {
  line-height: 1em;
  font-size: var(--font-size-title);
}

#site-info-div {
  text-align: left;
}

.info-header {
  font-size: var(--font-size-control);
}

.info-content {
  font-size: var(--font-size-data)
}
li {
  list-style-type: disc;
}

#list {
  float:left;
}
</style>
