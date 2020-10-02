import _ from "lodash";
import axios from "axios";
import TokenStorage from "./tokens";
import Config from "../secrets";
import { keysToCamel, keysToCamelWithWarnings } from "@/json-tools";
import { ExportParams } from "@/store/typed-actions";

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
    constructor(message) {
        super(message);
        this.name = "TokenError";
    }
}

export class AuthenticationRequiredError extends TokenError {
    constructor() {
        super("authentication required");
        this.name = "AuthenticationRequiredError";
    }
}

export class MissingTokenError extends TokenError {
    constructor() {
        super("missing token");
        this.name = "MissingTokenError";
    }
}

export type OnNoAuth<T> = () => Promise<T>;

export const OnNoReject = () => Promise.reject(new MissingTokenError());

export class LoginPayload {
    constructor(public readonly email: string, public readonly password: string) {}
}

export class LoginResponse {
    constructor(public readonly token: string | null) {}
}

export interface ExportStatus {
    id: number;
    token: string;
    createdAt: number;
    completedAt: number | null;
    kind: string;
    statusUrl: string;
    downloadUrl: string | null;
    progress: number;
    args: any;
}

export interface UserExports {
    exports: ExportStatus[];
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

export interface StationRegion {
    name: string;
    shape: [number, number][][];
}

export interface StationLocation {
    readonly precise: [number, number] | null;
    readonly regions: StationRegion[] | null;
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
    updatedAt: number;
    battery: number | null;
    location: StationLocation | null;
    placeNameOther: string | null;
    placeNameNative: string | null;
    recordingStartedAt: Date | null;
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

export enum Auth {
    None,
    Required,
    Optional,
}

export interface InvokeParams {
    auth: Auth;
    method: string;
    url: string;
    data?: any;
    contentType?: string;
    refreshed?: boolean | null;
    blob?: boolean | null;
}

class FKApi {
    private readonly baseUrl: string = Config.API_HOST;
    private readonly token: TokenStorage = new TokenStorage();
    private refreshing: Promise<any> | null = null;

    authenticated() {
        return this.token.authenticated();
    }

    private makeParams(params: InvokeParams): any {
        const headers = {
            "Content-Type": "application/json",
        };
        if (params.auth === Auth.Optional) {
            if (this.token.authenticated()) {
                const token = this.token.getHeader();
                headers["Authorization"] = token;
            }
        }
        if (params.auth === Auth.Required) {
            if (!this.token.authenticated()) {
                throw new AuthenticationRequiredError();
            }
            const token = this.token.getHeader();
            headers["Authorization"] = token;
        }
        return {
            method: params.method,
            url: params.url,
            headers: headers,
            data: params.data,
            responseType: params.blob === true ? "blob" : "json",
        };
    }

    private invoke(params: InvokeParams): Promise<any> {
        return axios(this.makeParams(params)).then(
            (response) => this.handle(params, response),
            (error) => {
                const response = error.response;

                if (!response) {
                    console.log(`api: error: ${error}`);
                    return error;
                }

                if (response.status === 401) {
                    if (params.refreshed !== true) {
                        // NOTE I'd like a better way to test for this.
                        if (response.data && response.data.message && response.data.message.indexOf("token") >= 0) {
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

    private afterTokenRefresh(): Promise<any> {
        if (this.refreshing !== null) {
            console.log("api: already refreshing");
            return this.refreshing;
        }

        const parsed = this.parseToken(this.token.getHeader());
        if (!parsed) {
            console.log("api: refresh skipped, invalid token");
            return this.logout().then(() => Promise.reject(new AuthenticationRequiredError()));
        }

        console.log("api: refreshing");

        const requestBody = {
            refreshToken: parsed.refresh_token, // eslint-disable-line
        };

        this.refreshing = axios({
            method: "POST",
            url: this.baseUrl + "/refresh",
            data: requestBody,
        })
            .then(
                (response) => this.handleLogin(response),
                (error) => this.logout().then(() => Promise.reject(error))
            )
            .finally(() => {
                this.refreshing = null;
            });

        return this.refreshing;
    }

    private refreshExpiredToken(original: InvokeParams): Promise<any> {
        return this.afterTokenRefresh().then(() => {
            console.log("api: retry original");
            return this.invoke(_.extend({ refreshed: true }, original));
        });
    }

    private handle(params: InvokeParams, response) {
        if (response.status == 200) {
            if (params.blob === true) {
                return response.data;
            }
            return keysToCamel(response.data);
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
            auth: Auth.None,
            method: "POST",
            url: this.baseUrl + "/users",
            data: user,
        });
    }

    resendCreateAccount(userId) {
        return this.invoke({
            auth: Auth.None,
            method: "POST",
            url: this.baseUrl + "/users/" + userId + "/validate-email",
        });
    }

    updatePassword(data) {
        return this.invoke({
            auth: Auth.Required,
            method: "PATCH",
            url: this.baseUrl + "/users/" + data.userId + "/password",
            data: { newPassword: data.newPassword, oldPassword: data.oldPassword },
        });
    }

    sendResetPasswordEmail(email) {
        return this.invoke({
            auth: Auth.None,
            method: "POST",
            url: this.baseUrl + "/user/recovery/lookup",
            data: { email: email },
        });
    }

    resetPassword(data) {
        return this.invoke({
            auth: Auth.None,
            method: "POST",
            url: this.baseUrl + "/user/recovery",
            data: { password: data.password, token: data.token },
        });
    }

    getStation(id): Promise<Station> {
        return this.invoke({
            auth: Auth.Required,
            method: "GET",
            url: this.baseUrl + "/stations/" + id,
        });
    }

    getUserStations(onNoAuth: OnNoAuth<StationsResponse>): Promise<StationsResponse> {
        if (!this.token.authenticated()) {
            return onNoAuth();
        }
        return this.invoke({
            auth: Auth.Required,
            method: "GET",
            url: this.baseUrl + "/user/stations",
        });
    }

    getCurrentUser(): Promise<CurrentUser> {
        return this.invoke({
            auth: Auth.Required,
            method: "GET",
            url: this.baseUrl + "/user",
        });
    }

    getUsersByProject(projectId): Promise<ProjectUsers> {
        return this.invoke({
            auth: Auth.Required,
            method: "GET",
            url: this.baseUrl + "/users/project/" + projectId,
        });
    }

    sendInvite(data) {
        return this.invoke({
            auth: Auth.Required,
            method: "POST",
            url: this.baseUrl + "/projects/" + data.projectId + "/invite",
            data: { email: data.email, role: data.role },
        });
    }

    getInvitesByToken(inviteToken) {
        return this.invoke({
            auth: Auth.Required,
            method: "GET",
            url: this.baseUrl + "/projects/invites/" + inviteToken,
        });
    }

    getInvitesByUser() {
        return this.invoke({
            auth: Auth.Required,
            method: "GET",
            url: this.baseUrl + "/projects/invites/pending",
        });
    }

    acceptProjectInvite(payload: { projectId: number }) {
        return this.invoke({
            auth: Auth.Required,
            method: "POST",
            url: this.baseUrl + "/projects/" + payload.projectId + "/invites/accept",
        });
    }

    declineProjectInvite(payload: { projectId: number }) {
        return this.invoke({
            auth: Auth.Required,
            method: "POST",
            url: this.baseUrl + "/projects/" + payload.projectId + "/invites/reject",
        });
    }

    acceptInvite(payload: { id: number; token: string }) {
        return this.invoke({
            auth: Auth.Required,
            method: "POST",
            url: this.baseUrl + "/projects/invites/" + payload.id + "/accept?token=" + payload.token,
        });
    }

    declineInvite(payload: { id: number; token: string }) {
        return this.invoke({
            auth: Auth.Required,
            method: "POST",
            url: this.baseUrl + "/projects/invites/" + payload.id + "/reject?token=" + payload.token,
            data: {
                token: payload.token,
            },
        });
    }

    getStationsByProject(projectId) {
        return this.invoke({
            auth: Auth.Required,
            method: "GET",
            url: this.baseUrl + "/projects/" + projectId + "/stations",
        });
    }

    addStationToProject(data) {
        return this.invoke({
            auth: Auth.Required,
            method: "POST",
            url: this.baseUrl + "/projects/" + data.projectId + "/stations/" + data.stationId,
        });
    }

    removeStationFromProject(data) {
        return this.invoke({
            auth: Auth.Required,
            method: "DELETE",
            url: this.baseUrl + "/projects/" + data.projectId + "/stations/" + data.stationId,
        });
    }

    removeUserFromProject(data) {
        return this.invoke({
            auth: Auth.Required,
            method: "DELETE",
            url: this.baseUrl + "/projects/" + data.projectId + "/members",
            data: { email: data.email },
        });
    }

    uploadUserImage(data) {
        return this.invoke({
            auth: Auth.Required,
            method: "POST",
            url: this.baseUrl + "/user/media",
            data: data.file,
        });
    }

    updateUser(data) {
        return this.invoke({
            auth: Auth.Required,
            method: "PATCH",
            url: this.baseUrl + "/users/" + data.id,
            data: data,
        });
    }

    getUserProjects(onNoAuth: OnNoAuth<ProjectsResponse>): Promise<ProjectsResponse> {
        if (!this.token.authenticated()) {
            return onNoAuth();
        }
        return this.invoke({
            auth: Auth.Required,
            method: "GET",
            url: this.baseUrl + "/user/projects",
        });
    }

    getPublicProjects(): Promise<ProjectsResponse> {
        return this.invoke({
            auth: Auth.Optional,
            method: "GET",
            url: this.baseUrl + "/projects",
        });
    }

    getProject(id): Promise<Project> {
        return this.invoke({
            auth: Auth.Optional,
            method: "GET",
            url: this.baseUrl + "/projects/" + id,
        });
    }

    getProjectActivity(id): Promise<ProjectActivityResponse> {
        return this.invoke({
            auth: Auth.Optional,
            method: "GET",
            url: this.baseUrl + "/projects/" + id + "/activity",
        });
    }

    addDefaultProject() {
        return this.invoke({
            auth: Auth.Required,
            method: "POST",
            url: this.baseUrl + "/projects",
            data: {
                name: "Default FieldKit Project",
                description: "Any FieldKit stations you add, start life here.",
                slug: "default-proj-" + Date.now(),
            },
        });
    }

    addProject(data) {
        return this.invoke({
            auth: Auth.Required,
            method: "POST",
            url: this.baseUrl + "/projects",
            data: data,
        });
    }

    updateProject(data) {
        return this.invoke({
            auth: Auth.Required,
            method: "PATCH",
            url: this.baseUrl + "/projects/" + data.id,
            data: data,
        });
    }

    uploadProjectImage(payload) {
        return this.invoke({
            auth: Auth.Required,
            method: "POST",
            url: this.baseUrl + "/projects/" + payload.id + "/media",
            contentType: payload.type,
            data: payload.file,
        });
    }

    deleteProject(data) {
        return this.invoke({
            auth: Auth.Required,
            method: "DELETE",
            url: this.baseUrl + "/projects/" + data.projectId,
        });
    }

    getProjectFollows(projectId): Promise<ProjectFollowers> {
        return this.invoke({
            auth: Auth.Optional,
            method: "GET",
            url: this.baseUrl + "/projects/" + projectId + "/followers",
        });
    }

    followProject(projectId) {
        return this.invoke({
            auth: Auth.Required,
            method: "POST",
            url: this.baseUrl + "/projects/" + projectId + "/follow",
        });
    }

    unfollowProject(projectId) {
        return this.invoke({
            auth: Auth.Required,
            method: "POST",
            url: this.baseUrl + "/projects/" + projectId + "/unfollow",
        });
    }

    deleteFieldNote(data) {
        return this.invoke({
            auth: Auth.Required,
            method: "DELETE",
            url: this.baseUrl + "/stations/" + data.stationId + "/field-notes/" + data.fieldNoteId,
        });
    }

    getModulesMeta() {
        return this.invoke({
            auth: Auth.Required,
            method: "GET",
            url: this.baseUrl + "/modules/meta",
        });
    }

    getStationDataSummaryByDeviceId(deviceId, start, end) {
        if (!start) {
            start = new Date("1/1/2019").getTime();
        }
        if (!end) {
            end = new Date().getTime();
        }

        return this.invoke({
            auth: Auth.Required,
            method: "GET",
            url: this.baseUrl + "/data/devices/" + deviceId + "/summary/json?start=" + start + "&end=" + end,
        });
    }

    getStationDataByDeviceId(deviceId) {
        return this.invoke({
            auth: Auth.Required,
            method: "GET",
            url: this.baseUrl + "/data/devices/" + deviceId + "/data",
        });
    }

    getJSONDataByDeviceId(deviceId, page, pageSize) {
        return this.invoke({
            auth: Auth.Required,
            method: "GET",
            url: this.baseUrl + "/data/devices/" + deviceId + "/data/json?page=" + page + "&pageSize=" + pageSize,
        });
    }

    getFieldNotes(stationId) {
        return this.invoke({
            auth: Auth.Required,
            method: "GET",
            url: this.baseUrl + "/stations/" + stationId + "/field-notes",
        });
    }

    addProjectUpdate(data) {
        return this.invoke({
            auth: Auth.Required,
            method: "POST",
            url: this.baseUrl + "/projects/" + data.projectId + "/updates",
            data: data,
        });
    }

    updateProjectUpdate(data) {
        return this.invoke({
            auth: Auth.Required,
            method: "POST",
            url: this.baseUrl + "/projects/" + data.projectId + "/updates/" + data.updateId,
            data: data,
        });
    }

    deleteProjectUpdate(data) {
        return this.invoke({
            auth: Auth.Required,
            method: "DELETE",
            url: this.baseUrl + "/projects/" + data.projectId + "/updates/" + data.updateId,
        });
    }

    getAllSensors() {
        return this.invoke({
            auth: Auth.Required,
            method: "GET",
            url: this.baseUrl + "/sensors",
        });
    }

    public exportData(queryParams: URLSearchParams, params: ExportParams): Promise<ExportStatus> {
        console.log("api:exporting", queryParams, params);
        const getUrl = () => {
            if (params.csv) return "/export/csv";
            if (params.jsonLines) return "/export/json-lines";
            throw new Error("unexecpted export params");
        };
        return this.invoke({
            auth: Auth.Required,
            method: "POST",
            url: this.baseUrl + getUrl() + "?" + queryParams.toString(),
        });
    }

    public exportStatus(url: string): Promise<ExportStatus> {
        return this.invoke({
            auth: Auth.Required,
            method: "GET",
            url: this.baseUrl + url,
        });
    }

    public getUserExports(): Promise<UserExports> {
        return this.invoke({
            auth: Auth.Required,
            method: "GET",
            url: this.baseUrl + "/export",
        });
    }

    public sensorData(params: URLSearchParams): Promise<any> {
        return this.invoke({
            auth: Auth.Required,
            method: "GET",
            url: this.baseUrl + "/sensors/data?" + params.toString(),
        });
    }

    public tailSensorData(params: URLSearchParams): Promise<any> {
        return this.invoke({
            auth: Auth.Required,
            method: "GET",
            url: this.baseUrl + "/sensors/data?" + params.toString(),
        });
    }

    public getQuickSensors(stations: number[]) {
        const qp = new URLSearchParams();
        qp.append("stations", stations.join(","));
        return this.invoke({
            auth: Auth.Required,
            method: "GET",
            url: this.baseUrl + "/sensors/data?" + qp.toString(),
        });
    }

    public adminDeleteUser(payload) {
        return this.invoke({
            auth: Auth.Required,
            method: "DELETE",
            url: this.baseUrl + "/admin/user",
            data: payload,
        });
    }

    public getStationNotes(stationId: number): Promise<any> {
        return this.invoke({
            auth: Auth.Required,
            method: "GET",
            url: this.baseUrl + "/stations/" + stationId + "/notes",
        });
    }

    public patchStationNotes(stationId: number, payload: any): Promise<any> {
        return this.invoke({
            auth: Auth.Required,
            method: "PATCH",
            url: this.baseUrl + "/stations/" + stationId + "/notes",
            data: { notes: payload },
        });
    }

    public uploadStationMedia(stationId: number, key: string, file: any): Promise<any> {
        const qp = new URLSearchParams();
        qp.append("key", key);
        return this.invoke({
            auth: Auth.Required,
            method: "POST",
            url: this.baseUrl + "/stations/" + stationId + "/media" + "?" + qp.toString(),
            data: file,
        });
    }

    public getAllStations(page: number, pageSize: number): Promise<PageOfStations> {
        const qp = new URLSearchParams();
        qp.append("page", page.toString());
        qp.append("pageSize", pageSize.toString());
        return this.invoke({
            auth: Auth.Required,
            method: "GET",
            url: this.baseUrl + "/admin/stations?" + qp.toString(),
        });
    }

    public deleteStation(stationId: number): Promise<any> {
        return this.invoke({
            auth: Auth.Required,
            method: "DELETE",
            url: this.baseUrl + "/admin/stations/" + stationId,
        });
    }

    public loadMedia(url: string, params: { size: number } | null = null): Promise<any> {
        const getUrl = () => {
            if (params) {
                const qp = new URLSearchParams();
                qp.append("size", params.size.toString());
                return this.baseUrl + url + "?" + qp.toString();
            }
            return this.baseUrl + url;
        };
        return this.invoke({
            auth: Auth.Required,
            method: "GET",
            blob: true,
            url: getUrl(),
        }).then((blob) => {
            const reader = new FileReader();
            return new Promise((resolve, reject) => {
                reader.onloadend = () => resolve(reader.result as string);
                reader.readAsDataURL(blob);
            });
        });
    }

    public getStatus(): Promise<PortalDeployStatus> {
        return this.invoke({
            auth: Auth.None,
            method: "GET",
            url: this.baseUrl + "/status",
        });
    }
}

export interface PortalDeployStatus {
    serverName: string;
    name: string;
    tag: string;
    git: { hash: string };
}

export interface EssentialStation {
    id: number;
    name: string;
    deviceId: string;
    owner: { id: number; name: string };
}

export interface PageOfStations {
    stations: EssentialStation[];
    total: number;
}

export default FKApi;
