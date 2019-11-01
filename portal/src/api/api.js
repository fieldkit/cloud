import axios from "axios";
import TokenStorage from "./tokens";
import { API_HOST } from "../secrets";

class FKApi {
    constructor() {
        this.baseUrl = API_HOST;
        this.token = new TokenStorage();
    }

    authenticated() {
        return this.token.authenticated();
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
        }).then(this._handleLoginResponse.bind(this));
    }

    logout() {
        this.token.clear();
    }

    _handleLoginResponse(response) {
        try {
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

    getStation(id) {
        const token = this.token.getToken();
        return axios({
            method: "GET",
            url: this.baseUrl + "/stations/@/" + id,
            headers: {
                "Content-Type": "application/json",
                Authorization: token
            }
        }).then(this._handleResponse.bind(this));
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
        }).then(this._handleResponse.bind(this));
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
        }).then(this._handleResponse.bind(this));
    }

    getUsers() {
        const token = this.token.getToken();
        return axios({
            method: "GET",
            url: this.baseUrl + "/users",
            headers: {
                "Content-Type": "application/json",
                Authorization: token
            }
        }).then(this._handleResponse.bind(this));
    }

    uploadUserImage(data) {
        const token = this.token.getToken();
        return axios({
            method: "POST",
            url: this.baseUrl + "/user/media",
            headers: {
                "Content-Type": data.type,
                Authorization: token
            },
            data: data.image
        }).then(this._handleResponse.bind(this));
    }

    updateUser(data) {
        const token = this.token.getToken();
        return axios({
            method: "PATCH",
            url: this.baseUrl + "/users/" + data.id,
            headers: {
                "Content-Type": "application/json",
                Authorization: token
            },
            data: data
        }).then(this._handleResponse.bind(this));
    }

    getProjects() {
        const token = this.token.getToken();
        return axios({
            method: "GET",
            url: this.baseUrl + "/user/projects",
            headers: {
                "Content-Type": "application/json",
                Authorization: token
            }
        }).then(this._handleResponse.bind(this));
    }

    getProject(id) {
        const token = this.token.getToken();
        return axios({
            method: "GET",
            url: this.baseUrl + "/projects/" + id,
            headers: {
                "Content-Type": "application/json",
                Authorization: token
            }
        }).then(this._handleResponse.bind(this));
    }

    addDefaultProject() {
        const token = this.token.getToken();
        return axios({
            method: "POST",
            url: this.baseUrl + "/projects",
            headers: {
                "Content-Type": "application/json",
                Authorization: token
            },
            data: {
                name: "Default FieldKit Project",
                description: "Any FieldKit stations you add, start life here.",
                slug: "default-proj-" + Date.now()
            }
        }).then(this._handleResponse.bind(this));
    }

    addProject(data) {
        const token = this.token.getToken();
        return axios({
            method: "POST",
            url: this.baseUrl + "/projects",
            headers: {
                "Content-Type": "application/json",
                Authorization: token
            },
            data: data
        }).then(this._handleResponse.bind(this));
    }

    updateProject(data) {
        const token = this.token.getToken();
        return axios({
            method: "PATCH",
            url: this.baseUrl + "/projects/" + data.id,
            headers: {
                "Content-Type": "application/json",
                Authorization: token
            },
            data: data
        }).then(this._handleResponse.bind(this));
    }

    uploadProjectImage(data) {
        const token = this.token.getToken();
        return axios({
            method: "POST",
            url: this.baseUrl + "/projects/" + data.id + "/media",
            headers: {
                "Content-Type": data.type,
                Authorization: token
            },
            data: data.image
        }).then(this._handleResponse.bind(this));
    }

    getStationDataSummaryByDeviceId(deviceId) {
        const token = this.token.getToken();
        return axios({
            method: "GET",
            url: this.baseUrl + "/data/devices/" + deviceId + "/summary",
            headers: {
                "Content-Type": "application/json",
                Authorization: token
            }
        }).then(this._handleResponse.bind(this));
    }

    getStationDataByDeviceId(deviceId) {
        const token = this.token.getToken();
        return axios({
            method: "GET",
            url: this.baseUrl + "/data/devices/" + deviceId + "/data",
            headers: {
                "Content-Type": "application/json",
                Authorization: token
            }
        }).then(this._handleResponse.bind(this));
    }

    getJSONDataByDeviceId(deviceId, page, pageSize) {
        const token = this.token.getToken();
        return axios({
            method: "GET",
            url:
                this.baseUrl +
                "/data/devices/" +
                deviceId +
                "/data/json?page=" +
                page +
                "&pageSize=" +
                pageSize,
            headers: {
                "Content-Type": "application/json",
                Authorization: token
            }
        }).then(this._handleResponse.bind(this));
    }

    _handleResponse(response) {
        if (response.status == 200) {
            return response.data;
        } else {
            throw new Error("Api failed");
        }
    }
}

export default FKApi;
