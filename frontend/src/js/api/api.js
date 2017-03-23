
import 'whatwg-fetch'
import {BaseError} from '../utils.js'
class APIError extends BaseError {}
class AuthenticationError extends APIError {}
const SIGNED_IN_KEY = 'signedIn'

class FKApiClient {
  constructor(baseUrl) {
    this.baseUrl = baseUrl
  }

  async get(path, params) {
    try {
      const url = new URL(path, this.baseUrl)
      if (params) {
        for (const key in params) {
          url.searchParams.set(key, params[key])
        }
      }
      let res
      try {
        res = await fetch(url.toString(), {
          method: 'GET',
          credentials: 'include'
        })
      } catch (e) {
        console.log('Threw while GETing', url.toString(), e)
        throw new APIError('HTTP error')
      }
      if (res.status == 404) {
        console.log('not found', url.toString(), await res.text())
        throw new APIError('not found')
      }
      if (res.status == 401) {
        console.log('Bad auth while GETing', url.toString(), await res.text())
        throw new AuthenticationError()
      } else if (!res.ok) {
        const err = await res.json()
        console.log('Non-OK response while GETing', url.toString(), err)
        throw new APIError(err)
      }
      return res
    } catch (e) {
      if (e instanceof AuthenticationError) {
        this.onAuthError(e)
      }
      throw e
    }
  }

  async getJSON(path: string, params?: Object): Promise<any> {
    const res = await this.get(path, params)
    return res.json()
  }

  async post(path, body) {
    try {
      const url = new URL(path, this.baseUrl)
      let res
      try {
        res = await fetch(url.toString(), {
          method: 'POST',
          credentials: 'include',
          body
        })
      } catch (e) {
        console.log('Threw while POSTing', url.toString(), e)
        throw new APIError('HTTP error')
      }
      if (res.status == 401) {
        console.log('Bad auth while POSTing', url.toString(), await res.text())
        throw new AuthenticationError()
      } else if (!res.ok) {
        const err = await res.json()
        console.log('Non-OK response while POSTing', url.toString(), err)
        throw new APIError(err)
      }
      return res
    } catch (e) {
      if (e instanceof AuthenticationError) {
        this.onAuthError(e)
      }
      throw e
    }
  }

  async postForm(path, body) {
    const data = new FormData()
    if (body) {
      for (const key in body) {
        data.append(key, body[key])
      }
    }
    const res = await this.post(path, data)
    return res.text()
  }

  async getProject (project) {
    const res = await this.getJSON('/api/project/' + projectID)
    return res
  }

  async getExpeditions (projectID) {
    const res = await this.getJSON('/api/project/' + projectID + '/expeditions')
    return res
  }

  async getExpedition (projectID, expeditionID) {
    const res = await this.getJSON('/api/project/' + projectID + '/expedition/' + expeditionID)
    return res
  }

  async getDocuments (projectID, expeditionID) {
    const res = await this.getJSON('/api/project/' + projectID + '/expedition/' + expeditionID + '/documents')
    return res
  }
}

// TODO switch back as soon as we can collect data locally
// const hostname = location.hostname.split('.')[location.hostname.split('.').length-1] === 'localhost' ? 'http://localhost:8080' : 'https://fieldkit.org'
const hostname = 'https://fieldkit.org'
export default new FKApiClient(hostname)