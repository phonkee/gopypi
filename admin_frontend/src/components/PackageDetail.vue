<template>
    <page title="Package detail">
        <row>
            <column :lg="3">
                <panel title="Package info" icon="archive">
                    <dl>
                        <dt>Name</dt>
                        <dd>{{pack.name}}</dd>
                    </dl>
                    <dl>
                        <dt>Author</dt>
                        <dd>{{pack.author | fullName}}</dd>
                    </dl>
                    <dl>
                        <dt>Created</dt>
                        <dd>{{pack.created_at | date}}</dd>
                    </dl>
                    <dl>
                        <dt>Updated</dt>
                        <dd>{{pack.updated_at | date}}</dd>
                    </dl>
                </panel>
            </column>
            <column :lg="9" v-if="downloadStatsEnabled()">
                <row>
                    <column :lg="4" v-if="packageDownloadStats.all && packageDownloadStats.all.length">
                        <panel title="All downloads" icon="cloud-download">
                            <download-stats-chart :chart-data="packageDownloadStats.all"
                                                  style="height:300px"></download-stats-chart>
                        </panel>
                    </column>
                    <column :lg="4" v-if="packageDownloadStats.monthly && packageDownloadStats.monthly.length">
                        <panel title="Monthly downloads" icon="cloud-download">
                            <download-stats-chart :chart-data="packageDownloadStats.monthly"
                                                  style="height:300px"></download-stats-chart>
                        </panel>
                    </column>
                    <column :lg="4" v-if="packageDownloadStats.weekly && packageDownloadStats.weekly.length">
                        <panel title="Weekly downloads" icon="cloud-download">
                            <download-stats-chart :chart-data="packageDownloadStats.weekly"
                                                  style="height:300px"></download-stats-chart>
                        </panel>
                    </column>
                </row>
            </column>
        </row>
        <row>
            <column>
                <panel title="Maintainers">
                    <select2 :options="activeUsers" :selected="pack.maintainers || []" :label="fullName" @add="addMaintainer($event)" @remove="removeMaintainer($event)"></select2>
                </panel>
            </column>
        </row>
        <row v-if="pack.versions && pack.versions.length">
            <column>
                <panel title="Package versions" icon="code-fork">
                    <div class="panel-group" id="accordion" role="" aria-multiselectable="true">
                        <div class="panel panel-default" v-for="(version, index) in pack.versions">
                            <div class="panel-heading" role="tab" :id="'heading' + index">
                                <a role="button" data-toggle="collapse" data-parent="#accordion"
                                   :href="'#collapse' + index" aria-expanded="true" :aria-controls="'collapse' + index">
                                    {{pack.name}} v{{version.version}}
                                </a>

                                <small class="pull-right">
                                    Created: {{version.created_at | date}},
                                    Author: {{version.author | fullName}}
                                </small>
                            </div>
                            <div :id="'collapse' + index" class="panel-collapse collapse" role="tabpanel"
                                 :aria-labelledby="'heading' + index">
                                <div class="panel-body">
                                    <row>
                                        <column :lg="4">
                                            <dl>
                                                <dt>Created</dt>
                                                <dd>{{version.created_at | date}}</dd>
                                            </dl>
                                            <dl>
                                                <dt>Author</dt>
                                                <dd>{{version.author | fullName}}</dd>
                                            </dl>
                                            <dl v-if="version.home_page">
                                                <dt>Homepage</dt>
                                                <dd><a :href="version.home_page"
                                                       target="_blank">{{version.home_page}}</a></dd>
                                            </dl>
                                        </column>
                                        <column :lg="8" v-if="version.files && version.files.length > 0">
                                            <table width="100%" class="table table-striped table-bordered table-hover"
                                                   id="dataTables-example">
                                                <thead>
                                                    <tr>
                                                        <th>File</th>
                                                        <th>Created</th>
                                                        <th>Author</th>
                                                        <th>Md5 hash</th>
                                                    </tr>
                                                </thead>
                                                <tbody>
                                                    <tr class="odd" v-for="file in version.files">
                                                        <td><a :href="file.download_url" target="_blank">{{file.filename}}</a></td>
                                                        <td>{{file.created_at | date}}</td>
                                                        <td>{{file.author | fullName}}</td>
                                                        <td>{{file.md5_digest}}</td>
                                                    </tr>
                                                </tbody>
                                            </table>
                                        </column>
                                    </row>
                                </div>
                            </div>
                        </div>
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
  import store from '../store'
  import FormField from './forms/FormField.vue'
  import Well from './layout/Well.vue'
  import DownloadStatsChart from './charts/DownloadStatsChart.vue'
  import feature from '../api/feature'
  import Select2 from './layout/Select2.vue'
  import filters from '../filters'
  import pack from '../api/pack'

  export default {
    components: {
      Page,
      Panel,
      Well,
      Column,
      Row,
      FormField,
      Select2,
      DownloadStatsChart
    },
    computed: mapGetters({
      pack: 'singlePackage',
      packageDownloadStats: 'packageDownloadStats',
      info: 'allInfo',
      activeUsers: 'allUsers'
    }),
    beforeRouteEnter (route, redirect, next) {
      var pack = store.dispatch('getPackage', {
        id: route.params.id
      })
      var pdStats = store.dispatch('getPackageDownloadStats', {
        id: route.params.id
      })
      var activeUsers = store.dispatch('getActiveUsers')

      Promise.all([pack, pdStats, activeUsers]).then((pack) => {
        next()
      }).catch((response) => {
        next(false)
      })
    },
    methods: {
      downloadStatsEnabled () {
        return feature.hasFeature(this.info, feature.FEATURE_DOWNLOAD_STATS)
      },
      addMaintainer (data) {
        var vm = this
        pack.addMaintainer(vm.pack.id, data.id).then(() => {
          vm.$store.dispatch('messageSuccess', 'Maintainer added.')
          vm.$store.dispatch('getPackage', {
            id: vm.$route.params.id
          })
        })
      },
      removeMaintainer (data) {
        var vm = this
        pack.removeMaintainer(vm.pack.id, data.id).then(() => {
          vm.$store.dispatch('messageSuccess', 'Maintainer removed.')
          vm.$store.dispatch('getPackage', {
            id: vm.$route.params.id
          })
        })
      },
      fullName (data) {
        return filters.fullName(data)
      }
    }
  }
</script>
<style scoped>
    .charts {
        height: 100%;
    }
    .draggable {
        cursor: move;
    }
    .dropzone {
        min-height: 140px;
        border: 1px solid #eeeeee;
    }

</style>