<template>
  <div id="app">
    <TopBar />
    <router-view />
  </div>
</template>

<script>
import TopBar from './components/Topbar.vue';

export default {
  name: 'App',
  components: { TopBar },
};
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}
</style>
