// @flow

import 'whatwg-fetch'
import {BaseError} from '../common/util'
import log from 'loglevel';
import jwtDecode from 'jwt-decode';

const JWT_KEY = 'jwt';
const JWT_EXPIRATION_TIME_KEY = 'jwt-expiration';

export class APIError extends BaseError {
  response: ?Response;
  body: ?string;

  constructor(msg: string, res: ?Response, body: ?string) {
    super(msg);
    this.response = res;
    this.body = body;
  }
}
export class AuthenticationError extends APIError {
  constructor(res: ?Response, body: ?string) {
    super('Unauthorized', res, body);
  }
}

export class APIClient {
  baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  async post(path: string, body?: FormData | string | null, headers: Object = {}): Promise<Response> {
    const url = new URL(path, this.baseUrl);

    log.info('POST', path, body, headers);

    let res;
    try {
      res = await fetch(url.toString(), {
        method: 'POST',
        body,
        headers
      });
    } catch (e) {
      log.error(e);
      log.error('Threw while POSTing', url.toString());
      throw new APIError('HTTP error', res);
    }

    if (!res.ok) {
      const body = await res.text();
      if (res.status === 401) {
        log.error('Bad auth while POSTing', url.toString(), body);
        throw new AuthenticationError(res, body);
      } else {
        log.error('Non-OK response while POSTing', url.toString(), body);
        throw new APIError('HTTP error', res, body);
      }
    }

    return res;
  }

  async get(path: string, params?: Object, headers: Object = {}): Promise<Response> {
    const url = new URL(path, this.baseUrl);

    if (params) {
      for (const key in params) {
        url.searchParams.set(key, params[key]);
      }
    }

    log.info('GET', path, params);

    let res;
    try {
      res = await fetch(url.toString(), {
        method: 'GET',
        headers
      });
    } catch (e) {
      log.error(e);
      log.error('Threw while GETing', url.toString());
      throw new APIError('HTTP error', res);
    }

    if (!res.ok) {
      const body = await res.text();
      if (res.status === 401) {
        log.error('Bad auth while GETing', url.toString(), body);
        throw new AuthenticationError(res, body);
      } else {
        log.error('Non-OK response while GETing', url.toString(), body);
        throw new APIError('HTTP error', res, body);
      }
    }

    return res;
  }

  async postForm(path: string, body?: Object, headers: Object = {}): Promise<string> {
    const data = new FormData();

    if (body) {
      for (const key in body) {
        data.append(key, body[key]);
      }
    }

    const res = await this.post(path, data, headers);
    return res.text();
  }

  async postJSON(path: string, body?: Object, headers: Object = {}): Promise<any> {
    const res = await this.post(path, JSON.stringify(body), headers);
    if (res.status === 204) {
      return null
    }
    return res.json();
  }

  async getText(path: string, params?: Object, headers: Object = {}): Promise<string> {
    const res = await this.get(path, params, headers);
    return res.text();
  }

  async getJSON(path: string, params?: Object, headers: Object = {}): Promise<any> {
    const res = await this.get(path, params, headers);
    return res.json();
  }
}

export class JWTAPIClient extends APIClient {
  async refreshJWT(): Promise<void> {
    const jwt = this.loadDecodedJWT();
    if (!jwt) return;

    const refreshToken: string = jwt.refresh_token;
    try {
      await super.postJSON('/refresh', { payload: { refresh_token: refreshToken } }, { Authorization: null });
    } catch (e) {
      console.error('JWT refresh error', e);
      this.clearJWT();
      throw e;
    }
  }

  loadJWT(): ?string {
    return localStorage.getItem(JWT_KEY);
  }

  loadJWTExpiration(): ?number {
    const expStr = localStorage.getItem(JWT_EXPIRATION_TIME_KEY);
    if (!expStr) {
      return null;
    }

    return parseInt(expStr, 10);
  }

  loadDecodedJWT(): ?Object {
    const jwt = this.loadJWT();
    if (!jwt) return null;

    return jwtDecode(jwt);
  }

  saveJWT(jwt: string) {
    const decoded: Object = jwtDecode(jwt);
    const expiration = decoded.exp;

    localStorage.setItem(JWT_KEY, jwt);
    localStorage.setItem(JWT_EXPIRATION_TIME_KEY, expiration.toString(10));
  }

  saveJWTFromResponse(res: Response) {
    window.res = res;
    const authorization: ?string = res.headers.get('Authorization');
    if (authorization && authorization.startsWith('Bearer ')) {
      this.saveJWT(authorization.substring(7));
    }
  }

  clearJWT() {
    localStorage.removeItem(JWT_KEY)
    localStorage.removeItem(JWT_EXPIRATION_TIME_KEY);
  }

  async checkRefreshJWT(): Promise<void> {
    const hasJWT = this.loadJWT() != null;
    const jwtExpiration = this.loadJWTExpiration();

    if (hasJWT && jwtExpiration) {
      const nowSeconds = Math.floor(new Date().getTime() / 1000);
      const jwtExpired = nowSeconds > jwtExpiration;

      if (jwtExpired) {
        await this.refreshJWT();
      }
    }
  }

  async post(path: string, body?: FormData | string | null = null, headers: Object = {}): Promise<any> {
    if (headers['Authorization'] === undefined) {
      await this.checkRefreshJWT();
      const jwt = this.loadJWT();
      if (jwt) {
        headers['Authorization'] = `Bearer ${jwt}`;
      }
    }

    const res = await super.post(path, body, headers);
    this.saveJWTFromResponse(res);
    return res;
  }

  async get(path: string, params: Object = {}, headers: Object = {}): Promise<any> {
    if (headers['Authorization'] === undefined) {
      await this.checkRefreshJWT();
      const jwt = this.loadJWT();
      if (jwt) {
        headers['Authorization'] = `Bearer ${jwt}`;
      }
    }

    const res = await super.get(path, params, headers);
    this.saveJWTFromResponse(res);
    return res;
  }
}
