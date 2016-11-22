import Vue from 'vue'

import VueCharts from 'vue-charts'
import VueFilter from 'vue-filter'
import VueRouter from 'vue-router'
import VueResource from 'vue-resource'

Vue.use(VueCharts)
Vue.use(VueFilter)
Vue.use(VueRouter)
Vue.use(VueResource)
Vue.use(require('vue-moment'))

import $ from 'jquery'
window.$ = $
window.jQuery = $

require('imports-loader')
require('bootstrap/dist/js/bootstrap.min.js')
require('bootstrap/dist/css/bootstrap.css')
require('../template/js/sb-admin-2.js')
require('../template/dist/css/sb-admin-2.min.css')
require('../static/gopypi.css')
require('../template/vendor/font-awesome/css/font-awesome.min.css')
require('jquery/dist/jquery.js')
require('../template/vendor/flot/excanvas.min.js')
require('../template/vendor/flot/jquery.flot.js')
require('../template/vendor/flot-tooltip/jquery.flot.tooltip.min.js')
require('imports?this=>window!../template/vendor/flot/jquery.flot.categories.js')
require('imports?this=>window!../template/vendor/flot/jquery.flot.resize.js')
require('../template/vendor/flot/jquery.flot.time.js')
require('select2/dist/js/select2.full.min.js')
require('select2/dist/css/select2.min.css')

import Dashboard from './components/Dashboard.vue'
import DownloadStatsAll from './components/DownloadStatsAll.vue'
import Admin from './components/Admin.vue'
import FeatureList from './components/FeatureList.vue'
import Howto from './components/Howto.vue'
import Login from './components/Login.vue'
import LicenseList from './components/LicenseList.vue'
import PackageList from './components/PackageList.vue'
import PackageDetail from './components/PackageDetail.vue'
import Sponsors from './components/Sponsors.vue'
import UserAdd from './components/UserAdd.vue'
import UserList from './components/UserList.vue'
import UserPassword from './components/UserPassword.vue'
import UserProfile from './components/UserProfile.vue'
import UserUpdate from './components/UserUpdate.vue'
import filters from './filters'

import store from './store'
import auth from './api/auth'
import spinner from './api/spinner'
import system from './api/system'

Vue.filter('fullName', filters.fullName)

// Prepare router
export const router = new VueRouter({
  mode: 'history',
  routes: [{
    path: '/admin/login',
    component: Login,
    name: 'login'
  }, {
    path: '/admin',
    component: Admin,
    children: [{
      component: UserAdd,
      path: 'user/add/',
      name: 'admin.user.add'
    }, {
      component: UserList,
      path: 'user/',
      name: 'admin.user.list'
    }, {
      component: UserUpdate,
      path: 'user/:id/',
      name: 'admin.user.update'
    }, {
      component: UserProfile,
      path: 'user/profile/',
      name: 'admin.user.profile'
    }, {
      component: UserPassword,
      path: 'user/password/',
      name: 'admin.user.password'
    }, {
      component: PackageList,
      path: 'package/',
      name: 'admin.package.list'
    }, {
      component: PackageDetail,
      path: 'package/:id/',
      name: 'admin.package.detail'
    }, {
      component: LicenseList,
      path: 'license/',
      name: 'admin.license.list'
    }, {
      path: 'dashboard/',
      component: Dashboard,
      name: 'admin.dashboard'
    }, {
      path: 'sponsors/',
      component: Sponsors,
      name: 'admin.sponsors'
    }, {
      path: 'feature/',
      component: FeatureList,
      name: 'admin.feature.list'
    }, {
      path: 'how-to/',
      component: Howto,
      name: 'admin.howto'
    }, {
      component: DownloadStatsAll,
      path: 'stats/download/',
      name: 'admin.stats.download.list'
    }]
  }, {
    path: '*',
    redirect: '/admin/dashboard'
  }]
})

Vue.http.interceptors.push(auth.interceptor(router))
Vue.http.interceptors.push(spinner.interceptor(store))
Vue.http.interceptors.push(system.interceptor(router))

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  store,
  render: h => h(require('./App.vue'))
})
