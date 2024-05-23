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
body {
  margin: 0;
  height: calc(100vh - 10px);
}
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
  padding: 5px;

}

.dragging {
  background-color: #eee;
  border: 3px dashed #aaa;
}
</style>
