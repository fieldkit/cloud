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
