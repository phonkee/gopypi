/**
 * Created by phonkee on 19.10.16.
 */

import store from '../store'
export default {
  interceptor (router) {
    return (request, next) => {
      next((response) => {
        // Check for expired token response, if expired, navigate to login
        if (response.status === 500) {
          store.dispatch('messageError', response.statusText + ' [' + request.url + ']')
          router.push('login')
        } else {
          return response
        }
      })
    }
  }
}
