<template>
    <page title="Packages">
        <row>
            <column :lg="12">
                <panel title="Packages" icon="archive">
                    <table class="table table-striped table-bordered table-hover">
                        <thead>
                        <tr>
                            <th>#</th>
                            <th>Package</th>
                            <th>Author</th>
                            <th>Maintainers</th>
                            <th>Versions</th>
                            <th>Created at</th>
                            <th>Actions</th>
                        </tr>
                        </thead>
                        <tbody>
                        <tr v-for="pack in packages">
                            <td>{{ pack.id }}</td>
                            <td>{{ pack.name }}</td>
                            <td>{{ pack.author | fullName }}</td>
                            <td>
                                <div v-for="(maintainer, index) in pack.maintainers">
                                    {{maintainer | fullName}}
                                    <span v-if="index !== (pack.maintainers.length-1)">, </span>
                                </div>
                            </td>
                            <td>
                                <span v-for="(version, index) in pack.versions">{{version.version}}<span v-if="index !== (pack.versions.length-1)">, </span></span>
                            </td>
                            <td>{{ pack.created_at | date }}</td>
                            <td>
                                <router-link class="btn btn-default btn-xs"
                                             :to="{ name: 'admin.package.detail', params: { id: pack.id } }"
                                             exact>
                                    <i class="fa fa-search"></i></router-link>
                                </button>
                            </td>
                        </tr>
                        </tbody>
                    </table>
                </panel>
            </column>
        </row>
    </page>
</template>
<script>
  import {mapGetters} from 'vuex'
  import Column from './layout/Column.vue'
  import Page from './layout/Page.vue'
  import Panel from './layout/Panel.vue'
  import Row from './layout/Row.vue'
  import store from '../store'
  export default {
    components: {Column, Page, Panel, Row},
    computed: mapGetters({
      packages: 'allPackages'
    }),
    beforeRouteEnter (route, redirect, next) {
      store.dispatch('getAllPackages').then((packages) => {
        next()
      }).catch((response) => {
        next(false)
      })
    }
  }
</script>
