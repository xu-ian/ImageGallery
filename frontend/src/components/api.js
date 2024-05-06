import axios from 'axios';

const SERVER = 'http://localhost:8000'

export function signup(username, password, next, err) {
  axios.post(SERVER + "/signup", {"username": username, "password": password}, { headers: { "Content-Type": "application/json"}})
  .then(next)
  .catch(err)
}

export function login(username, password, next, err) {
  axios.post(SERVER + "/login", {"username": username, "password": password}, { withCredentials: true, headers: { "Content-Type": "application/json"}})
  .then(next)
  .catch(err)
}

export function logout(next) {
  axios.post(SERVER + "/logout", {}, { withCredentials: true})
  .then(next)
  .catch(function(err) {
    console.log(err)
  })
}

export function addFile(file, username, next) {
  let form = new FormData()
  form.append("file", file)
  axios.post(SERVER + '/'+username+'/images', form, { withCredentials: true, headers: { "Content-Type": "multipart/form-data" }})
  .then(next)
  .catch(function(err) {
    console.log(err)
  })
}

export function addLink(link, username, next) {

  axios.post(SERVER + '/'+username+'/link', {"link": link}, { withCredentials: true, headers: { "Content-Type": "application/json"}})
  .then(next)
  .catch(function(err) {
    console.log(err)
  })
}

export function getImagesByUser(user, next) {
  axios.get(SERVER + '/'+user+'/images')
  .then(next)
  .catch(function(err) {
    console.log(err)
  })
}

export function removeImageByUser(user, filepath, next) {
  axios.delete(SERVER + '/'+user+'/images/' + filepath, { withCredentials: true})
  .then(next)
  .catch(function(err) {
    console.log(err)
  })
}