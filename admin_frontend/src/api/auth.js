/**
 * Created by phonkee on 19.10.16.
 */

import Vue from 'vue'

export default {
  createUser (user, cb, errCb) {
    Vue.http.post('/api/user/', user).then((response) => {
      cb(response.data)
    }, (response) => {
      errCb(response)
    })
  },
  getMe () {
    return new Promise((resolve, reject) => {
      Vue.http.get('/api/me/').then((response) => {
        resolve(response.data.result)
      }, (response) => {
        reject(response)
      })
    })
  },
  /*
  GetUser returns promise which is resolved when user is found, otherwise it's rejected
   */
  getUser (id) {
    return new Promise((resolve, reject) => {
      Vue.http.get('/api/user/' + id + '/').then((response) => {
        resolve(response.data.result)
      }, (response) => {
        reject(response)
      })
    })
  },
  /*
  updateMe updates currently logged user
   */
  updateMe (data) {
    return new Promise((resolve, reject) => {
      Vue.http.post('/api/me/', data).then((response) => {
        resolve(response.data)
      }, (response) => {
        reject(response)
      })
    })
  },
  /*
  updatePassword updates password for currently logged user
   */
  updatePassword (data) {
    return new Promise((resolve, reject) => {
      Vue.http.post('/api/me/password/', data).then((response) => {
        resolve(response.data)
      }, (response) => {
        reject(response)
      })
    })
  },
  /*
  updateUser updates user on backend
   */
  updateUser (user, id) {
    return new Promise((resolve, reject) => {
      Vue.http.post('/api/user/' + id + '/', user).then((response) => {
        resolve(response.data.result)
      }, (response) => {
        reject(response)
      })
    })
  },
  listUsers (page, filter) {
    var params = Object.assign({}, filter || {})
    if (page) {
      params['page'] = page
    }
    return new Promise((resolve, reject) => {
      Vue.http.get('/api/user/', {params: params}).then((response) => {
        resolve(response.data)
      }, (response) => {
        reject(response)
      })
    })
  },
  /*
  login calls login with given username and password and returns promise.
  resolve resolves with token
  reject rejects with response
   */
  login (username, password) {
    return new Promise((resolve, reject) => {
      Vue.http.post('/api/login/', JSON.stringify({username: username, password: password})).then((response) => {
        window.localStorage.setItem('auth_token', response.headers.get('authorization'))
        resolve(response.headers.get('authorization'))
      }, (response) => {
        reject(response)
      })
    })
  },
  logout () {
    window.localStorage.removeItem('auth_token')
  },
  interceptor (router) {
    return (request, next) => {
      var authToken = window.localStorage.getItem('auth_token')
      if (authToken) {
        request.headers.set('Authorization', authToken)
      }
      next((response) => {
        // Check for expired token response, if expired, navigate to login
        if (response.status === 401) {
          router.push({name: 'login'})
        } else {
          return response
        }
      })
    }
  },
  getAthorizationHeader () {
    return 'Bearer ' + window.localStorage.getItem('auth_token')
  }
}
