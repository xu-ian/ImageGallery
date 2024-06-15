
<script>
import { ref } from 'vue'
import { signup, login, logout } from './api'

export default {
  name: 'AuthenticateUser',
  setup(){
    const username = ref('')
    const password = ref('')
    const signupstate = ref("nav-link")
    const loginstate = ref("nav-link active")
    const logoutstate = ref("nav-link")
    const state = ref(1)
    const errorMessage = ref('')
    return {
      username, password, state, signupstate, loginstate, logoutstate, errorMessage
    }
  },
  methods: {
    setState: function(state) {
      this.errorMessage = ""
      this.username = ""
      this.password = ""
      this.state = state
      if(state == 0) {
        this.signupstate = "nav-link active"
        this.loginstate = "nav-link"
        this.logoutstate = "nav-link"
      } else if(state == 1) {
        this.signupstate = "nav-link"
        this.loginstate = "nav-link active"
        this.logoutstate = "nav-link"
      } else {
        this.signupstate = "nav-link"
        this.loginstate = "nav-link"
        this.logoutstate = "nav-link active"
      }
    },
    signup: function(){
      const self = this
      signup(this.username, this.password, 
      function() {
        self.setState(1)        
      },
      function(error) {
        self.errorMessage = error.response.data
        console.log(error)
      })
    },
    login: function(){
      const self = this
      login(this.username, this.password, 
      function() {
        self.$router.push({ path: '/files'})
      },
      function(error) {
        console.log(error)
        self.errorMessage = error.response.data
      })
    },
    signout: function(){
      const self = this
      logout(function(){
        self.setState(1)
      })
    }
  }
}

</script>

<template>
  <div class="card text-center m-2">
  <div class="card-header">
    <ul class="nav nav-tabs card-header-tabs">
      <li v-if="state != 2" class="nav-item">
        <a :class="signupstate" href="#" @click="setState(0)">Sign up</a>
      </li>
      <li v-if="state != 2" class="nav-item">
        <a :class="loginstate" href="#" @click="setState(1)">Log in</a>
      </li>
      <li v-if="state == 2" class="nav-item">
        <a :class="logoutstate" href="#" @click="setState(2)">Log out</a>
      </li>
    </ul>
  </div>
  <form v-if="state==0" class="card-body" @submit.prevent="signup">
    <div v-if="errorMessage != ''" class="alert alert-danger" role="alert">{{ errorMessage }}</div>
    <div class="my-2">
      <label class="m-2" for="username">Enter your username: </label>
      <input class="m-2" v-model="username" id="username" type="text" />
    </div>
    <div class="my-2">
      <label class="m-2" for="password">Enter your password: </label>
      <input v-model="password" class="m-2" id="password" type="password" />
    </div>
    <div class="my-2">
      <button class="btn btn-primary" type="submit">Sign up</button>
    </div>
  </form>
  <form v-else-if="state==1" class="card-body" @submit.prevent="login">
    <div v-if="errorMessage != ''" class="alert alert-danger" role="alert">{{ errorMessage }}</div>
    <div class="my-2">
      <label class="m-2" for="username">Enter your username: </label>
      <input class="m-2" v-model="username" id="username" type="text" />
    </div>
    <div class="my-2">
      <label class="m-2" for="password">Enter your password: </label>
      <input class="m-2" v-model="password" id="password" type="password" />
    </div>
    <div class="my-2">
      <button class="btn btn-primary" type="submit">Log in</button>
    </div>
  </form>
  <form v-else class="card-body" @submit.prevent="signout">
    <div class="my-2">
      <button class="btn btn-primary" type="submit">Log out</button>
    </div>
  </form>
</div>
</template>