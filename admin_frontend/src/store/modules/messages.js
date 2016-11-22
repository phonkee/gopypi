/**
 * Created by phonkee on 22.10.16.
 */

import * as types from '../mutation-types'

// initial state
const state = {
  all: []
}

// mutations
const mutations = {
  [types.ADD_FLASH_MESSAGE] (state, payload) {
    state.all.push(payload)
  },

  [types.HIDE_FLASH_MESSAGE] (state, guid) {
    state.all.forEach(function (item, index, array) {
      if (item.guid === guid) {
        array.splice(item, 1)
      }
    })
  }
}

export default {
  state,
  mutations
}
