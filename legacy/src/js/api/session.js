import { FkPromisedApi } from '../api/calls';

import TokenStorage  from './tokens';

export default class UserSession {
    constructor() {
        this.tokens = new TokenStorage();
    }

    authenticated() {
        return this.tokens.authenticated();
    }

    async login(email, password) {
        try {
            await FkPromisedApi.login(email, password);

            return await FkPromisedApi.getAuthenticatedUser();
        }
        catch (err) {
            console.log(err);

            return null;
        }
    }

    async logout() {
        this.tokens.clear();
    }
};
