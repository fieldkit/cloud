
export function dateToString (d, short) {
  var month = d.getUTCMonth() + 1
  var day = d.getUTCDate()
  var hour = d.getUTCHours() + ''
  if (hour.length === 1) hour = '0' + hour
  var minute = d.getUTCMinutes() + ''
  if (minute.length === 1) minute = '0' + minute

  var monthString = ''
  if (month === 1) monthString = 'January'
  if (month === 2) monthString = 'February'
  if (month === 3) monthString = 'March'
  if (month === 4) monthString = 'April'
  if (month === 5) monthString = 'May'
  if (month === 6) monthString = 'June'
  if (month === 7) monthString = 'July'
  if (month === 8) monthString = 'August'
  if (month === 9) monthString = 'September'
  if (month === 10) monthString = 'October'
  if (month === 11) monthString = 'November'
  if (month === 12) monthString = 'December'
  if (short) monthString = monthString.slice(0, 3)

  return monthString + ' ' + day + ', ' + hour + ':' + minute
}

export function lerp (start, end, ratio) {
  return start + (end - start) * ratio
}

export function rgb2hex (rgb) {
  return ((rgb[0] * 255 << 16) + (rgb[1] * 255 << 8) + rgb[2] * 255)
}

export function getURLParameters () {
  let queryString = {}
  let query = window.location.search.substring(1)
  let vars = query.split('&')
  for (let i = 0; i < vars.length; i++) {
    let pair = vars[i].split('=')
        // If first entry with this name
    if (typeof queryString[pair[0]] === 'undefined') {
      queryString[pair[0]] = decodeURIComponent(pair[1])
        // If second entry with this name
    } else if (typeof queryString[pair[0]] === 'string') {
      let arr = [ queryString[pair[0]], decodeURIComponent(pair[1]) ]
      queryString[pair[0]] = arr
        // If third or later entry with this name
    } else {
      queryString[pair[0]].push(decodeURIComponent(pair[1]))
    }
  }
  return queryString
}
