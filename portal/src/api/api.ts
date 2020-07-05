import _ from "lodash";
import axios from "axios";
import TokenStorage from "./tokens";
import Config from "../secrets";
import { keysToCamel, keysToCamelWithWarnings } from "@/json-tools";

export class ApiError extends Error {
    constructor(message) {
        super(message);
        this.name = "ApiError";
    }
}

export class ApiUnexpectedStatus extends ApiError {
    constructor(public readonly status: number) {
        super("unexpected status");
        this.name = "ApiUnexpectedStatus ";
    }
}

export class TokenError extends ApiError {
    authenticated: boolean;

    constructor(message) {
        super(message);
        this.name = "TokenError";
    }
}

export class LoginPayload {
    constructor(public readonly email: string, public readonly password: string) {}
}

export class LoginResponse {
    constructor(public readonly token: string | null) {}
}

export interface ProjectRef {
    id: number;
    name: string;
}

export interface StationRef {
    id: number;
    name: string;
}

export interface UploaderRef {
    id: number;
    name: string;
}

export interface UploadedData {
    id: number;
    records: number;
}

export interface UploadedMeta {
    uploader: UploaderRef;
    data: UploadedData;
    errors: boolean;
}

export interface Activity {
    id: number;
    key: string;
    project: ProjectRef;
    station: StationRef;
    createdAt: number;
    type: string;
    meta: UploadedMeta;
}

export interface ProjectActivityResponse {
    activities: Activity[];
    total: number;
    page: number;
}

export interface Avatar {
    url: string;
}

export interface Follower {
    id: number;
    name: string;
    avatar: Avatar;
}

export interface ProjectFollowers {
    followers: Follower[];
    total: number;
    page: number;
}

export interface SimpleUser {
    id: number;
    email: string;
    name: string;
    bio: string;
    mediaUrl: string;
    mediaContentType: string;
}

export interface ProjectUser {
    membership: string;
    role: string;
    user: SimpleUser;
}

export interface ProjectUsers {
    users: ProjectUser[];
}

export class CurrentUser {
    id: number;
    email: string;
    name: string;
    bio: string;
    mediaUrl: string;
}

export interface Project {
    description: string;
    goal: string;
    id: number;
    location: string;
    name: string;
    private: boolean;
    readOnly: boolean;
    slug: string;
    tags: string;
    mediaContentType: string;
    mediaUrl: string;
    startTime?: Date;
    endTime?: Date;
    numberOfFollowers: number;
}

export interface Owner {
    id: number;
    name: string;
}

export interface Upload {
    id: number;
    time: number;
    uploadId: string;
    size: number;
    url: string;
    type: string;
    blocks: number[];
}

export interface Photos {
    small: string;
}

export interface SensorReading {
    last: number;
    time: number;
}

export interface ModuleSensor {
    name: string;
    unitOfMeasure: string;
    key: string;
    ranges: null;
    reading: SensorReading | null;
}

export interface StationModule {
    id: number;
    name: string;
    hardwareId: string;
    position: number;
    internal: boolean;
    flags: number;
    sensors: ModuleSensor[];
}

export interface StationConfiguration {
    id: number;
    time: number;
    provisionId: number;
    modules: StationModule[];
}

export interface Configurations {
    all: StationConfiguration[];
}

export interface HasLocation {
    readonly latitude: number | null;
    readonly longitude: number | null;
}

export interface Station {
    id: number;
    name: string;
    owner: Owner;
    deviceId: string;
    uploads: Upload[];
    images: any[];
    photos: Photos;
    readOnly: boolean;
    configurations: Configurations;
    updated: number;
    location: HasLocation | null;
    placeNameOther: string | null;
    placeNameNative: string | null;
}

export interface ProjectsResponse {
    projects: Project[];
}

export interface StationsResponse {
    stations: Station[];
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

    private invoke(params): Promise<any> {
        return axios(params).then(
            (response) => this.handle(response),
            (error) => {
                const response = error.response;

                if (response.status === 401) {
                    if (params.refreshed !== true) {
                        // NOTE I'd like a better way to test for this.
                        if (response.data && response.data.detail && response.data.detail.indexOf("expired") >= 0) {
                            console.log("api: token expired");
                            return this.refreshExpiredToken(params);
                        }
                    }

                    // Token is super bad, so no use to us if we can't refresh.
                    this.token.clear();

                    console.log("api: refresh failed");
                    return Promise.reject(new TokenError("unauthorized"));
                }

                console.log("api: error", error.response.status, error.response.data);
                return Promise.reject(new ApiUnexpectedStatus(error.response.status));
            }
        );
    }

    private parseToken(token) {
        try {
            const encoded = token.split(".")[1];
            const decoded = Buffer.from(encoded, "base64").toString();
            return JSON.parse(decoded);
        } catch (e) {
            console.log("api: error parsing token", e, "token", token);
            return null;
        }
    }

    private refreshExpiredToken(original) {
        const parsed = this.parseToken(this.token.getHeader());
        const requestBody = {
            refresh_token: parsed.refresh_token, // eslint-disable-line
        };

        console.log("api: refreshing");
        return axios({
            method: "POST",
            url: this.baseUrl + "/refresh",
            data: requestBody,
        }).then(
            (response) => {
                return this.handleLogin(response).then(() => {
                    console.log("api: retry original");
                    original.headers.Authorization = this.token.getHeader();
                    return this.invoke(_.extend({ refreshed: true }, original));
                });
            },
            (error) => this.logout().then(() => Promise.reject(error))
        );
    }

    private handle(response) {
        if (response.status == 200) {
            // eslint-disable-next-line no-constant-condition
            if (false) {
                return keysToCamelWithWarnings(response.data);
            } else {
                return keysToCamel(response.data);
            }
        } else if (response.status == 204) {
            return true;
        } else {
            throw new ApiError("api: error: unknown");
        }
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
        }).then((response) => this.handleLogin(response));
    }

    private handleLogin(response): Promise<string> {
        try {
            if (response.status == 204) {
                this.token.setToken(response.headers.authorization);
                return Promise.resolve(response.headers.authorization);
            } else {
                throw new ApiError("login failed");
            }
        } catch (err) {
            console.log("api: login error:", err, err.stack);
            throw new ApiError("login failed");
        }
    }

    logout() {
        this.token.clear();
        return Promise.resolve();
    }

    register(user) {
        return this.invoke({
            method: "POST",
            url: this.baseUrl + "/users",
            data: user,
        });
    }

    resendCreateAccount(userId) {
        return this.invoke({
            method: "POST",
            url: this.baseUrl + "/users/" + userId + "/validate-email",
        });
    }

    updatePassword(data) {
        const token = this.token.getHeader();
        return this.invoke({
            method: "PATCH",
            url: this.baseUrl + "/users/" + data.userId + "/password",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
            data: { newPassword: data.newPassword, oldPassword: data.oldPassword },
        });
    }

    sendResetPasswordEmail(email) {
        const token = this.token.getHeader();
        return this.invoke({
            method: "POST",
            url: this.baseUrl + "/user/recovery/lookup",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
            data: { email: email },
        });
    }

    resetPassword(data) {
        return this.invoke({
            method: "POST",
            url: this.baseUrl + "/user/recovery",
            data: { password: data.password, token: data.token },
        });
    }

    getStationFromVuex(id): Promise<Station> {
        const token = this.token.getHeader();
        return this.invoke({
            method: "GET",
            url: this.baseUrl + "/stations/" + id,
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        });
    }

    getStation(id): Promise<Station> {
        const token = this.token.getHeader();
        return this.invoke({
            method: "GET",
            url: this.baseUrl + "/stations/@/" + id,
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        });
    }

    getStations(): Promise<StationsResponse> {
        const token = this.token.getHeader();
        return this.invoke({
            method: "GET",
            url: this.baseUrl + "/stations",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        });
    }

    getCurrentUser(): Promise<CurrentUser> {
        const token = this.token.getHeader();
        return this.invoke({
            method: "GET",
            url: this.baseUrl + "/user",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        });
    }

    getUsersByProject(projectId): Promise<ProjectUsers> {
        const token = this.token.getHeader();
        return this.invoke({
            method: "GET",
            url: this.baseUrl + "/users/project/" + projectId,
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        });
    }

    sendInvite(data) {
        const token = this.token.getHeader();
        return this.invoke({
            method: "POST",
            url: this.baseUrl + "/projects/" + data.projectId + "/invite",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
            data: { email: data.email, role: data.role },
        });
    }

    getInvitesByToken(inviteToken) {
        const token = this.token.getHeader();
        return this.invoke({
            method: "GET",
            url: this.baseUrl + "/projects/invites/" + inviteToken,
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        });
    }

    getInvitesByUser() {
        const token = this.token.getHeader();
        return this.invoke({
            method: "GET",
            url: this.baseUrl + "/projects/invites/pending",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        });
    }

    acceptInvite(inviteId) {
        const token = this.token.getHeader();
        return this.invoke({
            method: "POST",
            url: this.baseUrl + "/projects/invites/" + inviteId + "/accept",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        });
    }

    declineInvite(inviteId) {
        const token = this.token.getHeader();
        return this.invoke({
            method: "POST",
            url: this.baseUrl + "/projects/invites/" + inviteId + "/reject",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        });
    }

    getStationsByProject(projectId) {
        const token = this.token.getHeader();
        return this.invoke({
            method: "GET",
            url: this.baseUrl + "/projects/" + projectId + "/stations",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        });
    }

    addStationToProject(data) {
        const token = this.token.getHeader();
        return this.invoke({
            method: "POST",
            url: this.baseUrl + "/projects/" + data.projectId + "/stations/" + data.stationId,
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        });
    }

    removeStationFromProject(data) {
        const token = this.token.getHeader();
        return this.invoke({
            method: "DELETE",
            url: this.baseUrl + "/projects/" + data.projectId + "/stations/" + data.stationId,
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        });
    }

    removeUserFromProject(data) {
        const token = this.token.getHeader();
        return this.invoke({
            method: "DELETE",
            url: this.baseUrl + "/projects/" + data.projectId + "/members",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
            data: { email: data.email },
        });
    }

    uploadUserImage(data) {
        const token = this.token.getHeader();
        return this.invoke({
            method: "POST",
            url: this.baseUrl + "/user/media",
            headers: {
                "Content-Type": data.type,
                Authorization: token,
            },
            data: data.image,
        });
    }

    updateUser(data) {
        const token = this.token.getHeader();
        return this.invoke({
            method: "PATCH",
            url: this.baseUrl + "/users/" + data.id,
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
            data: data,
        });
    }

    getUserProjects(): Promise<ProjectsResponse> {
        const token = this.token.getHeader();
        return this.invoke({
            method: "GET",
            url: this.baseUrl + "/user/projects",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        });
    }

    getPublicProjects(): Promise<ProjectsResponse> {
        const token = this.token.getHeader();
        return this.invoke({
            method: "GET",
            url: this.baseUrl + "/projects",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        });
    }

    getProject(id): Promise<Project> {
        const token = this.token.getHeader();
        return this.invoke({
            method: "GET",
            url: this.baseUrl + "/projects/" + id,
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        });
    }

    getProjectActivity(id): Promise<ProjectActivityResponse> {
        const token = this.token.getHeader();
        return this.invoke({
            method: "GET",
            url: this.baseUrl + "/projects/" + id + "/activity",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        });
    }

    addDefaultProject() {
        const token = this.token.getHeader();
        return this.invoke({
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
        });
    }

    addProject(data) {
        const token = this.token.getHeader();
        return this.invoke({
            method: "POST",
            url: this.baseUrl + "/projects",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
            data: data,
        });
    }

    updateProject(data) {
        const token = this.token.getHeader();
        return this.invoke({
            method: "PATCH",
            url: this.baseUrl + "/projects/" + data.id,
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
            data: data,
        });
    }

    uploadProjectImage(data) {
        const token = this.token.getHeader();
        return this.invoke({
            method: "POST",
            url: this.baseUrl + "/projects/" + data.id + "/media",
            headers: {
                "Content-Type": data.type,
                Authorization: token,
            },
            data: data.image,
        });
    }

    deleteProject(data) {
        const token = this.token.getHeader();
        return this.invoke({
            method: "DELETE",
            url: this.baseUrl + "/projects/" + data.projectId,
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        });
    }

    getProjectFollows(projectId): Promise<ProjectFollowers> {
        return this.invoke({
            method: "GET",
            url: this.baseUrl + "/projects/" + projectId + "/followers",
            headers: {
                "Content-Type": "application/json",
            },
        });
    }

    followProject(projectId) {
        const token = this.token.getHeader();
        return this.invoke({
            method: "POST",
            url: this.baseUrl + "/projects/" + projectId + "/follow",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        });
    }

    unfollowProject(projectId) {
        const token = this.token.getHeader();
        return this.invoke({
            method: "POST",
            url: this.baseUrl + "/projects/" + projectId + "/unfollow",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        });
    }

    deleteFieldNote(data) {
        const token = this.token.getHeader();
        return this.invoke({
            method: "DELETE",
            url: this.baseUrl + "/stations/" + data.stationId + "/field-notes/" + data.fieldNoteId,
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        });
    }

    getModulesMeta() {
        const token = this.token.getHeader();
        return this.invoke({
            method: "GET",
            url: this.baseUrl + "/modules/meta",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        });
    }

    getStationDataSummaryByDeviceId(deviceId, start, end) {
        if (!start) {
            start = new Date("1/1/2019").getTime();
        }
        if (!end) {
            end = new Date().getTime();
        }

        const token = this.token.getHeader();
        return this.invoke({
            method: "GET",
            url: this.baseUrl + "/data/devices/" + deviceId + "/summary/json?start=" + start + "&end=" + end,
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        });
    }

    getStationDataByDeviceId(deviceId) {
        const token = this.token.getHeader();
        return this.invoke({
            method: "GET",
            url: this.baseUrl + "/data/devices/" + deviceId + "/data",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        });
    }

    getJSONDataByDeviceId(deviceId, page, pageSize) {
        const token = this.token.getHeader();
        return this.invoke({
            method: "GET",
            url: this.baseUrl + "/data/devices/" + deviceId + "/data/json?page=" + page + "&pageSize=" + pageSize,
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        });
    }

    getFieldNotes(stationId) {
        const token = this.token.getHeader();
        return this.invoke({
            method: "GET",
            url: this.baseUrl + "/stations/" + stationId + "/field-notes",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        });
    }

    addProjectUpdate(data) {
        const token = this.token.getHeader();
        return this.invoke({
            method: "POST",
            url: this.baseUrl + "/projects/" + data.projectId + "/updates",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
            data: data,
        });
    }

    updateProjectUpdate(data) {
        const token = this.token.getHeader();
        return this.invoke({
            method: "POST",
            url: this.baseUrl + "/projects/" + data.projectId + "/updates/" + data.updateId,
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
            data: data,
        });
    }

    deleteProjectUpdate(data) {
        const token = this.token.getHeader();
        return this.invoke({
            method: "DELETE",
            url: this.baseUrl + "/projects/" + data.projectId + "/updates/" + data.updateId,
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        });
    }

    getAllSensors() {
        const token = this.token.getHeader();
        return this.invoke({
            url: this.baseUrl + "/sensors",
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        });
    }

    sensorData(params: any): Promise<any> {
        const queryParams = new URLSearchParams();
        queryParams.append("start", params.when.start.getTime());
        queryParams.append("end", params.when.end.getTime());
        queryParams.append("stations", params.stations.join(","));
        queryParams.append("sensors", params.sensors.join(","));
        const token = this.token.getHeader();
        return this.invoke({
            url: this.baseUrl + "/sensors/data?" + queryParams.toString(),
            headers: {
                "Content-Type": "application/json",
                Authorization: token,
            },
        });
    }
}

export default FKApi;
