<template>
    <div v-bind:class="getCssClass()">
        <slot></slot>
    </div>
</template>
<script>
  export default{
    props: {
      xs: {
        type: Number,
        default: 0
      },
      sm: {
        type: Number,
        default: 0
      },
      md: {
        type: Number,
        default: 0
      },
      lg: {
        type: Number,
        default: 0
      }
    },
    methods: {
      getCssClass () {
        var props = ['xs', 'sm', 'md', 'lg']
        var classes = []
        props.forEach((item, index, array) => {
          if (this[item] > 0) {
            classes.push('col-' + item + '-' + this[item])
          }
        })
        if (classes.length === 0) {
          classes.push('col-lg-12')
        }
        return classes.join(' ')
      }
    }
  }
</script>