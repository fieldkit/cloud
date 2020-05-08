export default class TokenStorage {
    constructor() {
        this.token = this.getToken();
    }

    getToken() {
        this.token = this._sanitize(JSON.parse(window.localStorage["fktoken"] || "null"));
        return this.token;
    }

    getHeader() {
        return "Bearer " + this.getToken();
    }

    authenticated() {
        return this.token != null;
    }

    setToken(token) {
        const sanitized = this._sanitize(token);
        window.localStorage["fktoken"] = JSON.stringify(sanitized);
        this.token = sanitized;
    }

    clear() {
        this.token = null;
        delete window.localStorage["fktoken"];
    }

    _sanitize(token) {
        if (token) {
            return token.replace("Bearer ", "");
        }
        return null;
    }
}
