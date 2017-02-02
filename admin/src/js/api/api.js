
import 'whatwg-fetch'
import {BaseError} from '../utils.js'
import { protocol, hostname } from '../constants/APIBaseURL'
class APIError extends BaseError {}
class AuthenticationError extends APIError {}
const SIGNED_IN_KEY = 'signedIn'

class FKApiClient {

  onAuthError(e) {
    if (this.signedIn) {
      this.onSignOut()
    }
  }

  onSignIn() {
    localStorage.setItem(SIGNED_IN_KEY, 'signedIn')
  }

  onSignOut() {
    localStorage.removeItem(SIGNED_IN_KEY)
  }

  signedIn() {
    return localStorage.getItem(SIGNED_IN_KEY) != null
  }

  async get(path, params) {
    try {
      const url = new URL(path, protocol + hostname)
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
      const url = new URL(path, protocol + hostname)
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

  async signUp(params) {
    await this.postForm('/api/user/sign-up', params)
  }

  async signIn(username, password) {
    await this.postForm('/api/user/sign-in', { username, password })
  }

  async signOut() {
    await this.postForm('/api/user/sign-out')
    this.onSignOut()
  }

  async getProjects () {
    const res = await this.getJSON('/api/projects')
    return res
  }

  async createProjects (name) {
    const res = await this.getJSON('/api/projects/add?name=' + name)
    return res
  }

  async getExpeditions (projectID) {
    const res = await this.getJSON('/api/project/' + projectID + '/expeditions')
    return res
  }

  async postGeneralSettings (projectID, expeditionName) {
    const res = await this.getJSON('/api/project/' + projectID + '/expeditions/add?name=' + expeditionName)
    return res
  }

  async postInputs (projectID, expeditionID, inputName) {
    const res = await this.getJSON('/api/project/' + projectID + '/expedition/' + expeditionID + '/inputs/add?name=' + inputName)
    return res
  }

  async addExpeditionToken (projectID, expeditionID) {
    const res = await this.getJSON('/api/project/' + projectID + '/expedition/' + expeditionID + '/tokens/add')
    return res
  }

  async addInput (projectID, expeditionID, inputName) {
    const res = await this.getJSON('/api/project/' + projectID + '/expedition/' + expeditionID + '/inputs/add?name=' + inputName)
    return res
  }
}

export default new FKApiClient()