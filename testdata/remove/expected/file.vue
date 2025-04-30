<template>
  <div class="hello-world">
    <h1>{{ message }}</h1>
    <p>Welcome to your Vue.js application</p>
  </div>
</template>

<script>
export default {
  name: 'HelloWorld',
  data() {
    return {
      message: 'Hello, World!'
    }
  }
}
</script>

<style scoped>
.hello-world {
  margin-top: 20px;
  text-align: center;
}
</style>