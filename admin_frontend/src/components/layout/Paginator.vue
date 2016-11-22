<template>
    <div class="col-lg12 text-center">
        <nav aria-label="Page navigation" v-if="source.num_pages > 1">
            <ul class="pagination">
                <li v-if="source.page > 1">
                    <a href="#" v-on:click.prevent="onChange(source.page-1)" aria-label="Previous">
                        <span aria-hidden="true">&laquo;</span>
                    </a>
                </li>
                <li v-for="i in range(1, source.num_pages)" v-bind:class="[i==source.page ? 'active' : '']"><a href="#" v-on:click.prevent="onChange(i)">{{i}}</a></li>
                <li v-if="source.page < source.num_pages">
                    <a href="#" v-on:click.prevent="onChange(source)" aria-label="Next">
                        <span aria-hidden="true">&raquo;</span>
                    </a>
                </li>
            </ul>
        </nav>
    </div>
</template>
<script>
  export default {
    props: {
      source: {
        type: Object,
        default: 1
      },
      urlParam: {
        type: String
      }
    },
    methods: {
      /*
      This function reads information from url when url-param is provided as prop
       */
      getValueFromURL: function () {
        if (!this.urlParam) {
          return
        }
        var page = parseInt(this.$route.query[this.urlParam], 10)
        if (isNaN(page) || page < 1) {
          page = 1
        }
        this.onChange(page)
      },
      onChange: function (value) {
        if (this.urlParam !== '') {
          var query = {}
          query[this.urlParam] = value
          this.$router.replace({ query })
        }
        this.source['page'] = value
        this.$emit('change', this.source)
      },
      range (start, end) {
        var result = []
        for (var i = start; i <= end; i++) {
          result.push(i)
        }
        return result
      }
    },
    mounted () {
      this.getValueFromURL()
    }
  }
</script>
