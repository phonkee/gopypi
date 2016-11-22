/**
 * Created by phonkee on 19.10.16.
 */
import Vue from 'vue'

export default {
  listPackages () {
    return new Promise((resolve, reject) => {
      Vue.http.get('/api/package/').then((response) => {
        resolve(response.data.result)
      }, (response) => {
        reject(response)
      })
    })
  },
  listMyPackages () {
    return new Promise((resolve, reject) => {
      Vue.http.get('/api/me/package/').then((response) => {
        resolve(response.data.result)
      }, (response) => {
        reject(response)
      })
    })
  },
  getPackage (id) {
    return new Promise((resolve, reject) => {
      Vue.http.get('/api/package/' + id + '/').then((response) => {
        resolve(response.data.result)
      }, (response) => {
        reject(response)
      })
    })
  },
  addMaintainer (packageId, userId) {
    return Vue.http.post('/api/package/' + packageId + '/maintainer/' + userId + '/')
  },
  removeMaintainer (packageId, userId) {
    return Vue.http.delete('/api/package/' + packageId + '/maintainer/' + userId + '/')
  }
}
