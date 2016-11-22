<template>
    <page title="User profile">
        <row>
            <column :lg="12">
                <panel icon="user" title="User profile">
                    <form role="form">
                        <div class="row">
                            <div class="col-lg-12">
                                <div class="form-group">
                                    <label>Username</label>
                                    <input class="form-control" disabled readonly v-model="me.username">
                                    <p class="help-block">Unique username.</p>
                                </div>
                            </div>
                        </div>

                        <div class="row">
                            <div class="col-lg-6">
                                <form-field v-bind:errors="errors" label="First name" field="first_name">
                                    <input class="form-control" v-model="me.first_name">
                                </form-field>
                            </div>
                            <div class="col-lg-6">
                                <form-field v-bind:errors="errors" label="Last name" field="last_name">
                                    <input class="form-control" v-model="me.last_name">
                                </form-field>
                            </div>
                        </div>

                        <div class="row">
                            <div class="col-lg-12">
                                <form-field v-bind:errors="errors" label="Email" field="email">
                                    <input class="form-control" v-model="me.email">
                                </form-field>
                            </div>
                        </div>

                        <button class="btn btn-default btn-primary" v-on:click.prevent="submit()">Submit</button>
                    </form>
                </panel>
            </column>
        </row>
    </page>
</template>
<script>
  import {mapGetters} from 'vuex'
  import auth from '../api/auth'

  import Column from './layout/Column.vue'
  import FormField from './forms/FormField.vue'
  import Page from './layout/Page.vue'
  import Panel from './layout/Panel.vue'
  import Row from './layout/Row.vue'

  export default {
    data () {
      return {
        errors: {}
      }
    },
    components: {
      Column, FormField, Page, Panel, Row
    },
    computed: mapGetters({
      me: 'me'
    }),
    created () {
      this.$store.dispatch('getMe')
    },
    methods: {
      submit () {
        var data = {
          'first_name': this.me.first_name,
          'last_name': this.me.last_name,
          'email': this.me.email
        }
        auth.updateMe(data).then((me) => {
          this.$store.dispatch('getMe')
          this.$store.dispatch('messageSuccess', 'Profile updated.')
          this.errors = {}
        }).catch((response) => {
          this.errors = response.data.error
        })
      }
    }
  }
</script>
