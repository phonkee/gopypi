<template>
    <page title="Add user">
        <row>
            <column :lg="12">
                <panel title="Add new user">
                    <form role="form">
                        <form-field v-bind:errors="errors" field="username" label="Username" required="true">
                            <input class="form-control" autofocus placeholder="Enter unique username"
                                   v-model.trim="newUser.username">
                        </form-field>

                        <form-field v-bind:errors="errors" field="first_name" label="First name">
                            <input class="form-control" v-model.trim="newUser.first_name">
                        </form-field>

                        <form-field v-bind:errors="errors" field="last_name" label="Last name">
                            <input class="form-control" v-model.trim="newUser.last_name">
                        </form-field>

                        <form-field v-bind:errors="errors" field="email" label="Email">
                            <input class="form-control" v-model.trim="newUser.email">
                        </form-field>
                        <div class="row">
                            <div class="col-md-6">
                                <form-field v-bind:errors="errors" field="password" label="Password" required="true">
                                    <input class="form-control" type="password" placeholder="Enter password"
                                           v-model.trim="newUser.password">
                                </form-field>
                            </div>
                            <div class="col-md-6">
                                <form-field v-bind:errors="errors" field="password2" label="Repeat password"
                                            required="true">
                                    <input class="form-control" type="password" placeholder="Retype password"
                                           v-model.trim="newUser.password2">
                                </form-field>
                            </div>
                        </div>

                        <div class="form-group">
                            <label>Flags</label>
                            <div class="checkbox">
                                <label>
                                    <input type="checkbox" v-model="newUser.is_active">User is active
                                    <small>(use this to set user as deleted)</small>
                                </label>
                            </div>
                            <div class="checkbox">
                                <label>
                                    <input type="checkbox" v-model="newUser.is_admin"> User is admin
                                    <small>(can access admin interface)</small>
                                </label>
                            </div>
                        </div>

                        <div class="form-group">
                            <label>Permissions</label>
                            <div class="checkbox">
                                <label>
                                    <input type="checkbox" v-model="newUser.can_list">Can list packages
                                </label>
                            </div>
                            <div class="checkbox">
                                <label>
                                    <input type="checkbox" v-model="newUser.can_download">Can download packages
                                </label>
                            </div>
                            <div class="checkbox">
                                <label>
                                    <input type="checkbox" v-model="newUser.can_create">Can create packages
                                </label>
                            </div>
                            <div class="checkbox">
                                <label>
                                    <input type="checkbox" v-model="newUser.can_update">Can update packages (if is
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
  import { mapGetters } from 'vuex'
  import Column from './layout/Column.vue'
  import FormField from './forms/FormField.vue'
  import Page from './layout/Page.vue'
  import Panel from './layout/Panel.vue'
  import Row from './layout/Row.vue'

  import auth from '../api/auth'
  export default {
    components: {
      Column, FormField, Page, Panel, Row
    },
    computed: mapGetters({
      info: 'allInfo'
    }),
    data () {
      return {
        errors: {},
        newUser: {
          username: '',
          first_name: '',
          last_name: '',
          email: '',
          password: '',
          password2: '',
          is_active: true,
          is_admin: false,
          can_list: true,
          can_download: true,
          can_create: false,
          can_update: false
        }
      }
    },
    methods: {
      submit () {
        var newUser = {
          username: this.newUser.username,
          first_name: this.newUser.first_name,
          last_name: this.newUser.last_name,
          email: this.newUser.email,
          password: this.newUser.password,
          password2: this.newUser.password2,
          is_active: this.newUser.is_active,
          is_admin: this.newUser.is_admin,
          can_list: this.newUser.can_list,
          can_create: this.newUser.can_create,
          can_download: this.newUser.can_download,
          can_update: this.newUser.can_update
        }

        auth.createUser(newUser, (data) => {
          this.errors = {}
        }, (response) => {
          this.errors = response.data.error
        })
      }
    }
  }
</script>