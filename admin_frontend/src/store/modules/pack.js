/**
 * Created by phonkee on 19.10.16.
 */

import * as types from '../mutation-types'

// initial state
const state = {
  all: [],
  pack: {},
  myList: []
}

// mutations
const mutations = {
  [types.RECEIVE_PACKAGE_LIST] (state, { packages }) {
    state.all = packages
  },
  [types.RECEIVE_SINGLE_PACKAGE] (state, { pack }) {
    state.pack = pack
  },
  [types.RECEIVE_MY_PACKAGE_LIST] (state, { packages }) {
    state.myList = packages
  }
}

export default {
  state,
  mutations
}
