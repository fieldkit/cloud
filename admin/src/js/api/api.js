// @flow weak

import { JWTAPIClient, APIError, AuthenticationError } from './base-api';

import type { ErrorMap } from '../common/util';

export type FKAPIResponse = {
  type: 'ok' | 'err';
  payload?: any;
  errors?: ErrorMap;
  raw?: string;
}

let apiClientInstance;
export class FKApiClient extends JWTAPIClient {
  signinCb: ?() => void;
  signoutCb: ?() => void;
  unauthorizedHandler: ?() => void;

  static setup(
    baseUrl: string,
    unauthorizedHandler: () => void,
    { onSignin, onSignout }: { onSignin?: () => void, onSignout?: () => void } = {}
  ): FKApiClient
  {
    if (!apiClientInstance) {
      apiClientInstance = new FKApiClient(baseUrl);
      apiClientInstance.signinCb = onSignin;
      apiClientInstance.signoutCb = onSignout;
      apiClientInstance.unauthorizedHandler = unauthorizedHandler;
    }

    return apiClientInstance;
  }

  onSignin() {
    if (this.signinCb) {
      this.signinCb();
    }
  }

  onSignout() {
    this.clearJWT();

    if (this.signoutCb) {
      this.signoutCb();
    }
  }

  signedIn(): boolean {
    return this.loadJWT() != null;
  }

  onAuthError(e: AuthenticationError) {
    if (this.signedIn()) {
      if (this.unauthorizedHandler) {
        this.unauthorizedHandler();
      }
      this.onSignout();
    }
  }

  async post(path: string, body?: FormData | string | null = null, headers: Object = {}): Promise<any> {
    try {
      return await super.post(path, body, headers);
    } catch (e) {
      if (e instanceof AuthenticationError) {
        this.onAuthError(e);
      }

      throw e;
    }
  }

  async get(path: string, params: Object = {}, headers: Object = {}): Promise<any> {
    try {
      return await super.get(path, params, headers);
    } catch (e) {
      if (e instanceof AuthenticationError) {
        this.onAuthError(e);
      }

      throw e;
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
        if (e.body) {
          return { type: 'err', errors: JSON.parse(e.body) };
        } else {
          return { type: 'err', errors: e.msg };
        }
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

  signUp(email: string, username: string, password: string, invite_token: string): Promise<FKAPIResponse> {
    return this.postJSONWithJSONErrors('/user', { email, username, password, invite_token });
  }

  async signIn(username, password): Promise<FKAPIResponse> {
    const response = await this.postJSONWithJSONErrors('/login', { username, password });
    if (response.type === 'ok') {
      this.onSignin();
    }
    return response;
  }

  async signOut(): Promise<void> {
    await this.postForm('/logout')
    this.onSignout();
  }

  getUser(): Promise<FKAPIResponse> {
    return this.getWithJSONErrors('/user')
  }

  getProjects(): Promise<FKAPIResponse> {
    return this.getWithJSONErrors('/projects')
  }

  getProjectBySlug(slug: string): Promise<FKAPIResponse> {
    return this.getWithJSONErrors(`/projects/@/${slug}`)
  }

  createProject(values: { name: string, slug: string, description: string }): Promise<FKAPIResponse> {
    return this.postJSONWithJSONErrors('/project', values)
  }

  updateProject(id: string, values: { name: string, slug: string, description: string }): Promise<FKAPIResponse> {
    return this.postJSONWithJSONErrors(`/projects/${id}`, values)
  }

  getExpeditionsByProjectSlug(projectSlug: string): Promise<FKAPIResponse> {
    return this.getWithJSONErrors(`/projects/@/${projectSlug}/expeditions`)
  }

  getExpeditionBySlugs(projectSlug: string, expeditionSlug: string): Promise<FKAPIResponse> {
    return this.getWithJSONErrors(`/projects/@/${projectSlug}/expedition/@/${expeditionSlug}`)
  }

  createExpedition(projectId: number, values: { name: string, slug: string, description: string }): Promise<FKAPIResponse> {
    return this.postJSONWithJSONErrors(`/projects/${projectId}/expedition`, values)
  }

  getTeamsBySlugs(projectSlug: string, expeditionSlug: string): Promise<FKAPIResponse> {
    return this.getWithJSONErrors(`/projects/@/${projectSlug}/expeditions/@/${expeditionSlug}`)
  }

  getTeamBySlugs(projectSlug: string, expeditionSlug: string, teamSlug: string): Promise<FKAPIResponse> {
    return this.getWithJSONErrors(`/projects/@/${projectSlug}/expedition/@/${expeditionSlug}/teams/@/${teamSlug}`)
  }

  createTeam(expeditionId: number, values: { name: string, slug: string, description: string }): Promise<FKAPIResponse> {
    return this.postJSONWithJSONErrors(`/expeditions/${expeditionId}/team`, values)
  }
}
