/**
 * Created by phonkee on 19.10.16.
 */

import Vue from 'vue'

export default {
  FEATURE_DOWNLOAD_STATS: 'download_stats',
  updateFeature (id, value) {
    return Vue.http.post('/api/feature/' + id + '/', {value: value})
  },
  hasFeature (info, id) {
    for (var f of info.features) {
      if (f.id === id) {
        return f.value
      }
    }
    return false
  }
}
