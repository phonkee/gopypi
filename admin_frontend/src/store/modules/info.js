/**
 * Created by phonkee on 19.10.16.
 */

/**
 * Created by phonkee on 19.10.16.
 */

import * as types from '../mutation-types'

// initial state
const state = {
  all: {}
}

// mutations
const mutations = {
  [types.RECEIVE_INFO] (state, { info }) {
    state.all = info
  }
}

export default {
  state,
  mutations
}
