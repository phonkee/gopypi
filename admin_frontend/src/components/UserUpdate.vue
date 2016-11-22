<template>
    <page title="Update user">
        <row>
            <column :lg="12">
                <panel title="Update user">
                    <form role="form">
                        <form-field v-bind:errors="errors" field="username" label="Username" required="true">
                            <input class="form-control" autofocus placeholder="Enter unique username"
                                   v-model.trim="user.username">
                        </form-field>

                        <form-field v-bind:errors="errors" field="first_name" label="First name">
                            <input class="form-control" v-model.trim="user.first_name">
                        </form-field>

                        <form-field v-bind:errors="errors" field="last_name" label="Last name">
                            <input class="form-control" v-model.trim="user.last_name">
                        </form-field>

                        <form-field v-bind:errors="errors" field="email" label="Email">
                            <input class="form-control" v-model.trim="user.email">
                        </form-field>

                        <div class="row">
                            <div class="col-md-6">
                                <form-field v-bind:errors="errors" field="password" label="Password">
                                    <input class="form-control" type="password" placeholder="Enter password" v-model.trim="user.password">
                                </form-field>
                            </div>
                            <div class="col-md-6">
                                <form-field v-bind:errors="errors" field="password2" label="Repeat password">
                                    <input class="form-control" type="password" placeholder="Retype password" v-model.trim="user.password2">
                                </form-field>
                            </div>
                        </div>

                        <div class="form-group">
                            <label>Flags</label>
                            <div class="checkbox">
                                <label>
                                    <input type="checkbox" v-model="user.is_active">User is active
                                    <small>(use this to set user as deleted)</small>
                                </label>
                            </div>
                            <div class="checkbox">
                                <label>
                                    <input type="checkbox" v-model="user.is_admin"> User is admin
                                    <small>(can access admin interface)</small>
                                </label>
                            </div>
                        </div>

                        <div class="form-group">
                            <label>Permissions</label>
                            <div class="checkbox">
                                <label>
                                    <input type="checkbox" v-model="user.can_list">Can list packages
                                </label>
                            </div>
                            <div class="checkbox">
                                <label>
                                    <input type="checkbox" v-model="user.can_download">Can download packages
                                </label>
                            </div>
                            <div class="checkbox">
                                <label>
                                    <input type="checkbox" v-model="user.can_create">Can create packages
                                </label>
                            </div>
                            <div class="checkbox">
                                <label>
                                    <input type="checkbox" v-model="user.can_update">Can update packages (if is
                                    maintainer of package)
                                </label>
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
  import Column from './layout/Column.vue'
  import FormField from './forms/FormField.vue'
  import Page from './layout/Page.vue'
  import Panel from './layout/Panel.vue'
  import Row from './layout/Row.vue'
  import {mapGetters} from 'vuex'

  import auth from '../api/auth'
  import store from '../store'

  export default {
    components: {Column, Row, FormField, Page, Panel},
    computed: mapGetters({
      user: 'user'
    }),
    beforeRouteEnter (route, redirect, next) {
      store.dispatch('getUser', {
        id: route.params.id
      }).then((user) => {
        next()
      }).catch((response) => {
        next(false)
      })
    },
    data () {
      return {
        errors: {}
      }
    },
    methods: {
      submit () {
        auth.updateUser(this.user, this.$route.params.id).then((user) => {
          this.$store.dispatch('messageSuccess', 'User has been updated!')
          this.$router.push({name: 'admin.user.list'})
        }).catch((response) => {
          this.$store.dispatch('messageError', 'User update returned error!')
          this.errors = response.data.error
        })
      }
    }
  }

</script>
