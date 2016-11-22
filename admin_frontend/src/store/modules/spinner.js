/**
 * Created by phonkee on 23.10.16.
 */
import * as types from '../mutation-types'

// initial state
const state = {
  pending: 0
}

// mutations
const mutations = {
  [types.SPINNER_ADD_PENDING] (state) {
    state.pending += 1
  },
  [types.SPINNER_REMOVE_PENDING] (state) {
    if (state.pending > 0) {
      state.pending -= 1
    }
  }
}

export default {
  state,
  mutations
}
