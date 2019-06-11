import { FkPromisedApi } from '../api/calls';

import TokenStorage  from './tokens';

export default class UserSession {
    constructor() {
        this.tokens = new TokenStorage();
    }

    async authenticated() {
        return this.tokens.authenticated();
    }

    async login(email, password) {
        try {
            await FkPromisedApi.login(email, password);

            const user = await FkPromisedApi.getAuthenticatedUser();

            console.log(user);
        }
        finally {
        }

        return this.user;
    }

    async logout() {
        this.tokens.clear();
    }
};
