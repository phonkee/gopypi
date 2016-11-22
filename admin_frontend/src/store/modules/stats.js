/**
 * Created by phonkee on 19.10.16.
 */

import * as types from '../mutation-types'

// initial state
const state = {
  all: {},
  downloadAll: {},
  pack: {}
}

// mutations
const mutations = {
  [types.RECEIVE_SERVER_STATS] (state, { stats }) {
    state.all = stats
  },
  [types.RECEIVE_ALL_DOWNLOAD_STATS] (state, {stats}) {
    state.downloadAll = stats
  },
  [types.RECEIVE_PACKAGE_DOWNLOAD_STATS] (state, {stats}) {
    state.pack = stats
  }
}

export default {
  state,
  mutations
}
