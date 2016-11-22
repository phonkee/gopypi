<template>
    <page title="Howto">
        <row>
            <column :lg="12">
                <well lg="true">
                    <h4>
                        Please follow instructions how to write proper setup.py with gopypi
                        <a href="http://blog.xelnor.net/private-pypi/" target="_blank">here</a>
                    </h4>
                </well>
            </column>
        </row>
    </page>
</template>
<script>
  import Column from './layout/Column.vue'
  import Page from './layout/Page.vue'
  import Row from './layout/Row.vue'
  import Well from './layout/Well.vue'
  export default {
    components: {Column, Page, Row, Well}
  }
</script>