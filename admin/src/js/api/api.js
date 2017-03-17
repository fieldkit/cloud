// @flow weak

import { JWTAPIClient, APIError, AuthenticationError } from './base-api';
import type { SupportedMethods } from './base-api';

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
  APIInput,
  APINewInput,
  APIInputs,
  APITeam,
  APINewTeam,
  APITeams,
  APIMember,
  APINewMember,
  APIMembers,
  APIAdministrator,
  APINewAdministrator,
  APIAdministrators
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

  static get(): FKApiClient {
    if (!apiClientInstance) {
      throw new APIError('API has not been set up!');
    }

    return apiClientInstance;
  }

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

  async exec(
    method: SupportedMethods,
    path: string,
    { params, body, headers = {} }: {
      params?: Object,
      body?: ?(Blob | FormData | URLSearchParams | string),
      headers: Object
    } = {}
  ): Promise<Response> {
    try {
      return await super.exec(method, path, { params, body, headers });
    } catch (e) {
      if (e instanceof AuthenticationError) {
        this.onAuthError(e);
      }

      throw e;
    }
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

  delWithErrors<T>(endpoint: string, values?: Object): Promise<FKAPIResponse<T>> {
    return this.execWithErrors(this.delJSON(endpoint, values));
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
    return this.getWithErrors(`/projects/@/${projectSlug}/expedition/@/${expeditionSlug}`)
  }

  createExpedition(projectId: number, values: APINewExpedition): Promise<FKAPIResponse<APIExpedition>> {
    return this.postWithErrors(`/projects/${projectId}/expedition`, values)
  }

  getInputsByProjectSlug(projectSlug: string): Promise<FKAPIResponse<APIInputs>> {
    return this.getWithErrors(`/projects/@/${projectSlug}/inputs`)
  }

  getInputBySlugs(projectSlug: string, inputSlug: string): Promise<FKAPIResponse<APIInput>> {
    return this.getWithErrors(`/projects/@/${projectSlug}/inputs/@/${inputSlug}`)
  }

  getInput(inputId: number): Promise<FKAPIResponse<APIInput>> {
    return this.getWithErrors(`/inputs/${inputId}`)
  }

  createInput(projectId: number, values: APINewInput): Promise<FKAPIResponse<APIInput>> {
    return this.postWithErrors(`/projects/${projectId}/input`, values)
  }

  getTeamsBySlugs(projectSlug: string, expeditionSlug: string): Promise<FKAPIResponse<APITeams>> {
    return this.getWithErrors(`/projects/@/${projectSlug}/expeditions/@/${expeditionSlug}`)
  }

  getTeamBySlugs(projectSlug: string, expeditionSlug: string, teamSlug: string): Promise<FKAPIResponse<APITeam>> {
    return this.getWithErrors(`/projects/@/${projectSlug}/expedition/@/${expeditionSlug}/teams/@/${teamSlug}`)
  }

  createTeam(expeditionId: number, values: APINewTeam): Promise<FKAPIResponse<APITeam>> {
    return this.postWithErrors(`/expeditions/${expeditionId}/team`, values)
  }

  addAdministrator(projectId: number, values: APINewAdministrator): Promise<FKAPIResponse<APIAdministrator>> {
    return this.postWithErrors(`/projects/${projectId}/administrator`, values)
  }

  getAdministrators(projectId: number): Promise<FKAPIResponse<APIAdministrators>> {
    return this.getWithErrors(`/projects/${projectId}/administrators`)
  }

  getAdministratorsBySlug(projectSlug: string): Promise<FKAPIResponse<APIAdministrators>> {
    return this.getWithErrors(`/projects/@/${projectSlug}/administrators`)
  }

  deleteAdministrator(projectId: number, userId: number): Promise<FKAPIResponse<APIAdministrator>> {
    return this.delWithErrors(`/projects/${projectId}/administrators/${userId}`)
  }

  addMember(teamId: number, values: APINewMember): Promise<FKAPIResponse<APIMember>> {
    return this.postWithErrors(`/teams/${teamId}/member`, values)
  }

  getMembers(teamId: number): Promise<FKAPIResponse<APIMembers>> {
    return this.getWithErrors(`/teams/${teamId}/members`)
  }

  getMembersBySlugs(projectSlug: string, expeditionSlug: string, teamSlug: string): Promise<FKAPIResponse<APIMembers>> {
    return this.getWithErrors(`/projects/@/${projectSlug}/expedition/@/${expeditionSlug}/teams/@/${teamSlug}/members`)
  }

  deleteMember(teamId: number, userId: number): Promise<FKAPIResponse<APIMember>> {
    return this.delWithErrors(`/teams/${teamId}/members/${userId}`)
  }
}
