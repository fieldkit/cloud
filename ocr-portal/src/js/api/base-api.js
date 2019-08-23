// @flow

import 'whatwg-fetch'
import { BaseError } from '../common/util'
import log from 'loglevel';
import jwtDecode from 'jwt-decode';

const JWT_KEY = 'jwt';
const JWT_EXPIRATION_TIME_KEY = 'jwt-expiration';

export type SupportedMethods = 'GET' | 'POST' | 'DELETE' | 'HEAD' | 'OPTIONS' | 'PUT' | 'PATCH';

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

    async exec(
        method: SupportedMethods,
        path: string, {params, body, headers = {} }: {
            params?: Object,
            body?: ?(Blob | FormData | URLSearchParams | string),
            headers: Object} = {}
    ): Promise<Response> {
        const url = new URL(path, this.baseUrl);

        if (params) {
            for (const key in params) {
                url.searchParams.set(key, params[key]);
            }
        }

        log.info(method, path, {
            params,
            body,
            headers
        });

        let res;
        try {
            res = await fetch(url.toString(), {
                method,
                params,
                body,
                headers
            });
        } catch (e) {
            log.error(e);
            log.error('Threw during', method, url.toString());
            throw new APIError('HTTP error', res);
        }

        if (!res.ok) {
            const body = await res.text();
            if (res.status === 401) {
                log.error('Bad auth during', method, url.toString(), body);
                throw new AuthenticationError(res, body);
            } else {
                log.error('Non-OK response during', method, url.toString(), body);
                throw new APIError('HTTP error', res, body);
            }
        }

        return res;
    }

    post(path: string, body?: FormData | string | null, headers: Object = {}): Promise<Response> {
        return this.exec('POST', path, {
            body,
            headers
        });
    }

    patch(path: string, body?: FormData | string | null, headers: Object = {}): Promise<Response> {
        return this.exec('PATCH', path, {
            body,
            headers
        });
    }

    put(path: string, body?: FormData | string | null, headers: Object = {}): Promise<Response> {
        return this.exec('PUT', path, {
            body,
            headers
        });
    }

    get(path: string, params?: Object, headers: Object = {}): Promise<Response> {
        return this.exec('GET', path, {
            params,
            headers
        });
    }

    del(path: string, body?: FormData | string | null, headers: Object = {}): Promise<Response> {
        return this.exec('DELETE', path, {
            body,
            headers
        });
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

    async patchJSON(path: string, body?: Object, headers: Object = {}): Promise<any> {
        const res = await this.patch(path, JSON.stringify(body), headers);
        if (res.status === 204) {
            return null
        }
        return res.json();
    }

    async putJSON(path: string, body?: Object, headers: Object = {}): Promise<any> {
        const res = await this.patch(path, JSON.stringify(body), headers);
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

    async delJSON(path: string, body?: Object, headers: Object = {}): Promise<any> {
        const res = await this.del(path, JSON.stringify(body), headers);
        if (res.status === 204) {
            return null
        }
        return res.json();
    }
}

export class JWTAPIClient extends APIClient {
    async refreshJWT(): Promise<void> {
        const jwt = this.loadDecodedJWT();
        if (!jwt) return;

        const refreshToken: string = jwt.refresh_token;
        try {
            await super.postJSON('/refresh', {
                refresh_token: refreshToken
            }, {
                Authorization: null
            });
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

        const nowSeconds = Math.floor(new Date().getTime() / 1000);
        console.log('Got JWT with expiration', expiration, 'now is', nowSeconds, decoded);

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
                console.log('JWT is expired, refreshing...', nowSeconds, jwtExpiration, this.loadDecodedJWT());
                await this.refreshJWT();
            }
        }
    }

    async exec(
        method: SupportedMethods,
        path: string, {params, body, headers = {} }: {
            params?: Object,
            body?: ?(Blob | FormData | URLSearchParams | string),
            headers: Object} = {}
    ): Promise<Response> {
        if (headers['Authorization'] === undefined) {
            await this.checkRefreshJWT();
            const jwt = this.loadJWT();
            if (jwt) {
                headers['Authorization'] = `Bearer ${jwt}`;
            }
        }
        const res = await super.exec(method, path, {
            params,
            body,
            headers
        });
        this.saveJWTFromResponse(res);
        return res;
    }
}
