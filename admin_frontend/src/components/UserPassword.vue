<template>
    <page title="Change password">
        <panel title="Change password" icon="key">
            <form role="form">
                <div class="row">
                    <div class="col-lg-12">
                        <form-field v-bind:errors="errors" label="Current password" field="current">
                            <input class="form-control" type="password" v-model="change.current">
                        </form-field>
                    </div>
                </div>
                <div class="row">
                    <div class="col-lg-6">
                        <form-field v-bind:errors="errors" label="New password" field="password">
                            <input class="form-control" type="password" v-model="change.password">
                        </form-field>
                    </div>
                    <div class="col-lg-6">
                        <form-field v-bind:errors="errors" label="Retype password" field="password2">
                            <input class="form-control" type="password" v-model="change.password2">
                        </form-field>
                    </div>
                </div>
                <button class="btn btn-default btn-primary" v-on:click.prevent="submit()">Submit</button>
            </form>
        </panel>
    </page>
</template>
<script>
  import auth from '../api/auth'
  import Page from './layout/Page.vue'
  import Panel from './layout/Panel.vue'
  import Row from './layout/Row.vue'
  import Column from './layout/Column.vue'
  import FormField from './forms/FormField.vue'
  export default {
    components: {
      Page, Panel, Row, Column, FormField
    },
    data () {
      return {
        change: {
          current: '',
          password: '',
          password2: ''
        },
        errors: {}
      }
    },
    methods: {
      submit () {
        var data = this.change
        auth.updatePassword(data).then((data) => {
          this.$store.dispatch('messageSuccess', 'Password changed.')
          this.errors = {}
          this.change = {}
        }).catch((response) => {
          this.errors = response.data.error
        })
      }
    }
  }
</script>