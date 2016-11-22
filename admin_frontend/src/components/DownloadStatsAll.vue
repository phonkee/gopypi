<template>
    <page title="Downloads">
        <row>
            <column :lg="4" v-if="allDownloadStats.all && allDownloadStats.all.length">
                <panel title="All downloads" icon="cloud-download">
                    <download-stats-chart :chart-data="allDownloadStats.all"></download-stats-chart>
                </panel>
            </column>
            <column :lg="4" v-if="allDownloadStats.monthly && allDownloadStats.monthly.length">
                <panel title="Monthly downloads" icon="cloud-download">
                    <download-stats-chart :chart-data="allDownloadStats.monthly" aggregation="monthly"></download-stats-chart>
                </panel>
            </column>
            <column :lg="4" v-if="allDownloadStats.weekly && allDownloadStats.weekly.length">
                <panel title="Weekly downloads" icon="cloud-download">
                    <download-stats-chart :chart-data="allDownloadStats.weekly" aggregation="weekly"></download-stats-chart>
                </panel>
            </column>
        </row>
    </page>
</template>
<script>
  import {mapGetters} from 'vuex'
  import store from '../store'
  import Page from './layout/Page.vue'
  import Row from './layout/Row.vue'
  import Column from './layout/Column.vue'
  import Panel from './layout/Panel.vue'
  import DownloadStatsChart from './charts/DownloadStatsChart.vue'
  export default {
    components: {
      Page, Row, Column, Panel, DownloadStatsChart
    },
    computed: mapGetters({
      allDownloadStats: 'allDownloadStats'
    }),
    beforeRouteEnter (route, redirect, next) {
      store.dispatch('getAllDownloadStats').then((stats) => {
        next()
      }).catch((response) => {
        next(false)
      })
    }
  }
</script>
