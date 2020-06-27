import axios from "axios";
import TokenStorage from "./tokens";
import Config from "../secrets";
import { keysToCamel } from "@/json-tools";

export class LoginPayload {
    constructor(public readonly email: string, public readonly password: string) {}
}

export class LoginResponse {
    constructor(public readonly token: string | null) {}
}

export class CurrentUser {
    id: number;
    email: string;
    name: string;
    bio: string;
    mediaUrl: string;
}

// Intentionally keeping this synchronous since it'll get used in
// VueJS stuff quite often to make URLs that don't require custom
// headers for authentication.
export function makeAuthenticatedApiUrl(url) {
    const tokens = new TokenStorage();
    const token = tokens.getToken();
    return Config.API_HOST + url + "?token=" + token;
}

class FKApi {
    private readonly baseUrl: string = Config.API_HOST;
    private readonly token: TokenStorage = new TokenStorage();

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
        }).then(response => this.handleLogin(response));
    }

    logout() {
        this.token.clear();
    }

    private handleLogin(response): string {
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

    register(user) {
        return axios({
            method: "POST",
            url: this.baseUrl + "/users",
            data: user,
        }).then(response => this.handle(response));
    }

    resendCreateAccount(userId) {
        return axios({
            method: "POST",
            url: this.baseUrl + "/users/" + userId + "/validate-email",
        }).then(response => this.handle(response));
    }

    updatePassword(data) {
        const token = this.token.getHeader();
        return axios({
            method: "PATCH",
            url: this.baseUrl + "/users/" + data.userId + "/password",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
            data: { newPassword: data.newPassword, oldPassword: data.oldPassword },
        }).then(response => this.handle(response));
    }

    sendResetPasswordEmail(email) {
        const token = this.token.getHeader();
        return axios({
            method: "POST",
            url: this.baseUrl + "/user/recovery/lookup",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
            data: { email: email },
        }).then(response => this.handle(response));
    }

    resetPassword(data) {
        return axios({
            method: "POST",
            url: this.baseUrl + "/user/recovery",
            data: { password: data.password, token: data.token },
        }).then(response => this.handle(response));
    }

    getStation(id) {
        const token = this.token.getHeader();
        return axios({
            method: "GET",
            url: this.baseUrl + "/stations/@/" + id,
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        }).then(response => this.handle(response));
    }

    getStations() {
        const token = this.token.getHeader();
        return axios({
            method: "GET",
            url: this.baseUrl + "/stations",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        }).then(response => this.handle(response));
    }

    getCurrentUser() {
        const token = this.token.getHeader();
        return axios({
            method: "GET",
            url: this.baseUrl + "/user",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        }).then(response => this.handle(response));
    }

    getUsersByProject(projectId) {
        const token = this.token.getHeader();
        return axios({
            method: "GET",
            url: this.baseUrl + "/users/project/" + projectId,
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        }).then(response => this.handle(response));
    }

    sendInvite(data) {
        const token = this.token.getHeader();
        return axios({
            method: "POST",
            url: this.baseUrl + "/projects/" + data.projectId + "/invite",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
            data: { email: data.email, role: data.role },
        }).then(response => this.handle(response));
    }

    getInvitesByToken(inviteToken) {
        const token = this.token.getHeader();
        return axios({
            method: "GET",
            url: this.baseUrl + "/projects/invites/" + inviteToken,
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        }).then(response => this.handle(response));
    }

    getInvitesByUser() {
        const token = this.token.getHeader();
        return axios({
            method: "GET",
            url: this.baseUrl + "/projects/invites/pending",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        }).then(response => this.handle(response));
    }

    acceptInvite(inviteId) {
        const token = this.token.getHeader();
        return axios({
            method: "POST",
            url: this.baseUrl + "/projects/invites/" + inviteId + "/accept",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        }).then(response => this.handle(response));
    }

    declineInvite(inviteId) {
        const token = this.token.getHeader();
        return axios({
            method: "POST",
            url: this.baseUrl + "/projects/invites/" + inviteId + "/reject",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        }).then(response => this.handle(response));
    }

    getStationsByProject(projectId) {
        const token = this.token.getHeader();
        return axios({
            method: "GET",
            url: this.baseUrl + "/projects/" + projectId + "/stations",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        }).then(response => this.handle(response));
    }

    addStationToProject(data) {
        const token = this.token.getHeader();
        return axios({
            method: "POST",
            url: this.baseUrl + "/projects/" + data.projectId + "/stations/" + data.stationId,
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        }).then(response => this.handle(response));
    }

    removeStationFromProject(data) {
        const token = this.token.getHeader();
        return axios({
            method: "DELETE",
            url: this.baseUrl + "/projects/" + data.projectId + "/stations/" + data.stationId,
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        }).then(response => this.handle(response));
    }

    removeUserFromProject(data) {
        const token = this.token.getHeader();
        return axios({
            method: "DELETE",
            url: this.baseUrl + "/projects/" + data.projectId + "/members",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
            data: { email: data.email },
        }).then(response => this.handle(response));
    }

    uploadUserImage(data) {
        const token = this.token.getHeader();
        return axios({
            method: "POST",
            url: this.baseUrl + "/user/media",
            headers: {
                "Content-Type": data.type,
                Authorization: token,
            },
            data: data.image,
        }).then(response => this.handle(response));
    }

    updateUser(data) {
        const token = this.token.getHeader();
        return axios({
            method: "PATCH",
            url: this.baseUrl + "/users/" + data.id,
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
            data: data,
        }).then(response => this.handle(response));
    }

    getUserProjects() {
        const token = this.token.getHeader();
        return axios({
            method: "GET",
            url: this.baseUrl + "/user/projects",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        }).then(response => this.handle(response));
    }

    getPublicProjects() {
        const token = this.token.getHeader();
        return axios({
            method: "GET",
            url: this.baseUrl + "/projects",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        }).then(response => {
            if (response.status == 200) {
                return response.data.projects.filter(p => {
                    return !p.private;
                });
            } else {
                throw new Error("Api failed");
            }
        });
    }

    getProject(id) {
        const token = this.token.getHeader();
        return axios({
            method: "GET",
            url: this.baseUrl + "/projects/" + id,
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        }).then(response => this.handle(response));
    }

    getProjectActivity(id) {
        const token = this.token.getHeader();
        return axios({
            method: "GET",
            url: this.baseUrl + "/projects/" + id + "/activity",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        }).then(response => this.handle(response));
    }

    addDefaultProject() {
        const token = this.token.getHeader();
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
        }).then(response => this.handle(response));
    }

    addProject(data) {
        const token = this.token.getHeader();
        return axios({
            method: "POST",
            url: this.baseUrl + "/projects",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
            data: data,
        }).then(response => this.handle(response));
    }

    updateProject(data) {
        const token = this.token.getHeader();
        return axios({
            method: "PATCH",
            url: this.baseUrl + "/projects/" + data.id,
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
            data: data,
        }).then(response => this.handle(response));
    }

    uploadProjectImage(data) {
        const token = this.token.getHeader();
        return axios({
            method: "POST",
            url: this.baseUrl + "/projects/" + data.id + "/media",
            headers: {
                "Content-Type": data.type,
                Authorization: token,
            },
            data: data.image,
        }).then(response => this.handle(response));
    }

    deleteProject(data) {
        const token = this.token.getHeader();
        return axios({
            method: "DELETE",
            url: this.baseUrl + "/projects/" + data.projectId,
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        }).then(response => this.handle(response));
    }

    getProjectFollows(projectId) {
        return axios({
            method: "GET",
            url: this.baseUrl + "/projects/" + projectId + "/followers",
            headers: {
                "Content-Type": "application/json",
            },
        }).then(response => this.handle(response));
    }

    followProject(projectId) {
        const token = this.token.getHeader();
        return axios({
            method: "POST",
            url: this.baseUrl + "/projects/" + projectId + "/follow",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        }).then(response => this.handle(response));
    }

    unfollowProject(projectId) {
        const token = this.token.getHeader();
        return axios({
            method: "POST",
            url: this.baseUrl + "/projects/" + projectId + "/unfollow",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        }).then(response => this.handle(response));
    }

    deleteFieldNote(data) {
        const token = this.token.getHeader();
        return axios({
            method: "DELETE",
            url: this.baseUrl + "/stations/" + data.stationId + "/field-notes/" + data.fieldNoteId,
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        }).then(response => this.handle(response));
    }

    getModulesMeta() {
        const token = this.token.getHeader();
        return axios({
            method: "GET",
            url: this.baseUrl + "/modules/meta",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        }).then(response => this.handle(response));
    }

    getStationDataSummaryByDeviceId(deviceId, start, end) {
        if (!start) {
            start = new Date("1/1/2019").getTime();
        }
        if (!end) {
            end = new Date().getTime();
        }

        const token = this.token.getHeader();
        return axios({
            method: "GET",
            url: this.baseUrl + "/data/devices/" + deviceId + "/summary/json?start=" + start + "&end=" + end,
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        }).then(response => this.handle(response));
    }

    getStationDataByDeviceId(deviceId) {
        const token = this.token.getHeader();
        return axios({
            method: "GET",
            url: this.baseUrl + "/data/devices/" + deviceId + "/data",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        }).then(response => this.handle(response));
    }

    getJSONDataByDeviceId(deviceId, page, pageSize) {
        const token = this.token.getHeader();
        return axios({
            method: "GET",
            url: this.baseUrl + "/data/devices/" + deviceId + "/data/json?page=" + page + "&pageSize=" + pageSize,
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        }).then(response => this.handle(response));
    }

    getFieldNotes(stationId) {
        const token = this.token.getHeader();
        return axios({
            method: "GET",
            url: this.baseUrl + "/stations/" + stationId + "/field-notes",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        }).then(response => this.handle(response));
    }

    getPlaceName(longLat) {
        return axios({
            method: "GET",
            url:
                "https://api.mapbox.com/geocoding/v5/mapbox.places/" +
                longLat +
                ".json?types=place&access_token=" +
                Config.MAPBOX_ACCESS_TOKEN,
            headers: {
                "Content-Type": "application/json",
            },
        }).then(response => this.handle(response));
    }

    getNativeLand(location) {
        return axios({
            method: "GET",
            url: "https://native-land.ca/api/index.php?maps=territories&position=" + location,
            headers: {
                "Content-Type": "application/json",
            },
        }).then(response => this.handle(response));
    }

    addProjectUpdate(data) {
        const token = this.token.getHeader();
        return axios({
            method: "POST",
            url: this.baseUrl + "/projects/" + data.projectId + "/updates",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
            data: data,
        }).then(response => this.handle(response));
    }

    updateProjectUpdate(data) {
        const token = this.token.getHeader();
        return axios({
            method: "POST",
            url: this.baseUrl + "/projects/" + data.projectId + "/updates/" + data.updateId,
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
            data: data,
        }).then(response => this.handle(response));
    }

    deleteProjectUpdate(data) {
        const token = this.token.getHeader();
        return axios({
            method: "DELETE",
            url: this.baseUrl + "/projects/" + data.projectId + "/updates/" + data.updateId,
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        }).then(response => this.handle(response));
    }

    private handle(response) {
        if (response.status == 200) {
            return keysToCamel(response.data);
        } else if (response.status == 204) {
            return true;
        } else {
            throw new Error("Api failed");
        }
    }
}

export default FKApi;
