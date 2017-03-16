
const isOriginLocalhost = location.hostname.split('.')[location.hostname.split('.').length-1] === 'localhost'
export const protocol = isOriginLocalhost ? 'http://' : 'https://'
export const hostname = isOriginLocalhost ? 'localhost:8080' : 'fieldkit.org'
