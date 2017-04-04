
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
          method: 'GET'
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

  async getProject (projectID) {
    const res = await this.getJSON('/projects/@/' + projectID)
    return res
  }

  async getExpeditions (projectID) {
    const res = await this.getJSON('/projects/@/' + projectID + '/expeditions')
    return res.expeditions
  }

  async getExpedition (projectID, expeditionID) {
    const res = await this.getJSON('/projects/@/' + projectID + '/expeditions/@/' + expeditionID)
    return res
  }

  async getDocuments (projectID, expeditionID) {
    const res = await this.getJSON('/projects/@/' + projectID + '/expeditions/@/' + expeditionID + '/documents')
    return res.documents
  }
}

let API_HOST = 'https://fieldkit.org';
if (process.env.NODE_ENV === 'development') {
  API_HOST = 'http://localhost:8080';
} else if (process.env.NODE_ENV === 'staging' || window.location.hostname.endsWith('fieldkit.team')) {
  API_HOST = 'https://api.fieldkit.team';
}

export default new FKApiClient(API_HOST)