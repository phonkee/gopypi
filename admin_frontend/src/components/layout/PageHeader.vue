<template>
    <div>
        <div class="row">
            <div class="col-lg-12">
                <h3 class="page-header"><slot></slot></h3>
            </div>
        </div>
        <messages></messages>
    </div>
</template>
<script>
  import Messages from '../Messages.vue'
  export default {
    components: {
      Messages
    }
  }
</script>