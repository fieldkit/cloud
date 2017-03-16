// @flow weak

import {AuthAPIClient, APIError} from './base-api';

import type { ErrorMap } from '../common/util';

export type FKAPIResponse = {
  type: 'ok' | 'err';
  payload?: any;
  errors?: ErrorMap;
  raw?: string;
}

let apiClientInstance;
export class FKApiClient extends AuthAPIClient {
  signinCb: ?() => void;
  signoutCb: ?() => void;

  static setup(
    baseUrl: string,
    unauthorizedHandler: () => void,
    { onSignin, onSignout }: { onSignin?: () => void, onSignout?: () => void } = {}
  ): FKApiClient
  {
    if (!apiClientInstance) {
      apiClientInstance = new FKApiClient(baseUrl, unauthorizedHandler);
      apiClientInstance.signinCb = onSignin;
      apiClientInstance.signoutCb = onSignout;
    }

    return apiClientInstance;
  }

  onSignin() {
    super.onSignin();
    if (this.signinCb) {
      this.signinCb();
    }
  }

  onSignout() {
    super.onSignout();
    if (this.signoutCb) {
      this.signoutCb();
    }
  }

  static get(): FKApiClient {
    if (!apiClientInstance) {
      throw new APIError('API has not been set up!');
    }

    return apiClientInstance;
  }

  async execWithJSONErrors(p: Promise<any>, parseJSON = false): Promise<FKAPIResponse> {
    try {
      const res = await p;
      if (res) {
        if (parseJSON) {
          return { type: 'ok', payload: JSON.parse(res) };
        } else {
          return { type: 'ok', payload: res };
        }
      } else {
        return { type: 'ok' }
      }
    } catch (e) {
      if (e instanceof APIError) {
        return { type: 'err', errors: JSON.parse(e.body) };
      } else {
        return { type: 'err', raw: e.body };
      }
    }
  }

  postFormWithJSONErrors(endpoint: string, values?: Object): Promise<FKAPIResponse> {
    return this.execWithJSONErrors(this.postForm(endpoint, values), true);
  }

  postJSONWithJSONErrors(endpoint: string, values?: Object): Promise<FKAPIResponse> {
    return this.execWithJSONErrors(this.postJSON(endpoint, values));
  }

  getWithJSONErrors(endpoint: string, values?: Object): Promise<FKAPIResponse> {
    return this.execWithJSONErrors(this.getJSON(endpoint, values));
  }

  signUp(email: string, username: string, password: string, invite: string): Promise<FKAPIResponse> {
    return this.postFormWithJSONErrors('/api/user/sign-up', { email, username, password, invite });
  }

  async signIn(username, password): Promise<FKAPIResponse> {
    const response = await this.postFormWithJSONErrors('/api/user/sign-in', { username, password });
    if (response.type === 'ok') {
      this.onSignin();
    }
    return response;
  }

  async signOut(): Promise<void> {
    await this.postForm('/api/user/sign-out')
    this.onSignout();
  }

  getUser(): Promise<FKAPIResponse> {
    return this.getWithJSONErrors('/api/user/current')
  }

  getProjects(): Promise<FKAPIResponse> {
    return this.getWithJSONErrors('/api/projects')
  }

  getProjectBySlug(slug: string): Promise<FKAPIResponse> {
    return this.getWithJSONErrors(`/api/project/${slug}`)
  }

  createProject(name, description): Promise<FKAPIResponse> {
    return this.postFormWithJSONErrors('/api/projects/add', { name, description })
  }

  getExpeditionsByProjectSlug(projectSlug: string): Promise<FKAPIResponse> {
    return this.getWithJSONErrors(`/api/project/${projectSlug}/expeditions`)
  }

  getExpeditionBySlugs(projectSlug: string, expeditionSlug: string): Promise<FKAPIResponse> {
    return this.getWithJSONErrors(`/api/project/${projectSlug}/expedition/${expeditionSlug}`)
  }

  async createExpedition(projectID: string, name: string, description: string): Promise<FKAPIResponse> {
    return this.postFormWithJSONErrors(`/api/project/${projectID}/expeditions/add`, { name, description })
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
