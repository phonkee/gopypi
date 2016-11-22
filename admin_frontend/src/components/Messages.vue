<template>
    <transition-group name="fade">
        <div v-bind:class="'alert alert-dismissable alert-' + message.level" v-for="(message, key)  in messages" :key=message.guid>
            <button type="button" class="close" aria-hidden="true" v-on:click.prevent="removeMessage(message.guid)">×</button>
            <i v-bind:class="getIcon(message)"></i> {{message.message}}
        </div>
    </transition-group>
</template>
<script>
  import {mapGetters} from 'vuex'
  export default {
    computed: mapGetters({
      messages: 'messages'
    }),
    methods: {
      removeMessage (guid) {
        this.$store.dispatch('hideMessage', guid)
      },
      getIcon (message) {
        var icons = {
          success: 'check',
          warning: 'exclamation-triangle',
          error: 'exclamation-triangle',
          info: 'info'
        }

        return 'fa fa-fw fa-' + icons[message.level]
      }
    }
  }
</script>
<style scoped>
    .fade-leave-active {
        transition: all .5s ease;
        overflow: hidden;
    }
    .fade-leave-active {
        /*height: 0;*/
        /*padding: 0;*/
        opacity: 0;
    }
</style>