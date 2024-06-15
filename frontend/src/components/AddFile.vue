<script>
import { ref } from 'vue'
import { addFile, addLink, getImagesByUser, removeImageByUser } from './api'

export default {
  name: 'AddFile',
  setup(){
    const type = ref('Link')
    const username = ref('')
    const file = ref(null)
    const link = ref('')
    const images = ref([])
    return {
      type, file, link, images, username
    }
  },
  mounted() {
    if(!this.$cookies.get("username")) {
      this.$router.push({ path: '/authentication'})
    }
    this.username = this.$cookies.get("username")
    getImagesByUser(this.username, this.setImages)
  },
  methods: {
    setFile: function(event) {
      this.file = event.target.files[0]
      this.link = ''
    },
    setLink: function(event) {
      this.link = event.target.value
      this.file = null
    },
    sendFile: function() {
      const self = this
      if(this.file != null) {
        addFile(this.file, self.username, function(){
          getImagesByUser(self.username, self.setImages);
        })
      } else if (this.link != '') {
        addLink(this.link, self.username, function(){
          getImagesByUser(self.username, self.setImages);
        })
      } else {
        console.warn("No image selected")
      }
    },
    setImages: function(img) {
      let temparr = img.data;
      this.images = [];
      if(!temparr) {
        return
      }
      while (temparr.length > 0) {
        let tt = temparr.splice(0, 4);
        this.images.push(tt);
      }
      console.log(this.images);
    },
    removeImage: function(image) {
      const self = this
      removeImageByUser(self.username, image, function(result){
        getImagesByUser(self.username, self.setImages)
        console.log(result);
      });
    }
  }
}

</script>

<template>
  <form class="card m-2" @submit.prevent="submit">
    <div class="my-2">
      <select @input="type = $event.target.value">
        <option>Link</option>
        <option>File</option>
      </select>
    </div>
    <div class="m-2">
      <input v-if="type==='File'" class="form-control" @change="setFile" type="file" accept="image/png, image/jpeg"/>
      <input v-else type="text" id="text-input" class="form-control" :value="link" @input="setLink" placeholder="paste link to image here"/>
      <label for="text-input"></label>
    </div>
    <div class="my-2">
      <button class="btn btn-primary" type="submit" @click="sendFile($event)">Add image</button>
    </div>
  </form>
  <div class="card m-2">
    <div v-for="imagerow in images" :key="imagerow" class="row mx-2">
      <template v-for="image in imagerow" :key="image">
          <div class="vstack border" style="max-width:25%">
            <img :id="image.id" :src="image.filepath" :alt="image.filename"/>
            <button class="btn btn-danger my-1" id="removeImage" @click="removeImage(image.id)">Remove Image</button>
          </div>
      </template>
    </div>
  </div>
</template>