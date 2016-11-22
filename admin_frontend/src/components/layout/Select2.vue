<template>
    <select name="" id="" multiple="multiple" style="width: 100%">
        <option :value="option.id" v-for="option in options" :selected="isSelected(option.id)">{{label(option)}}</option>
    </select>
</template>
<script>
  import $ from 'jquery'
  export default {
    props: {
      options: {
        type: Array
      },
      selected: {
        type: Array
      },
      label: {
        type: Function,
        default: 'label'
      }
    },
    mounted () {
      var vm = this
      $(this.$el).select2({
        data: vm.options,
        closeOnSelect: true
      }).on('select2:unselect', (event) => {
        setTimeout(() => {
          $(vm.$el).select2('close')
        }, 1)
        vm.$emit('remove', event.params.data)
      }).on('select2:select', (event) => {
        $(vm.$el).select2('close')
        vm.$emit('add', event.params.data)
      })
    },
    $destroy () {
      $(this.$el).select2('destroy')
    },
    methods: {
      isSelected (id) {
        for (var item of this.selected) {
          if (item.id === id) {
            return true
          }
        }
        return false
      }
    }
  }
</script>