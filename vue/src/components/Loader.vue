<template>
  <div v-if="size === 'small'" class="loader-small"></div>
  <div class="loader-large" v-else></div>
</template>

<script>
  export default {
    props: {
      size: ''
    }
  }
</script>

<style scoped>
.loader-large,
.loader-large:after {
  border-radius: 50%;
  width: 10em;
  height: 10em;
}
.loader-large {
  margin: 60px auto;
  font-size: 10px;
  position: relative;
  text-indent: -9999em;
  border-top: 1.1em solid var(--primary-sand);
  border-right: 1.1em solid var(--primary-sand);
  border-bottom: 1.1em solid var(--primary-sand);
  border-left: 1.1em solid var(--primary-green);
  -webkit-transform: translateZ(0);
  -ms-transform: translateZ(0);
  transform: translateZ(0);
  -webkit-animation: load8 1.1s infinite linear;
  animation: load8 1.1s infinite linear;
}

.loader-small,
.loader-small:after {
  border-radius: 50%;
  width: 3em;
  height: 3em;
}
.loader-small {
  margin: 60px auto;
  font-size: 10px;
  position: relative;
  text-indent: -9999em;
  border-top: .4em solid var(--primary-sand);
  border-right: .4em solid var(--primary-sand);
  border-bottom: .4em solid var(--primary-sand);
  border-left: .4em solid var(--primary-green);
  -webkit-transform: translateZ(0);
  -ms-transform: translateZ(0);
  transform: translateZ(0);
  -webkit-animation: load8 1.1s infinite linear;
  animation: load8 1.1s infinite linear;
}
@-webkit-keyframes load8 {
  0% {
    -webkit-transform: rotate(0deg);
    transform: rotate(0deg);
  }
  100% {
    -webkit-transform: rotate(360deg);
    transform: rotate(360deg);
  }
}
@keyframes load8 {
  0% {
    -webkit-transform: rotate(0deg);
    transform: rotate(0deg);
  }
  100% {
    -webkit-transform: rotate(360deg);
    transform: rotate(360deg);
  }
}
</style>
