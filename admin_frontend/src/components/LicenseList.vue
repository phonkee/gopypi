<template>
    <page title="Licenses">
        <row>
            <column :lg="12">
                <panel title="Licenses">
                    <div class="table-responsive">
                        <table class="table table-striped table-bordered table-hover">
                            <thead>
                            <tr>
                                <th>#</th>
                                <th>Code</th>
                                <th>Name</th>
                            </tr>
                            </thead>
                            <tbody>
                                <tr v-for="license in licenses">
                                    <td>{{ license.id }}</td>
                                    <td>{{ license.code }}</td>
                                    <td>{{ license.name }}</td>
                                </tr>
                            </tbody>
                        </table>
                    </div>
                </panel>
            </column>
        </row>
    </page>
</template>
<script>
  import {mapGetters} from 'vuex'
  import Page from './layout/Page.vue'
  import Panel from './layout/Panel.vue'
  import Column from './layout/Column.vue'
  import Row from './layout/Row.vue'
  import TrueFalseIcon from './layout/TrueFalseIcon.vue'

  import store from '../store'

  export default {
    components: {Column, Row, Panel, Page, TrueFalseIcon},
    computed: mapGetters({
      licenses: 'allLicenses'
    }),
    beforeRouteEnter (route, redirect, next) {
      store.dispatch('getAllLicenses').then((licenses) => {
        next()
      }).catch((response) => {
        next(false)
      })
    }
  }
</script>