// @flow weak

import { JWTAPIClient, APIError, AuthenticationError } from './base-api';
import type {
  APIErrors,
  APIUser,
  APIUsers,
  APINewUser,
  APIProject,
  APINewProject,
  APIProjects,
  APIExpedition,
  APINewExpedition,
  APIExpeditions,
  APITeam,
  APINewTeam,
  APITeams
} from './types';

export type FKAPIResponse<T> = {
  type: 'ok' | 'err';
  payload?: T;
  errors?: APIErrors;
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

  async execWithErrors<T>(p: Promise<any>, parseJSON = false): Promise<FKAPIResponse<T>> {
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
        if (e.body && parseJSON) {
          return { type: 'err', errors: JSON.parse(e.body) };
        } else {
          return { type: 'err', errors: e.msg };
        }
      } else {
        return { type: 'err', raw: e.body };
      }
    }
  }

  postWithErrors<T>(endpoint: string, values?: Object): Promise<FKAPIResponse<T>> {
    return this.execWithErrors(this.postJSON(endpoint, values));
  }

  getWithErrors<T>(endpoint: string, values?: Object): Promise<FKAPIResponse<T>> {
    return this.execWithErrors(this.getJSON(endpoint, values));
  }

  signUp(email: string, username: string, password: string, invite_token: string): Promise<FKAPIResponse<APIUser>> {
    return this.postWithErrors('/user', { email, username, password, invite_token });
  }

  async signIn(username, password): Promise<FKAPIResponse<void>> {
    const response = await this.postWithErrors('/login', { username, password });
    if (response.type === 'ok') {
      this.onSignin();
    }
    return response;
  }

  async signOut(): Promise<void> {
    await this.postForm('/logout')
    this.onSignout();
  }

  getCurrentUser(): Promise<FKAPIResponse<APIUser>> {
    return this.getWithErrors('/user');
  }

  getUserById(userId: number): Promise<FKAPIResponse<?APIUser>> {
    return this.getWithErrors(`/users/${userId}`);
  }

  getUserByUsername(username: string): Promise<FKAPIResponse<?APIUser>> {
    return this.getWithErrors(`/users/@/${username}`);
  }

  getUsers(): Promise<FKAPIResponse<APIUsers>> {
    return this.getWithErrors('/users');
  }

  getProjects(): Promise<FKAPIResponse<APIProjects>> {
    return this.getWithErrors('/projects')
  }

  getProjectBySlug(slug: string): Promise<FKAPIResponse<APIProject>> {
    return this.getWithErrors(`/projects/@/${slug}`)
  }

  createProject(values: APINewProject): Promise<FKAPIResponse<APIProject>> {
    return this.postWithErrors('/project', values)
  }

  updateProject(projectId: number, values: APINewProject): Promise<FKAPIResponse<APIProject>> {
    return this.postWithErrors(`/projects/${projectId}`, values)
  }

  getExpeditionsByProjectSlug(projectSlug: string): Promise<FKAPIResponse<APIExpeditions>> {
    return this.getWithErrors(`/projects/@/${projectSlug}/expeditions`)
  }

  getExpeditionBySlugs(projectSlug: string, expeditionSlug: string): Promise<FKAPIResponse<APIExpedition>> {
    return this.getWithErrors(`/projects/@/${projectSlug}/expeditions/@/${expeditionSlug}`)
  }

  createExpedition(projectId: number, values: APINewExpedition): Promise<FKAPIResponse<APIExpedition>> {
    return this.postWithErrors(`/projects/${projectId}/expedition`, values)
  }

  updateExpedition(expeditionId: number, values: APINewExpedition): Promise<FKAPIResponse<APIExpedition>> {
    return this.postWithErrors(`/expeditions/${expeditionId}`, values)
  }  

  getTeamsBySlugs(projectSlug: string, expeditionSlug: string): Promise<FKAPIResponse<APITeams>> {
    return this.getWithErrors(`/projects/@/${projectSlug}/expeditions/@/${expeditionSlug}/teams`)
  }

  getTeamBySlugs(projectSlug: string, expeditionSlug: string, teamSlug: string): Promise<FKAPIResponse<APITeam>> {
    return this.getWithErrors(`/projects/@/${projectSlug}/expedition/@/${expeditionSlug}/teams/@/${teamSlug}`)
  }

  createTeam(expeditionId: number, values: APINewTeam): Promise<FKAPIResponse<APITeam>> {
    return this.postWithErrors(`/expeditions/${expeditionId}/team`, values)
  }
}
