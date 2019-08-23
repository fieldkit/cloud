import axios from "axios";
import TokenStorage from "./tokens";
import { API_HOST } from "../secrets";

class FKApi {
    constructor() {
        this.baseUrl = API_HOST;
        this.token = new TokenStorage();
        this.handleResponse = this.handleResponse.bind(this);
    }

    login(email, password) {
        return axios({
            method: "POST",
            url: this.baseUrl + "/login",
            headers: { "Content-Type": "application/json" },
            data: {
                email: email,
                password: password
            }
        }).then(this.handleResponse);
    }

    handleResponse(response) {
        try {
            // console.log("login function response", response);
            if (response.status == 204) {
                this.token.setToken(response.headers.authorization);
                return response.headers.authorization;
            } else {
                throw new Error("Log In Failed");
            }
        } catch (err) {
            // console.log("Error:", err);
            throw new Error("Log In Failed");
        }
    }

    authenticated() {
        return this.token.authenticated();
    }

    logout() {
        this.token.clear();
    }

    getStations() {
        const token = this.token.getToken();
        return axios({
            method: "GET",
            url: this.baseUrl + "/stations",
            headers: {
                "Content-Type": "application/json",
                Authorization: token
            }
        }).then(handleResponse);

        function handleResponse(response) {
            if (response.status) {
                return response.data;
            } else {
                throw new Error("Get station failed");
            }
        }
    }

    getCurrentUser() {
        const token = this.token.getToken();
        return axios({
            method: "GET",
            url: this.baseUrl + "/user",
            headers: {
                "Content-Type": "application/json",
                Authorization: token
            }
        }).then(handleResponse);

        function handleResponse(response) {
            if (response.status == 200) {
                return response.data;
            } else {
                throw new Error("Get user failed");
            }
        }
    }
}

export default FKApi;
