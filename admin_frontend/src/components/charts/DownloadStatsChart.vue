<template>
    <div class="flot-chart">
        <div class="flot-chart-content"></div>
    </div>
</template>
<script>
  import $ from 'jquery'
  var moment = require('moment')
  export default {
    props: {
      chartData: Array,
      aggregation: {
        default: 'all'
      }
    },
    mounted () {
      var rows = []
      for (let item of this.chartData) {
        var parsed = moment(item.created_at)
        var key = parsed.format('YYYY')
        switch (this.aggregation) {
          case 'monthly':
            key = parsed.format('MMM/YYYY')
            break
          case 'weekly':
            key = parsed.format('W/YYYY')
            break
        }
        rows.push([key, item.downloads])
      }

      // prepare label
      var label = 'Year'
      switch (this.aggregation) {
        case 'monthly':
          label = 'Month'
          break
        case 'weekly':
          label = 'Week'
          break
      }

      var vm = this
      var barOptions = {
        series: {
          bars: {
            show: true,
            barWidth: 0.5,
            align: 'center'
          }
        },
        grid: {
          hoverable: true
        },
        tooltip: true,
        tooltipOpts: {
          content: label + ': %x <br>Downloads: %y'
        },
        xaxis: {
          mode: 'categories',
          tickSize: 0,
          tickDecimals: 0
        }
      }
      var barData = {
        label: label,
        data: rows
      }
      $.plot($(vm.$el).find('.flot-chart-content'), [barData], barOptions)
    }
  }
</script>
