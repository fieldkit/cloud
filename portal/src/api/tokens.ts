export default class TokenStorage {
    private token: string | null;

    constructor() {
        this.token = this.getToken();
    }

    public getToken() {
        this.token = this.sanitize(JSON.parse(window.localStorage["fktoken"] || "null"));
        return this.token;
    }

    public getHeader() {
        return "Bearer " + this.getToken();
    }

    public authenticated() {
        return this.token != null;
    }

    public setToken(token) {
        const sanitized = this.sanitize(token);
        window.localStorage["fktoken"] = JSON.stringify(sanitized);
        this.token = sanitized;
    }

    public clear() {
        this.token = null;
        delete window.localStorage["fktoken"];
    }

    private sanitize(token) {
        if (token) {
            return token.replace("Bearer ", "");
        }
        return null;
    }
}
