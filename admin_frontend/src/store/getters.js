/**
 * Created by phonkee on 19.10.16.
 */

export const allInfo = state => state.info.all
export const allStats = state => state.stats.all
export const allDownloadStats = state => state.stats.downloadAll
export const packageDownloadStats = state => state.stats.pack
export const allUsers = state => state.auth.all
export const allUsersPaginator = state => state.auth.paginator
export const allPackages = state => state.pack.all
export const singlePackage = state => state.pack.pack
export const allLicenses = state => state.license.all
export const me = state => state.auth.me
export const myPackages = state => state.pack.myList
export const user = state => state.auth.user
export const messages = state => state.messages.all
export const spinnerPendingRequests = state => state.spinner.pending
