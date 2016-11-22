/**
 * Created by phonkee on 22.10.16.
 */

/*
guid returns identifier that resembles to uuid. it has the same format and uses random.
 */
export const guid = () => {
  function s4 () {
    return Math.floor((1 + Math.random()) * 0x10000)
      .toString(16)
      .substring(1)
  }
  return s4() + s4() + '-' + s4() + '-' + s4() + '-' +
    s4() + '-' + s4() + s4() + s4()
}
