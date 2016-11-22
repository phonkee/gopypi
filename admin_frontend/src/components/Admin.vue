<template>
    <div id="wrapper">
        <navigation></navigation>
        <router-view></router-view>
        <foot></foot>
    </div>
</template>
<script>
  import Foot from './Foot.vue'
  import Navigation from './Navigation.vue'
  import store from '../store'

  export default {
    components: {Foot, Navigation},
    beforeRouteEnter (route, redirect, next) {
      store.dispatch('getAllInfo').then((info) => {
        next(vm => {
          vm.$router.push({name: 'admin.dashboard'})
        })
      }).catch((err) => {
        store.dispatch('messageError', 'Cannot fetch info: ' + err)
        next(false)
      })
    }
  }
</script>
