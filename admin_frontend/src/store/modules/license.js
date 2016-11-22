/**
 * Created by phonkee on 27.10.16.
 */

import * as types from '../mutation-types'

// initial state
const state = {
  all: []
}

// mutations
const mutations = {
  [types.RECEIVE_LICENSES] (state, { licenses }) {
    state.all = licenses
  }
}

export default {
  state,
  mutations
}
