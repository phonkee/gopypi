/**
 * Created by phonkee on 19.10.16.
 */
import Vue from 'vue'

export default {
  getInfo () {
    return new Promise((resolve, reject) => {
      Vue.http.get('/api/info/').then((response) => {
        resolve(response.data.result)
      }, (response) => {
        reject(response)
      })
    })
  }
}
