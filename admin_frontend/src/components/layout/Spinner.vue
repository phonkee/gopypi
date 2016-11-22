<template>
    <moon-loader :loading="spinnerPendingRequests > 0" color="#337ab7" size="14px" margin="4px"></moon-loader>
    {{spinnerPendingRequests}}
</template>
<script>
  import MoonLoader from 'vue-spinner/src/MoonLoader.vue'
  import { mapGetters } from 'vuex'
  export default {
    components: {
      MoonLoader
    },
    computed: mapGetters({
      spinnerPendingRequests: 'spinnerPendingRequests'
    })
  }
</script>
