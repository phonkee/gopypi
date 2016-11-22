/**
 * Created by phonkee on 27.10.16.
 */
import Vue from 'vue'

export default {
  listLicenses () {
    return new Promise((resolve, reject) => {
      Vue.http.get('/api/license/').then((response) => {
        resolve(response.data.result)
      }, (response) => {
        reject(response)
      })
    })
  }
}
