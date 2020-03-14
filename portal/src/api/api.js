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
                password: password,
            },
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
                Authorization: token,
            },
        }).then(this._handleResponse.bind(this));
    }

    getStations() {
        const token = this.token.getToken();
        return axios({
            method: "GET",
            url: this.baseUrl + "/stations",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        }).then(this._handleResponse.bind(this));
    }

    getCurrentUser() {
        const token = this.token.getToken();
        return axios({
            method: "GET",
            url: this.baseUrl + "/user",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        }).then(this._handleResponse.bind(this));
    }

    getUsersByProject(projectId) {
        const token = this.token.getToken();
        return axios({
            method: "GET",
            url: this.baseUrl + "/users/project/" + projectId,
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        }).then(this._handleResponse.bind(this));
    }

    sendInvite(data) {
        // temp! until there is a real email invitation in place,
        // only allow this to be tested with @conservify.org emails
        if (data.email.indexOf("@conservify.org") == -1) {
            return new Promise(resolve => {
                resolve("pseudo-invited!");
            });
        }
        const token = this.token.getToken();
        return axios({
            method: "POST",
            url: this.baseUrl + "/projects/" + data.projectId + "/invite",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
            data: { email: data.email },
        }).then(this._handleResponse.bind(this));
    }

    getStationsByProject(projectId) {
        const token = this.token.getToken();
        return axios({
            method: "GET",
            url: this.baseUrl + "/projects/" + projectId + "/stations",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        }).then(this._handleResponse.bind(this));
    }

    addStationToProject(data) {
        const token = this.token.getToken();
        return axios({
            method: "POST",
            url: this.baseUrl + "/projects/" + data.projectId + "/stations/" + data.stationId,
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        }).then(this._handleResponse.bind(this));
    }

    removeStationFromProject(data) {
        const token = this.token.getToken();
        return axios({
            method: "DELETE",
            url: this.baseUrl + "/projects/" + data.projectId + "/stations/" + data.stationId,
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        }).then(this._handleResponse.bind(this));
    }

    removeUserFromProject(data) {
        const token = this.token.getToken();
        return axios({
            method: "DELETE",
            url: this.baseUrl + "/projects/" + data.projectId + "/members/" + data.userId,
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
            data: { email: data.email },
        }).then(this._handleResponse.bind(this));
    }

    uploadUserImage(data) {
        const token = this.token.getToken();
        return axios({
            method: "POST",
            url: this.baseUrl + "/user/media",
            headers: {
                "Content-Type": data.type,
                Authorization: token,
            },
            data: data.image,
        }).then(this._handleResponse.bind(this));
    }

    updateUser(data) {
        const token = this.token.getToken();
        return axios({
            method: "PATCH",
            url: this.baseUrl + "/users/" + data.id,
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
            data: data,
        }).then(this._handleResponse.bind(this));
    }

    getProjects() {
        const token = this.token.getToken();
        return axios({
            method: "GET",
            url: this.baseUrl + "/user/projects",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        }).then(this._handleResponse.bind(this));
    }

    getProject(id) {
        const token = this.token.getToken();
        return axios({
            method: "GET",
            url: this.baseUrl + "/projects/" + id,
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        }).then(this._handleResponse.bind(this));
    }

    addDefaultProject() {
        const token = this.token.getToken();
        return axios({
            method: "POST",
            url: this.baseUrl + "/projects",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
            data: {
                name: "Default FieldKit Project",
                description: "Any FieldKit stations you add, start life here.",
                slug: "default-proj-" + Date.now(),
            },
        }).then(this._handleResponse.bind(this));
    }

    addProject(data) {
        const token = this.token.getToken();
        return axios({
            method: "POST",
            url: this.baseUrl + "/projects",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
            data: data,
        }).then(this._handleResponse.bind(this));
    }

    updateProject(data) {
        const token = this.token.getToken();
        return axios({
            method: "PATCH",
            url: this.baseUrl + "/projects/" + data.id,
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
            data: data,
        }).then(this._handleResponse.bind(this));
    }

    uploadProjectImage(data) {
        const token = this.token.getToken();
        return axios({
            method: "POST",
            url: this.baseUrl + "/projects/" + data.id + "/media",
            headers: {
                "Content-Type": data.type,
                Authorization: token,
            },
            data: data.image,
        }).then(this._handleResponse.bind(this));
    }

    deleteProject(data) {
        const token = this.token.getToken();
        return axios({
            method: "DELETE",
            url: this.baseUrl + "/projects/" + data.projectId,
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        }).then(this._handleResponse.bind(this));
    }

    deleteFieldNote(data) {
        const token = this.token.getToken();
        return axios({
            method: "DELETE",
            url: this.baseUrl + "/stations/" + data.stationId + "/field-notes/" + data.fieldNoteId,
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        }).then(this._handleResponse.bind(this));
    }

    getModulesMeta() {
        const token = this.token.getToken();
        return axios({
            method: "GET",
            url: this.baseUrl + "/modules/meta",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        }).then(this._handleResponse.bind(this));
    }

    getStationDataSummaryByDeviceId(deviceId, start, end) {
        if (!start) {
            start = new Date("1/1/2019").getTime();
        }
        if (!end) {
            end = new Date().getTime();
        }

        const token = this.token.getToken();
        return axios({
            method: "GET",
            url: this.baseUrl + "/data/devices/" + deviceId + "/summary/json?start=" + start + "&end=" + end,
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        }).then(this._handleResponse.bind(this));
    }

    getStationDataByDeviceId(deviceId) {
        const token = this.token.getToken();
        return axios({
            method: "GET",
            url: this.baseUrl + "/data/devices/" + deviceId + "/data",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        }).then(this._handleResponse.bind(this));
    }

    getJSONDataByDeviceId(deviceId, page, pageSize) {
        const token = this.token.getToken();
        return axios({
            method: "GET",
            url: this.baseUrl + "/data/devices/" + deviceId + "/data/json?page=" + page + "&pageSize=" + pageSize,
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        }).then(this._handleResponse.bind(this));
    }

    getFieldNotes(station) {
        const token = this.token.getToken();
        return axios({
            method: "GET",
            url: this.baseUrl + "/stations/" + station.id + "/field-notes",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        }).then(this._handleResponse.bind(this));
    }

    _handleResponse(response) {
        if (response.status == 200) {
            return response.data;
        } else if (response.status == 204) {
            return "success";
        } else {
            throw new Error("Api failed");
        }
    }
}

export default FKApi;
