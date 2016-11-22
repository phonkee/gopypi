<template>
    <div v-bind:class="'panel panel-' + type">
        <div class="panel-heading" v-if="title">
            <i v-if="icon" v-bind:class="'fa fa-fw fa-' + icon"></i>
            {{title}}
        </div>
        <div class="panel-body">
            <slot></slot>
        </div>
        <div class="panel-footer" v-if="$slots.footer">
            <slot name="footer"></slot>
        </div>
    </div>
</template>
<script>
  export default {
    props: {
      footer: {
        type: String
      },
      icon: {
        type: String
      },
      title: {
        type: String
      },
      type: {
        type: String,
        default: 'default',
        validator: (value) => {
          var availableStyle = ['default', 'primary', 'success', 'info', 'warning', 'danger', 'green', 'yellow', 'red']
          if (availableStyle.indexOf(value) === -1) {
            return availableStyle[0]
          }
          return value
        }
      }
    }
  }
</script>
