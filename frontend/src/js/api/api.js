
import 'whatwg-fetch';

import { BaseError } from '../common/errors';
import { API_HOST } from '../secrets';

import TokenStorage from './tokens';

class APIError extends BaseError {
}

class AuthenticationError extends APIError {
}

class FKApiClient {
    constructor(baseUrl) {
        this.baseUrl = baseUrl;
        this.tokens = new TokenStorage();
    }

    getHeaders() {
        if (!this.tokens.authenticated()) {
            return {
            };
        }

        const token = this.tokens.getToken();

        return {
            "Authorization": token,
        };
    }

    async get(path, params) {
        try {
            const url = new URL(path, this.baseUrl);
            if (params) {
                for (const key in params) {
                    url.searchParams.set(key, params[key]);
                }
            }
            let res;
            try {
                res = await fetch(url.toString(), {
                    method: 'GET',
                    headers: this.getHeaders()
                });
            } catch (e) {
                console.log('Threw while GETing', url.toString(), e);
                throw new APIError('HTTP error');
            }
            if (res.status === 404) {
                console.log('404', url.toString(), await res.text());
                throw new APIError('not found');
            }
            if (res.status === 401) {
                console.log('Authentication Error', url.toString(), await res.text());
                throw new AuthenticationError();
            } else if (!res.ok) {
                const err = await res.json();
                console.log('Unexpected Status', url.toString(), err);
                throw new APIError(err);
            }
            return res;
        } catch (e) {
            if (e instanceof AuthenticationError) {
                this.onAuthError(e);
            }
            throw e;
        }
    }

    async getJSON(path: string, params?: Object): Promise<any> {
        const res = await this.get(path, params);
        return res.json();
    }

    async post(path, body, contentType) {
        try {
            const url = new URL(path, this.baseUrl);
            let res;
            try {
                res = await fetch(url.toString(), {
                    method: 'POST',
                    headers: this.getHeaders(),
                    body
                });
            } catch (e) {
                console.log('Threw while POSTing', url.toString(), e);
                throw new APIError('HTTP error');
            }
            if (res.status === 401) {
                console.log('Bad auth while POSTing', url.toString(), await res.text());
                throw new AuthenticationError();
            } else if (!res.ok) {
                const err = await res.json();
                console.log('Non-OK response while POSTing', url.toString(), err);
                throw new APIError(err);
            }
            return res;
        } catch (e) {
            if (e instanceof AuthenticationError) {
                this.onAuthError(e);
            }
            throw e;
        }
    }

    async postJSON(path: string, params?: Object): Promise<any> {
        const res = await this.post(path, JSON.stringify(params), "application/json");

        let body = {};

        const contentLength = res.headers.get("Content-Length");
        if (contentLength) {
            body = res.json();
        }

        const token = res.headers.get("Authorization");
        if (token) {
            this.tokens.setToken(token);
            return {
                authorization: token,
                body
            };
        }

        return body;
    }

    async postForm(path, body) {
        const data = new FormData();
        if (body) {
            for (const key in body) {
                data.append(key, body[key]);
            }
        }
        const res = await this.post(path, data);
        return res.text();
    }

    onAuthError(error) {
        console.log("Authentication Error", error);
    }
};

export default new FKApiClient(API_HOST);
