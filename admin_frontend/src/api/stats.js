/**
 * Created by phonkee on 19.10.16.
 */
import Vue from 'vue'

export default {
  getServerStats () {
    return new Promise((resolve, reject) => {
      Vue.http.get('/api/stats/server/').then((response) => {
        resolve(response.data.result)
      }, (response) => {
        reject(response.data)
      })
    })
  },
  /*
  getDownloadStats fetches download stats for all packages
   */
  getDownloadStats () {
    return new Promise((resolve, reject) => {
      Vue.http.get('/api/stats/download/package/').then((response) => {
        resolve(response.data.result)
      }, (response) => {
        reject(response.data)
      })
    })
  },
  /*
  getPackageDownloadStats fetches download stats for given package id
   */
  getPackageDownloadStats (id) {
    return new Promise((resolve, reject) => {
      Vue.http.get('/api/stats/download/package/' + id + '/').then((response) => {
        resolve(response.data.result)
      }, (response) => {
        reject(response.data)
      })
    })
  },
  /*
  getPackageVersionDownloadStats fetches download stats for given package id and package version id
   */
  getPackageVersionDownloadStats (packageId, packageVersionId) {
    return new Promise((resolve, reject) => {
      Vue.http.get('/api/stats/download/package/' + packageId + '/version/' + packageVersionId + '/').then((response) => {
        resolve(response.data.result)
      }, (response) => {
        reject(response.data)
      })
    })
  }
}
