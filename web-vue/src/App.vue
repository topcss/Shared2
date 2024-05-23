<template>
  <div id="app">
    <FileList msg="hello world" />
  </div>
</template>

<script>
import FileList from './components/FileList.vue'

export default {
  name: 'App',
  components: {
    FileList: FileList
  }
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
  margin-top: 10px;
  padding: 5px;
}
</style>
