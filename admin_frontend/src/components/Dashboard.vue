<template>
    <page title="Dashboard">
        <row>
            <column :md="6" :lg="3">
                <div class="panel panel-primary">
                    <div class="panel-heading">
                        <div class="row">
                            <div class="col-xs-3">
                                <i class="fa fa-archive fa-5x"></i>
                            </div>
                            <div class="col-xs-9 text-right">
                                <div class="huge">
                                    {{ stats.packages }}
                                </div>
                                <div>Packages</div>
                            </div>
                        </div>
                    </div>
                    <router-link :to="{ name: 'admin.package.list' }">
                        <div class="panel-footer">
                            <span class="pull-left">View Details</span>
                            <span class="pull-right"><i class="fa fa-arrow-circle-right"></i></span>
                            <div class="clearfix"></div>
                        </div>
                    </router-link>
                </div>
            </column>
            <column :md="6" :lg="3" v-if="downloadStatsEnabled()">
                <div class="panel panel-green">
                    <div class="panel-heading">
                        <div class="row">
                            <div class="col-xs-3">
                                <i class="fa fa-cloud-download fa-5x"></i>
                            </div>
                            <div class="col-xs-9 text-right">
                                <div class="huge">{{ stats.downloads }}</div>
                                <div>Downloads</div>
                            </div>
                        </div>
                    </div>
                    <router-link :to="{ name: 'admin.stats.download.list' }">
                        <div class="panel-footer">
                            <span class="pull-left">View Details</span>
                            <span class="pull-right"><i class="fa fa-arrow-circle-right"></i></span>
                            <div class="clearfix"></div>
                        </div>
                    </router-link>
                </div>
            </column>
            <column :md="6" :lg="3">
                <div class="panel panel-red">
                    <div class="panel-heading">
                        <div class="row">
                            <div class="col-xs-3">
                                <i class="fa fa-user fa-5x"></i>
                            </div>
                            <div class="col-xs-9 text-right">
                                <div class="huge">{{ stats.active_users }}</div>
                                <div>Active users</div>
                            </div>
                        </div>
                    </div>
                    <router-link :to="{ name: 'admin.user.list' }">
                        <div class="panel-footer">
                            <span class="pull-left">View Details</span>
                            <span class="pull-right"><i class="fa fa-arrow-circle-right"></i></span>
                            <div class="clearfix"></div>
                        </div>
                    </router-link>
                </div>
            </column>
            <column :md="6" :lg="3">
                <div class="panel panel-yellow">
                    <div class="panel-heading">
                        <div class="row">
                            <div class="col-xs-3">
                                <i class="fa fa-certificate fa-5x"></i>
                            </div>
                            <div class="col-xs-9 text-right">
                                <div class="huge">{{ stats.licenses }}</div>
                                <div>Licenses</div>
                            </div>
                        </div>
                    </div>
                    <router-link :to="{ name: 'admin.license.list' }">
                        <div class="panel-footer">
                            <span class="pull-left">View Details</span>
                            <span class="pull-right"><i class="fa fa-arrow-circle-right"></i></span>
                            <div class="clearfix"></div>
                        </div>
                    </router-link>
                </div>
            </column>
        </row>

        <row v-if="myPackages.length > 0">
            <column :lg="12">
                <panel title="My packages" icon="archive">
                    <div class="table-responsive">
                        <table class="table table-striped">
                            <thead>
                            <tr>
                                <th>Name</th>
                                <th>Author</th>
                                <th>Maintainers</th>
                                <th>Versions</th>
                                <th>Created</th>
                                <th></th>
                            </tr>
                            </thead>
                            <tbody>
                            <tr v-for="pack in myPackages">
                                <td>{{pack.name}}</td>
                                <td>{{pack.author | fullName}}</td>
                                <td>
                                    <div v-for="(maintainer, index) in pack.maintainers">
                                        {{maintainer | fullName}}
                                        <span v-if="index !== (pack.maintainers.length-1)">, </span>
                                    </div>
                                </td>
                                <td>
                                    <span v-for="(version, index) in pack.versions">{{version.version}}<span
                                            v-if="index !== (pack.versions.length-1)">, </span></span>
                                </td>
                                <td>{{pack.created_at | date}}</td>
                                <td>
                                <router-link class="btn btn-default btn-xs"
                                             :to="{ name: 'admin.package.detail', params: { id: pack.id } }"
                                             exact>
                                    <i class="fa fa-search"></i></router-link>
                                </td>
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
  import Navigation from './Navigation.vue'

  import Column from './layout/Column.vue'
  import Row from './layout/Row.vue'
  import Page from './layout/Page.vue'
  import Panel from './layout/Panel.vue'
  import { mapGetters } from 'vuex'
  import store from '../store'
  import feature from '../api/feature'

  export default {
    components: {
      Column, Row, Navigation, Page, Panel
    },
    computed: mapGetters({
      stats: 'allStats',
      myPackages: 'myPackages',
      info: 'allInfo'
    }),
    data: function () {
      return {
        packageCount: 0
      }
    },
    beforeRouteEnter (route, redirect, next) {
      var stats = store.dispatch('getAllServerStats')
      var myPackages = store.dispatch('getMyPackages')

      Promise.all([stats, myPackages]).then(values => {
        next()
      }, reason => {
        next(false)
      })
    },
    methods: {
      downloadStatsEnabled () {
        return feature.hasFeature(this.info, feature.FEATURE_DOWNLOAD_STATS)
      }
    }
  }
</script>
