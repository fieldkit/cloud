export default class TokenStorage {
    constructor() {
        this.token = this.getToken();
    }

    getToken() {
        return (this.token = JSON.parse(window.localStorage["fktoken"] || "null"));
    }

    authenticated() {
        console.log("authenticated token", this.token);
        return this.token != null;
    }

    setToken(token) {
        window.localStorage["fktoken"] = JSON.stringify(token);
        this.token = token;
    }

    clear() {
        this.token = null;
        delete window.localStorage["fktoken"];
    }
}
