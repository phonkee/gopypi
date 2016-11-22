<template>
    <div v-bind:class="getCssClass()">
        <slot></slot>
    </div>
</template>
<script>
  export default {
    props: ['lg'],
    methods: {
      getCssClass () {
        var classes = ['well']
        if (this.lg) {
          classes.push('well-lg')
        }
        return classes.join(' ')
      }
    }
  }
</script>
