/*
Set of multiple filters
 */
export default {
  fullName (user) {
    if (!user.id) {
      return ''
    }
    if (user.first_name && user.last_name) {
      return user.first_name + ' ' + user.last_name
    }
    return user.username
  }
}
