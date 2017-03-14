// @flow

import 'whatwg-fetch'
import {BaseError} from '../common/util'
import log from 'loglevel';

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

  async post(path: string, body?: FormData | string): Promise<Response> {
    const url = new URL(path, this.baseUrl);

    log.info('POST', path, body, body);

    let res;
    try {
      res = await fetch(url.toString(), {
        method: 'POST',
        credentials: 'include',
        body
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

  async get(path: string, params?: Object): Promise<Response> {
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
        credentials: 'include'
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

  async postForm(path: string, body?: Object): Promise<string> {
    const data = new FormData();

    if (body) {
      for (const key in body) {
        data.append(key, body[key]);
      }
    }

    const res = await this.post(path, data);
    return res.text();
  }

  async postJSON(path: string, body?: Object): Promise<any> {
    const res = await this.post(path, JSON.stringify(body));
    if (res.status === 204) {
      return null
    }
    return res.json();
  }

  async getText(path: string, params?: Object): Promise<string> {
    const res = await this.get(path, params);
    return res.text();
  }

  async getJSON(path: string, params?: Object): Promise<any> {
    const res = await this.get(path, params);
    return res.json();
  }
}

const SIGNED_IN_KEY = 'signedIn';
export class AuthAPIClient extends APIClient {
  unauthorizedHandler: () => void;

  constructor(baseUrl: string, unauthorizedHandler: () => void) {
    super(baseUrl);

    this.unauthorizedHandler = unauthorizedHandler;
  }

  onAuthError(e: AuthenticationError) {
    if (this.signedIn()) {
      this.unauthorizedHandler();
      this.onSignout();
    }
  }

  onSignin() {
    localStorage.setItem(SIGNED_IN_KEY, 'loggedIn');
  }

  onSignout() {
    localStorage.removeItem(SIGNED_IN_KEY);
  }

  signedIn(): boolean {
    return localStorage.getItem(SIGNED_IN_KEY) != null;
  }

  async post(path: string, body?: FormData | string): Promise<any> {
    try {
      return await super.post(path, body);
    } catch (e) {
      if (e instanceof AuthenticationError) {
        this.onAuthError(e);
      }

      throw e;
    }
  }

  async get(path: string, params?: Object): Promise<any> {
    try {
      return await super.get(path, params);
    } catch (e) {
      if (e instanceof AuthenticationError) {
        this.onAuthError(e);
      }

      throw e;
    }
  }
}
