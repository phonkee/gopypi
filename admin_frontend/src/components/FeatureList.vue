<template>
    <page title="Features">
        <row>
            <column>
                <panel title="Features" icon="cogs">
                    <div class="checkbox" v-for="(feature, index) in info.features">
                        <label>
                            <input type="checkbox" :checked="feature.value" @change="toggleFeature(index, feature.id)">{{feature.description}}
                        </label>
                    </div>
                </panel>
            </column>
        </row>
    </page>
</template>
<script>
  import { mapGetters } from 'vuex'
  import Page from './layout/Page.vue'
  import Panel from './layout/Panel.vue'
  import Column from './layout/Column.vue'
  import Row from './layout/Row.vue'
  import store from '../store'
  import feature from '../api/feature'
  export default {
    computed: mapGetters({
      info: 'allInfo'
    }),
    components: {
      Page, Panel, Column, Row
    },
    beforeRouteEnter (route, redirect, next) {
      store.dispatch('getAllInfo').then((info) => {
        next()
      }).catch((err) => {
        store.dispatch('messageError', 'Cannot fetch info: ' + err)
        next(false)
      })
    },
    methods: {
      toggleFeature (index, id) {
        var newValue = !this.info.features[index].value
        feature.updateFeature(id, newValue).then(() => {
          var message = ''
          if (newValue) {
            message = 'Enabled feature: ' + this.info.features[index].description
          } else {
            message = 'Disabled feature: ' + this.info.features[index].description
          }
          store.dispatch('messageSuccess', message)
          store.dispatch('getAllInfo')
        })
      }
    }
  }
</script>