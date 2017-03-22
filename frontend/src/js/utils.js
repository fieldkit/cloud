
export class BaseError {
  name: string;
  message: string;
  stack: ?string;

  constructor(message: string = 'Error') {
    this.name = this.constructor.name;
    this.message = message;
    if (typeof Error.captureStackTrace === 'function') {
      Error.captureStackTrace(this, this.constructor);
    } else {
      this.stack = (new Error(message)).stack;
    }
  }
}

export function map (val, domain1, domain2, range1, range2) {
  if (domain1 === domain2) return range1
  return (val - domain1) / (domain2 - domain1) * (range2 - range1) + range1
}

export function constrain (val, min, max) {
  return Math.min(Math.max(min, val), max)
}

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

export function getSampleData () {
  return [{"id":"GIUJ2WATF2HRCWBEKI34FUAVFZ4S7OHE","message_id":"PGMOCCWTIHR527UKISQ4EV4JYG55DJDU","request_id":"DZKYM26OJL2ENEQ5DG6VLDHMIXMTXQN2","input_id":"SCX6732QLCM7VVKUR3EBFPCF2MDJ5JPV","data":{"Geometry":{"coordinates":[38.90714645385742,-77.04476165771484],"type":"Point"},"date":1.484660606e+09,"type":"Feature"}},{"id":"PI5KOJMYQPMA6VPUEPYGNUB2K3MULR45","message_id":"XGYLMLC5HHQE7LR3RG43LV5M4ER75C25","request_id":"53OTLO72IIXCTMN55PKUH3ZNCMTOD5KM","input_id":"SCX6732QLCM7VVKUR3EBFPCF2MDJ5JPV","data":{"Geometry":{"coordinates":[38.90659713745117,-77.04178619384766],"type":"Point"},"date":1.484660868e+09,"type":"Feature"}},{"id":"UDTRVNZGAYTWEVJLJWV6NYG4SRFVGT4Z","message_id":"NLAJGOVGZQMCADGQEVD5JTXV3HAPBHR2","request_id":"SVYSPQGB37YPFNHI4U2COZEENV4O4ZEG","input_id":"SCX6732QLCM7VVKUR3EBFPCF2MDJ5JPV","data":{"Geometry":{"coordinates":[38.904422760009766,-77.04060363769531],"type":"Point"},"date":1.484661089e+09,"type":"Feature"}},{"id":"4HNYNXJVQYGCB6OEQVGH3AK6TKUQGZPS","message_id":"2TB3BCSU5Y6K2J3MTZ6NPZFFWTGLJOUK","request_id":"RXYNJ6IUP5FEAGXH5EAHUZMSOTPLPYI6","input_id":"SCX6732QLCM7VVKUR3EBFPCF2MDJ5JPV","data":{"Geometry":{"coordinates":[38.90373611450195,-77.03819274902344],"type":"Point"},"date":1.484661273e+09,"type":"Feature"}},{"id":"W6CTJG6AGPNFLJATQRSISBWUIDE7GEQR","message_id":"QUL6V5QSLHJ2Y5CGKYUSC3NM3APHQCBD","request_id":"3KGKE44IVXRQOJKZERTSU5EHO44GTT4U","input_id":"SCX6732QLCM7VVKUR3EBFPCF2MDJ5JPV","data":{"date":1.484665769e+09,"geometry":{"coordinates":[38.903324127197266,-77.03483581542969],"type":"Point"},"type":"Feature"}},{"id":"4UPKWQKXF6DXS2WELGOKQGN2A6EBQRCP","message_id":"XM7IS7ADOSYVGOOGG2VCX5FOAJ2ZLF4R","request_id":"GHJR5SSJXWYZPZTQZQOUHQ5X2ZNHLHL2","input_id":"SCX6732QLCM7VVKUR3EBFPCF2MDJ5JPV","data":{"date":1.484666355e+09,"geometry":{"coordinates":[38.90122985839844,-77.03242492675781],"type":"Point"},"type":"Feature"}},{"id":"IKQM4JVCYYRUUMKEJKU3337YGTYKCKDP","message_id":"XEY7FRYPKUOYMFGVI73Y7FMLPRNYG6MO","request_id":"HOJ6S3X7XU4HCNWPC3NFP5ZEMK57TMRO","input_id":"SCX6732QLCM7VVKUR3EBFPCF2MDJ5JPV","data":{"date":1.484666632e+09,"geometry":{"coordinates":[38.898929595947266,-77.0322036743164],"type":"Point"},"type":"Feature"}},{"id":"IEO4AHCV5YG66ISODOXEMRWLJSQKRSSE","message_id":"ULJEQPA5EICPMLW3OBQDBE6FWVITEXC2","request_id":"NWIWIZWHQU4U3ZLLB2OPK6RDKQ6MDXF7","input_id":"SCX6732QLCM7VVKUR3EBFPCF2MDJ5JPV","data":{"date":1.484666942e+09,"geometry":{"coordinates":[38.895423889160156,-77.03177642822266],"type":"Point"},"type":"Feature"}},{"id":"VWL2QQU4WHLVEXEJS44TYACS5NA2JDZC","message_id":"I4D4IKYJQ2H7PKM5U4QGFENZHJ7HEUZD","request_id":"O6T75DBTISNXCH6ZTWDBZ5EPYYCVUGUC","input_id":"SCX6732QLCM7VVKUR3EBFPCF2MDJ5JPV","data":{"date":1.484666943e+09,"geometry":{"coordinates":[38.895423889160156,-77.03177642822266],"type":"Point"},"type":"Feature"}},{"id":"ZUEARF42RIEPXS4WF2SLIJPBDX3AVBSH","message_id":"76QTIEBEWQJREJ4KKKP65KSYIKNN5XMR","request_id":"FG6AVUVLLSAHS5JIONXJVVNXPZRNI7EU","input_id":"SCX6732QLCM7VVKUR3EBFPCF2MDJ5JPV","data":{"date":1.484667222e+09,"geometry":{"coordinates":[38.8921012878418,-77.03201293945312],"type":"Point"},"type":"Feature"}},{"id":"F3I74F2IKMNYWC6UW5KXWHUD5BQRF4MR","message_id":"VRROPLG4GJB6QOSY4U6GKYH6RYCQDMF3","request_id":"R2UFNTJTV6GVY3S4T4ZZU24NFXRBUW7U","input_id":"SCX6732QLCM7VVKUR3EBFPCF2MDJ5JPV","data":{"date":1.484667222e+09,"geometry":{"coordinates":[38.8921012878418,-77.03201293945312],"type":"Point"},"type":"Feature"}},{"id":"N7UF6ZMMUV7AZ47QLNUFCUSIJ56Q4472","message_id":"P5D553WE2XDZZHONZQHSEHP4GTQMIKV2","request_id":"U37WZSLHHW7WZP33TAX2V3XJGMSIIKWC","input_id":"SCX6732QLCM7VVKUR3EBFPCF2MDJ5JPV","data":{"date":1.484667513e+09,"geometry":{"coordinates":[38.88996124267578,-77.03324127197266],"type":"Point"},"type":"Feature"}},{"id":"ANCFSUOIBB2X2O5HMDZD3EIBMFUVDSZV","message_id":"W2A6MND3WNAJIHK5IYUTVTFVDI3PH2FS","request_id":"QTORCPX2EXOLAHC5D4J3FJXTEYPO6QZY","input_id":"SCX6732QLCM7VVKUR3EBFPCF2MDJ5JPV","data":{"date":1.484667654e+09,"geometry":{"coordinates":[38.89020919799805,-77.03443908691406],"type":"Point"},"type":"Feature"}},{"id":"JJJAL37D3LGNKRFTOSC7RYPGEE4SZVB7","message_id":"X3DJJMWFFCKBMING2H2PF5E2XAQYQBF4","request_id":"MAFNBRJTDYG4ZQQ7OGPGO5I5YWOO27KZ","input_id":"SCX6732QLCM7VVKUR3EBFPCF2MDJ5JPV","data":{"date":1.484667654e+09,"geometry":{"coordinates":[38.89020919799805,-77.03443908691406],"type":"Point"},"type":"Feature"}},{"id":"CGSSZVB4FCFIBCKHE2MKTE2HS6SHCHS2","message_id":"LNJC7SAN2IEQC6Y6NBH7SBWT2T6LYN2H","request_id":"JKHV3NBO5FQZNH2JBHAUBRYGBI7WTFZ4","input_id":"SCX6732QLCM7VVKUR3EBFPCF2MDJ5JPV","data":{"date":1.484667941e+09,"geometry":{"coordinates":[38.88999557495117,-77.036865234375],"type":"Point"},"type":"Feature"}},{"id":"E33FUBH4K43OI5YC3ZLACVAK7N3OV3OT","message_id":"KBQ2VWSHDMTZ2SNHDA7BUMG4TKYOPXJ5","request_id":"JW5IJEOHG622F6MIHKVTI3DHMIUOM6YU","input_id":"SCX6732QLCM7VVKUR3EBFPCF2MDJ5JPV","data":{"date":1.484668077e+09,"geometry":{"coordinates":[38.889949798583984,-77.03812408447266],"type":"Point"},"type":"Feature"}},{"id":"XV34CCAQR4HS2QIMIZWDKXR755JL3TWR","message_id":"FRN7NUFLWUAXCTRPJEI5WSTQRSOWLANH","request_id":"TSI3UAFOXQJ27Q6PD3STXS7WHFMAMAA5","input_id":"SCX6732QLCM7VVKUR3EBFPCF2MDJ5JPV","data":{"date":1.484668214e+09,"geometry":{"coordinates":[38.88996887207031,-77.03929901123047],"type":"Point"},"type":"Feature"}},{"id":"4RVD5JHXAJYEIO2UJNATMHG5ZZ443ZAG","message_id":"ASORVOLEAPKL6GQG6UPFMWTO5L7CLNOJ","request_id":"UJDATIJRFIENQNHO4J6FXAVAA44HKAXL","input_id":"SCX6732QLCM7VVKUR3EBFPCF2MDJ5JPV","data":{"date":1.484668214e+09,"geometry":{"coordinates":[38.88996887207031,-77.03929901123047],"type":"Point"},"type":"Feature"}},{"id":"JSIDNPQVXZI76BUG2ZOK3CGORHHQFCWF","message_id":"NXXN7J5XB7TIN4YMQVICXCPJIJPXUMZY","request_id":"CJPM6TT4LFHFM3UOCRDZUSOQI3ADXDER","input_id":"SCX6732QLCM7VVKUR3EBFPCF2MDJ5JPV","data":{"date":1.484668532e+09,"geometry":{"coordinates":[38.88978576660156,-77.04022979736328],"type":"Point"},"type":"Feature"}},{"id":"ZRWV3QVJTNM7EBDV74J7GEQFCGXEFZKS","message_id":"ELPFS6NVTUZYNC73X5X3CWN4BIFJU5SU","request_id":"AE5BX4Q4HK4WDBZFFNUFM5ZY2Q3CU2RR","input_id":"SCX6732QLCM7VVKUR3EBFPCF2MDJ5JPV","data":{"date":1.48466868e+09,"geometry":{"coordinates":[38.889163970947266,-77.04084777832031],"type":"Point"},"type":"Feature"}},{"id":"PJUNXICZFRE43FEDY7QIH5QCAABKU2U5","message_id":"O6UXFWTZ3D6RV6S4LPPPT7OOCIW6CGZF","request_id":"G66JP4SOCJ2EN6FWWXW2YKV4DV6INDSQ","input_id":"SCX6732QLCM7VVKUR3EBFPCF2MDJ5JPV","data":{"date":1.484668818e+09,"geometry":{"coordinates":[38.88915252685547,-77.04115295410156],"type":"Point"},"type":"Feature"}}]
}

