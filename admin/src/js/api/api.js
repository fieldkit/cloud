// @flow weak

import {AuthAPIClient, APIError} from './base-api';

import type { ErrorMap } from '../common/util';

let apiClientInstance;
export class FKApiClient extends AuthAPIClient {
  static setup(baseUrl: string, unauthorizedHandler: () => void): FKApiClient {
    if (!apiClientInstance) {
      apiClientInstance = new FKApiClient(baseUrl, unauthorizedHandler);
    }

    return apiClientInstance;
  }

  static get(): FKApiClient {
    if (!apiClientInstance) {
      throw new APIError('API has not been set up!');
    }

    return apiClientInstance;
  }

  async postFormWithJSONErrors(endpoint: string, values: Object): Promise<?ErrorMap> {
    try {
      await this.postForm(endpoint, values);
      return null;
    } catch (e) {
      if (e instanceof APIError) {
        return JSON.parse(e.body);
      } else {
        return { error: e.msg };
      }
    }
  }

  signUp(email: string, username: string, password: string, invite: string): Promise<?ErrorMap> {
    return this.postFormWithJSONErrors('/api/user/sign-up', { email, username, password, invite });
  }

  async signIn(username, password): Promise<?ErrorMap> {
    const errors = await this.postFormWithJSONErrors('/api/user/sign-in', { username, password });
    if (errors) {
      return errors;
    } else {
      this.onSignin();
      return null;
    }
  }

  async signOut() {
    await this.postForm('/api/user/sign-out')
    this.onSignout()
  }

  async getUser() {
    const res = await this.postJSON('/api/user/current')
    return res
  }

  async getProjects () {
    const res = await this.postJSON('/api/projects')
    return res
  }

  async createProjects (name) {
    const res = await this.postJSON('/api/projects/add', { name })
    return res
  }

  async getExpeditions (projectID) {
    const res = await this.getJSON(`/api/project/${projectID}/expeditions`)
    return res
  }

  async postGeneralSettings (projectID, name) {
    const res = await this.postJSON(`/api/project/${projectID}/expeditions/add`, { name })
    return res
  }

  async postInputs (projectID, expeditionID, name) {
    const res = await this.postJSON(`/api/project/${projectID}/expedition/${expeditionID}/inputs/add`, { name })
    return res
  }

  async addExpeditionToken (projectID, expeditionID) {
    const res = await this.postJSON(`/api/project/${projectID}/expedition/${expeditionID}/tokens/add`)
    return res
  }

  async addInput(projectID, expeditionID, name) {
    const res = await this.postJSON(`/api/project/${projectID}/expedition/${expeditionID}/inputs/add`, { name })
    return res
  }

  async getInputs(projectID, expeditionID) {
    const res = await this.getJSON(`/api/project/${projectID}/expedition/${expeditionID}/inputs`)
    return res
  }
}
