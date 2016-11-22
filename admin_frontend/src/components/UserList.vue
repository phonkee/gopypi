<template>
    <page title="Users">
        <row>
            <column :lg="12">
                <panel icon="user" title="Users">
                    <div class="table-responsive">
                        <table class="table table-striped table-bordered table-hover">
                            <thead>
                            <tr>
                                <th>#</th>
                                <th>Username</th>
                                <th>First name</th>
                                <th>Last name</th>
                                <th>Email</th>
                                <th>Active</th>
                                <th>Admin</th>
                                <th>L/D/C/U</th>
                                <th>Actions</th>
                            </tr>
                            </thead>
                            <tbody>
                            <tr v-for="user in users">
                                <td>{{ user.id }}</td>
                                <td>{{ user.username }}</td>
                                <td>{{ user.first_name }}</td>
                                <td>{{ user.last_name }}</td>
                                <td><a v-bind:href="'mailto:' + user.email" v-if="user.email">{{ user.email }}</a></td>
                                <td>
                                    <true-false-icon :value="user.is_active"></true-false-icon>
                                </i>
                                </td>
                                <td>
                                    <true-false-icon :value="user.is_admin"></true-false-icon>
                                </td>
                                <td>
                                    <true-false-icon :value="user.can_list"></true-false-icon> /
                                    <true-false-icon :value="user.can_download"></true-false-icon> /
                                    <true-false-icon :value="user.can_create"></true-false-icon> /
                                    <true-false-icon :value="user.can_update"></true-false-icon>
                                </td>
                                <td>
                                    <router-link class="btn btn-default btn-xs" :to="{ name: 'admin.user.update', params: { id: user.id } }" :title="'Update user ' + user.username" exact>
                                        <i class="fa fa-edit"></i></router-link>
                                    </button>
                                </td>
                            </tr>
                            </tbody>
                        </table>
                        <paginator v-bind:source="paginator" v-on:change="changedPaging" url-param="page"></paginator>
                    </div>
                </panel>
            </column>
        </row>
    </page>
</template>
<script>
  import {mapGetters} from 'vuex'
  import Column from './layout/Column.vue'
  import Page from './layout/Page.vue'
  import Panel from './layout/Panel.vue'
  import Row from './layout/Row.vue'

  import Paginator from './layout/Paginator.vue'
  import TrueFalseIcon from './layout/TrueFalseIcon.vue'

  import store from '../store'
  export default {
    components: {Column, Page, Paginator, Panel, Row, TrueFalseIcon},
    computed: mapGetters({
      users: 'allUsers',
      paginator: 'allUsersPaginator'
    }),
    beforeRouteEnter (route, redirect, next) {
      store.dispatch('getAllUsers').then((users) => {
        next()
      }).catch((response) => {
        store.dispatch('messageError', 'Cannot fetch users: ' + response.status)
        next(false)
      })
    },
    methods: {
      changedPaging (paginator) {
        store.dispatch('getAllUsers', paginator.page)
      }
    }
  }
</script>