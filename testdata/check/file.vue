<!--
 Copyright 2025 BINARY Members

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
-->

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