<template>
  <nav class="navbar navbar-expand-lg navbar-light bg-light">
    <span class="navbar-brand mb-0 h1">ImageGallery</span>
    <ul class="navbar-nav d-flex justify-content-around w-100">
      <li class="nav-item active">
        <RouterLink class="nav-link py-4" to="/gallery">Gallery</RouterLink>
      </li>
      <li v-if="!logged" class="nav-item active">
        <RouterLink class="nav-link py-4" to="/authentication">Sign up/Log in</RouterLink>
      </li>
      <li v-if="logged" class="nav-item active">
        <RouterLink class="nav-link py-4" to="/files">Add Images</RouterLink>
      </li>
      <li v-if="logged" class="nav-item active">
        <a href=# class="nav-link py-4 end-0" @click="exit">Log out</a>
      </li>
    </ul>
  </nav>
  <main>
    <img alt="Vue logo" src="./assets/logo.png">
    <RouterView />
  </main>
</template>

<script setup>
import { logout } from './components/api';
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router'
import { useCookies } from 'vue3-cookies';

var logged = ref(false)
const router = useRouter()
const { cookies } = useCookies()

function exit(){
  logout(()=>{
    router.push("/authentication")
  });
}

function after() {
  console.log(cookies.get('username'))
  logged.value = !(!cookies.get('username'))
}

onMounted(() => {
  logged.value = !(!cookies.get('username'))
  router.afterEach(after)
})
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
