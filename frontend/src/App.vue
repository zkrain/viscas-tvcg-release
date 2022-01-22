<template>
  <div id="app">
    <Main />
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
// import HelloWorld from './components/HelloWorld.vue';
// import MapView from './components/Map/MapView.vue'
import Main from './components/Main.vue'

@Component({
  components: {
    // MapView,
    Main
  }
})
export default class App extends Vue {}
</script>

<style lang="scss">
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  width: 1920px;
  // width: 1700px;
  height: 1080px;
}

body {
  margin: 0;
}
::-webkit-scrollbar {
  display: none;
}
</style>
