<script>
import { ref } from 'vue'
import { getImagesByUser } from './api'

export default {
  name: 'ViewImages',
  setup(){
    const user = ref('')
    const images = ref([])
    images.value.push([])
    
    return {
      user, images
    }
  },
  methods: {
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
    },
    submit: function() {
      getImagesByUser(this.user, this.setImages);
    }
  }
}

</script>

<template>
  <form class="card" @submit.prevent="submit">
    <div class="form-floating my-2">
      <input id="user" class="form-control" v-model="user" type="text" placeholder="example_username">
      <label for="user">Name of User</label>
    </div>
    <div class="my-2">
      <button type="submit">Search</button>
    </div>
  </form>
  <div class="card">
    <div v-if="images[0].length == 0">No Images</div>
    <div v-for="imagerow in images" :key="imagerow" class="row">
      <template v-for="image in imagerow" :key="image">
        <div class="col border">
          <img :src="image.filepath" :alt="image.filename"/>
        </div>
      </template>
    </div>
  </div>
</template>