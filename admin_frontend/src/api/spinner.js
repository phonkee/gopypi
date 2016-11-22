/**
 * Created by phonkee on 23.10.16.
 */

import * as types from '../store/mutation-types'

export default {
  interceptor (store) {
    return (request, next) => {
      store.commit(types.SPINNER_ADD_PENDING)
      next((response) => {
        store.commit(types.SPINNER_REMOVE_PENDING)
      })
    }
  }
}
