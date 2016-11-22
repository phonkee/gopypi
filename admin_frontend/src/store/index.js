/**
 * Created by phonkee on 19.10.16.
 */

import Vue from 'vue'
import Vuex from 'vuex'

import * as actions from './actions'
import * as getters from './getters'
import auth from './modules/auth'
import info from './modules/info'
import license from './modules/license'
import messages from './modules/messages'
import pack from './modules/pack'
import spinner from './modules/spinner'
import stats from './modules/stats'

Vue.use(Vuex)

export default new Vuex.Store({
  actions,
  getters,
  modules: {
    auth,
    info,
    license,
    messages,
    pack,
    spinner,
    stats
  }
})
