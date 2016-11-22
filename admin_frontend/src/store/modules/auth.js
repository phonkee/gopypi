/**
 * Created by phonkee on 19.10.16.
 */

import * as types from '../mutation-types'

// initial state
const state = {
  all: [],
  paginator: {},
  me: {},
  user: {}
}

// mutations
const mutations = {
  [types.RECEIVE_USERS] (state, { users }) {
    state.all = users.result
    state.paginator = users.paginator
  },
  [types.RECEIVE_ME] (state, { me }) {
    state.me = me
  },
  [types.RECEIVE_SINGLE_USER] (state, { user }) {
    state.user = user
  }
}

export default {
  state,
  mutations
}
