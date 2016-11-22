<template>
    <span>
        <i class="fa fa-check" v-if="value" v-bind:title="title"></i>
        <i class="fa fa-close" v-if="!value"></i>
    </span>
</template>
<script>
  export default {
    props: {
      title: {
        type: String,
        default: ''
      },
      value: {
        type: Boolean,
        default: false
      }
    }
  }
</script>
