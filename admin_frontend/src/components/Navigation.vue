<template>

    <!-- Navigation -->
    <nav class="navbar navbar-default navbar-static-top" role="navigation" style="margin-bottom: 0">
        <div class="navbar-header">
            <button type="button" class="navbar-toggle" data-toggle="collapse" data-target=".navbar-collapse">
                <span class="sr-only">Toggle navigation</span>
                <span class="icon-bar"></span>
                <span class="icon-bar"></span>
                <span class="icon-bar"></span>
            </button>
            <a class="navbar-brand" href="https://github.com/phonkee/gopypi" target="_blank">gopypi <small>{{ info.version }}</small></a>
        </div>
        <!-- /.navbar-header -->

        <ul class="nav navbar-top-links navbar-right">

            <li>
                <spinner></spinner>
            </li>

            <!-- /.dropdown -->
            <li class="dropdown">
                <a class="dropdown-toggle" data-toggle="dropdown" href="#">
                    <i class="fa fa-user fa-fw"></i> <i class="fa fa-caret-down"></i>
                </a>
                <ul class="dropdown-menu dropdown-user">
                    <li>
                        <router-link :to="{name: 'admin.user.profile'}" active-class="active" exact><i class="fa fa-user fa-fw"></i> User Profile</router-link>
                    </li>
                    <li>
                        <router-link :to="{name: 'admin.user.password'}" active-class="active" exact><i class="fa fa-key fa-fw"></i> Password</router-link>
                    </li>
                    <li class="divider"></li>
                    <li>
                        <a v-on:click.prevent="logout()" href="#"><i class="fa fa-sign-out fa-fw"></i> Logout</a>
                    </li>
                </ul>
                <!-- /.dropdown-user -->
            </li>
            <!-- /.dropdown -->
        </ul>
        <!-- /.navbar-top-links -->

        <div class="navbar-default sidebar" role="navigation">
            <div class="sidebar-nav navbar-collapse">
                <ul class="nav side-menu">
                    <!--
                    <li class="sidebar-search">
                        <div class="input-group custom-search-form">
                            <input type="text" class="form-control" placeholder="Search...">
                            <span class="input-group-btn">
                            <button class="btn btn-default" type="button">
                                <i class="fa fa-search"></i>
                            </button>
                        </span>
                        </div>
                    </li>
                    -->
                    <li>
                        <router-link :to="{name: 'admin.dashboard'}" active-class="active"><i class="fa fa-dashboard fa-fw" exact></i> Dashboard</router-link>
                    </li>
                    <li>
                        <a href="#"><i class="fa fa-user fa-fw"></i> Users<span class="fa arrow"></span></a>
                        <ul class="nav nav-second-level">
                            <li class="active">
                                <router-link :to="{name: 'admin.user.list'}" active-class="active" exact>List users</router-link>
                            </li>
                            <li>
                                <router-link :to="{name: 'admin.user.add'}" active-class="active" exact>Add user</router-link>
                            </li>
                        </ul>
                    </li>
                    <li>
                        <router-link :to="{name: 'admin.package.list'}" active-class="active" exact><i class="fa fa-archive fa-fw"></i> Packages</router-link>
                    </li>
                    <li v-if="downloadStatsEnabled()">
                        <router-link :to="{name: 'admin.stats.download.list'}" active-class="active" exact><i class="fa fa-cloud-download fa-fw"></i> Downloads</router-link>
                    </li>
                    <li>
                        <router-link :to="{name: 'admin.license.list'}" active-class="active" exact><i class="fa fa-certificate fa-fw"></i> Licenses</router-link>
                    </li>
                    <li>
                        <router-link :to="{name: 'admin.feature.list'}" active-class="active" exact><i class="fa fa-cogs fa-fw"></i> Features</router-link>
                    </li>
                    <li>
                        <router-link :to="{name: 'admin.sponsors' }" active-class="active" exact><i class="fa fa-support fa-fw"></i> Sponsors</router-link>
                    </li>
                    <li>
                        <router-link :to="{name: 'admin.howto' }" active-class="active" exact><i class="fa fa-question-circle fa-fw"></i> How to</router-link>
                    </li>
                </ul>
            </div>
            <!-- /.sidebar-collapse -->
        </div>
        <!-- /.navbar-static-side -->
    </nav>
</template>

<script>
  import $ from 'jquery'
  import Spinner from './layout/Spinner.vue'
  import { mapGetters } from 'vuex'
  import auth from '../api/auth'
  import feature from '../api/feature'
  require('metismenu')

  export default {
    components: {
      Spinner
    },
    computed: mapGetters({
      info: 'allInfo'
    }),
    mounted: function () {
      $(() => {
        $(this.$el).find('ul.side-menu').metisMenu()
      })
    },
    methods: {
      logout () {
        auth.logout()
        this.$router.push({name: 'login'})
      },
      downloadStatsEnabled () {
        return feature.hasFeature(this.info, feature.FEATURE_DOWNLOAD_STATS)
      }
    }
  }
</script>