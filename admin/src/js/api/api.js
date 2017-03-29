// @flow weak

import log from 'loglevel';

import { JWTAPIClient, APIError, AuthenticationError } from './base-api';
import type { SupportedMethods } from './base-api';

import type {
  APIErrors,
  APIUser,
  APIUsers,
  APIBaseUser,
  APINewUser,
  APIProject,
  APINewProject,
  APIProjects,
  APIExpedition,
  APINewExpedition,
  APIExpeditions,
  APITwitterInput,
  APITwitterInputCreateResponse,
  APIInputs,
  APITeam,
  APINewTeam,
  APITeams,
  APIBaseMember,
  APIMember,
  APINewMember,
  APIMembers,
  APIAdministrator,
  APINewAdministrator,
  APIAdministrators
} from './types';

export type FKAPIOKResponse<T> = {
  type: 'ok';
  payload: T;
};
export type FKAPIErrResponse = {
  type: 'err';
  errors: APIErrors;
}
export type FKAPIResponse<T> = FKAPIOKResponse<T> | FKAPIErrResponse;

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
      return { type: 'ok', payload: res };
    } catch (e) {
      if (e instanceof AuthenticationError) {
        const APIFakeAuthError: APIErrors = {
          code: 'AuthenticationError',
          detail: 'Unauthorized',
          id: '',
          meta: {},
          status: 401
        };
        return { type: 'err', errors: APIFakeAuthError };
      } else if (e instanceof APIError) {
        const APIFakeOtherError: APIErrors = {
          code: 'UnknownAPIError',
          detail: e.body || '',
          id: '',
          meta: {},
          status: 500
        }

        try {
          return { type: 'err', errors: APIFakeOtherError };
        } catch (e2) {
          return { type: 'err', errors: APIFakeOtherError };
        }
      } else {
        const APIFakeOtherError: APIErrors = {
          code: 'UnknownError',
          detail: e.msg,
          id: '',
          meta: {},
          status: 500
        }

        return { type: 'err', errors: APIFakeOtherError };
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

  patchWithErrors<T>(endpoint: string, values?: Object): Promise<FKAPIResponse<T>> {
    return this.execWithErrors(this.patchJSON(endpoint, values));
  }

  signUp(u: APINewUser): Promise<FKAPIResponse<APIUser>> {
    return this.postWithErrors('/users', u);
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

  updateUserById(userId: number, u: APIBaseUser): Promise<FKAPIResponse<?APIUser>> {
    return this.patchWithErrors(`/users/${userId}`, u);
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

  getCurrentUserProjects(): Promise<FKAPIResponse<APIProjects>> {
    return this.getWithErrors('/user/projects')
  }

  getProjectBySlug(slug: string): Promise<FKAPIResponse<APIProject>> {
    return this.getWithErrors(`/projects/@/${slug}`)
  }

  createProject(values: APINewProject): Promise<FKAPIResponse<APIProject>> {
    return this.postWithErrors('/projects', values)
  }

  updateProject(projectId: number, values: APINewProject): Promise<FKAPIResponse<APIProject>> {
    return this.patchWithErrors(`/projects/${projectId}`, values)
  }

  getExpeditionsByProjectSlug(projectSlug: string): Promise<FKAPIResponse<APIExpeditions>> {
    return this.getWithErrors(`/projects/@/${projectSlug}/expeditions`)
  }

  getExpeditionBySlugs(projectSlug: string, expeditionSlug: string): Promise<FKAPIResponse<APIExpedition>> {
    return this.getWithErrors(`/projects/@/${projectSlug}/expeditions/@/${expeditionSlug}`)
  }

  createExpedition(projectId: number, values: APINewExpedition): Promise<FKAPIResponse<APIExpedition>> {
    return this.postWithErrors(`/projects/${projectId}/expeditions`, values)
  }

  updateExpedition(expeditionId: number, values: APINewExpedition): Promise<FKAPIResponse<APIExpedition>> {
    return this.patchWithErrors(`/expeditions/${expeditionId}`, values)
  }

  getExpeditionInputs(expeditionId: number): Promise<FKAPIResponse<APIInputs>> {
    return this.getWithErrors(`/expeditions/${expeditionId}/inputs`)
  }

  getInputsBySlugs(projectSlug: string, expeditionSlug: string): Promise<FKAPIResponse<APIInputs>> {
    return this.getWithErrors(`/projects/@/${projectSlug}/expeditions/@/${expeditionSlug}/inputs`)
  }

  getTwitterInput(inputId: number): Promise<FKAPIResponse<APITwitterInput>> {
    return this.getWithErrors(`/inputs/twitter-accounts/${inputId}`)
  }

  createTwitterInput(expeditionId: number): Promise<FKAPIResponse<APITwitterInputCreateResponse>> {
    return this.postWithErrors(`/expeditions/${expeditionId}/inputs/twitter-accounts`)
  }

  getTeamsBySlugs(projectSlug: string, expeditionSlug: string): Promise<FKAPIResponse<APITeams>> {
    return this.getWithErrors(`/projects/@/${projectSlug}/expeditions/@/${expeditionSlug}/teams`)
  }

  getTeamBySlugs(projectSlug: string, expeditionSlug: string, teamSlug: string): Promise<FKAPIResponse<APITeam>> {
    return this.getWithErrors(`/projects/@/${projectSlug}/expedition/@/${expeditionSlug}/teams/@/${teamSlug}`)
  }

  createTeam(expeditionId: number, values: APINewTeam): Promise<FKAPIResponse<APITeam>> {
    return this.postWithErrors(`/expeditions/${expeditionId}/teams`, values)
  }

  deleteTeam(teamId: number): Promise<FKAPIResponse<APIMember>> {
    return this.delWithErrors(`/teams/${teamId}`)
  }

  updateTeam(teamId: number, values: APINewTeam): Promise<FKAPIResponse<APINewTeam>> {
    return this.patchWithErrors(`/teams/${teamId}`, values)
  }

  addAdministrator(projectId: number, values: APINewAdministrator): Promise<FKAPIResponse<APIAdministrator>> {
    return this.postWithErrors(`/projects/${projectId}/administrators`, values)
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
    return this.postWithErrors(`/teams/${teamId}/members`, values)
  }

  updateMember(teamId: number, userId: number, values: APIBaseMember): Promise<FKAPIResponse<APIMember>> {
    return this.postWithErrors(`/teams/${teamId}/member/${userId}`, values)
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
