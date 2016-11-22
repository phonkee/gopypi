<template>
    <page title="Sponsors">
        <row>
            <column :lg="12">
                <panel title="Gopypi development is sponsored by:">
                        Your company can be here.
                </panel>
            </column>
        </row>
    </page>
</template>
<script>
  import Column from './layout/Column.vue'
  import Page from './layout/Page.vue'
  import Panel from './layout/Panel.vue'
  import Row from './layout/Row.vue'

  export default {
    components: {Column, Page, Panel, Row}
  }
</script>