<template>
    <div id="page-wrapper">
        <div class="row">
            <div class="col-lg-12">
                <h3 class="page-header">{{title}}</h3>
            </div>
            <!-- /.col-lg-12 -->
        </div>
        <messages></messages>
        <slot></slot>
    </div>
</template>
<script>
  import Messages from '../Messages.vue'
  export default {
    components: {
      Messages
    },
    props: {
      title: {
        type: String,
        default: 'Title'
      }
    }
  }
</script>
