<template>
  <div id="content">
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
</style>
