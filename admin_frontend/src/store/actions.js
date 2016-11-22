/**
 * Created by phonkee on 19.10.16.
 */

import * as types from './mutation-types'

import auth from '../api/auth'
import info from '../api/info'
import license from '../api/license'
import pack from '../api/pack'
import stats from '../api/stats'
import * as utils from '../api/utils'

export const getAllServerStats = ({ commit }) => {
  return new Promise((resolve, reject) => {
    stats.getServerStats().then((stats) => {
      commit(types.RECEIVE_SERVER_STATS, { stats })
      resolve(stats)
    }).catch((err) => {
      reject(err)
    })
  })
}

export const getAllDownloadStats = ({ commit }) => {
  return new Promise((resolve, reject) => {
    stats.getDownloadStats().then((stats) => {
      commit(types.RECEIVE_ALL_DOWNLOAD_STATS, { stats })
      resolve(stats)
    }).catch((err) => {
      reject(err)
    })
  })
}

export const getPackageDownloadStats = ({ commit }, payload) => {
  return new Promise((resolve, reject) => {
    stats.getPackageDownloadStats(payload.id).then((stats) => {
      commit(types.RECEIVE_PACKAGE_DOWNLOAD_STATS, { stats })
      resolve(stats)
    }).catch((err) => {
      reject(err)
    })
  })
}

export const getAllInfo = ({ commit }) => {
  return info.getInfo().then(info => {
    commit(types.RECEIVE_INFO, { info })
  })
}

export const getAllLicenses = ({ commit, state }) => {
  license.listLicenses().then(licenses => {
    commit(types.RECEIVE_LICENSES, { licenses })
  })
}

export const getAllPackages = ({ commit, state }) => {
  pack.listPackages().then(packages => {
    commit(types.RECEIVE_PACKAGE_LIST, { packages })
  })
}

export const getMyPackages = ({ commit, state }) => {
  pack.listMyPackages().then(packages => {
    commit(types.RECEIVE_MY_PACKAGE_LIST, { packages })
  })
}

export const getPackage = ({ commit, dispatch }, payload) => {
  return new Promise((resolve, reject) => {
    pack.getPackage(payload.id).then((pack) => {
      resolve(pack)
      commit(types.RECEIVE_SINGLE_PACKAGE, { pack })
    }).catch((err) => {
      dispatch('messageError', 'Cannot get package: ' + err.statusText)
      reject(err)
    })
  })
}

export const getAllUsers = ({ commit, state }, payload) => {
  return auth.listUsers(payload).then(users => {
    commit(types.RECEIVE_USERS, { users })
  })
}

export const getActiveUsers = ({ commit, state }, payload) => {
  return auth.listUsers(payload, {is_active: true}).then(users => {
    commit(types.RECEIVE_USERS, { users })
  })
}
export const getMe = ({ commit }) => {
  auth.getMe().then((me) => {
    commit(types.RECEIVE_ME, { me })
  })
}

export const getUser = ({ commit, dispatch }, payload) => {
  return new Promise((resolve, reject) => {
    auth.getUser(payload.id).then((user) => {
      resolve(user)
      commit(types.RECEIVE_SINGLE_USER, { user })
    }).catch((err) => {
      dispatch('messageError', 'Cannot get user: ' + err.statusText)
      reject(err)
    })
  })
}

/*
Generic function to call message
 */
function message ({commit}, message, level) {
  var guid = utils.guid()
  commit(types.ADD_FLASH_MESSAGE, {
    message: message,
    level: level,
    guid: guid
  })

  setTimeout(() => {
    commit(types.HIDE_FLASH_MESSAGE, guid)
  }, 5000)
}

/*
Add error message
 */
export const messageError = ({ commit }, payload) => {
  message({commit}, payload, 'danger')
}

/*
Add info message
 */
export const messageInfo = ({ commit }, payload) => {
  message({commit}, payload, 'info')
}

/*
Add success message
 */
export const messageSuccess = ({ commit }, payload) => {
  message({commit}, payload, 'success')
}

/*
Add warning message
 */
export const messageWarning = ({ commit }, payload) => {
  message({commit}, payload, 'warning')
}

export const hideMessage = ({ commit }, payload) => {
  commit(types.HIDE_FLASH_MESSAGE, payload)
}
