import _ from "lodash";
import axios from "axios";

import Config from "../secrets";
import TokenStorage from "./tokens";
import Backoff from "backoff";

import { keysToCamel } from "@/json-tools";

import { ExportParams } from "@/store/typed-actions";
import { BoundingRectangle } from "@/store/map-types";

import { SensorInfoResponse } from "@/views/viz/api";

// Ew
import { NewComment } from "@/views/comments/model";
import { Comment } from "@/views/comments/model";
import { SensorsResponse, VizConfig } from "@/views/viz/api";
import {PortalStationFieldNotes} from '@/views/notes/model';

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
    owner: { id: number; name: string; email: string };
    uploads: { id: number }[];
}

export interface PageOfStations {
    stations: EssentialStation[];
    total: number;
}

export interface MentionableUser {
    id: number;
    name: string;
    photo: { url: string };
}

export class ApiError extends Error {
    constructor(message) {
        super(message);
        this.name = "ApiError";
    }

    public static isInstance(err: Error): boolean {
        return err.name === "ApiError";
    }
}

export class ApiUnexpectedStatus extends ApiError {
    constructor(public readonly status: number) {
        super("unexpected status");
        this.name = "ApiUnexpectedStatus ";
    }

    public static isInstance(err: Error): boolean {
        return err.name === "ApiUnexpectedStatus";
    }
}

export class TokenError extends ApiError {
    constructor(message) {
        super(message);
        this.name = "TokenError";
    }

    public static isInstance(err: Error): boolean {
        return err.name === "TokenError";
    }
}

export class AuthenticationRequiredError extends TokenError {
    constructor() {
        super("authentication required");
        this.name = "AuthenticationRequiredError";
    }

    public static isInstance(err: Error): boolean {
        return err.name === "AuthenticationRequiredError";
    }
}

export class MissingTokenError extends TokenError {
    constructor() {
        super("missing token");
        this.name = "MissingTokenError";
    }

    public static isInstance(err: Error): boolean {
        return err.name === "MissingTokenError";
    }
}

export class ForbiddenError extends ApiError {
    constructor(public readonly status: number) {
        super("403 status");
        this.name = "ForbiddenError";
    }

    public static isInstance(err: Error): boolean {
        return err.name === "ForbiddenError";
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
    args: Record<string, unknown>;
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
    photo: { url: string };
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
    tncDate: number;
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
    startTime?: Date;
    endTime?: Date;
    bounds?: BoundingRectangle;
    showStations: boolean;
    photo: string;
    following: {
        following: boolean;
        total: number;
    };
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
    fullKey: string;
    name: string;
    unitOfMeasure: string;
    key: string;
    ranges: null;
    reading: SensorReading | null;
    meta: {
        viz: VizConfig[];
    };
}

export interface StationModule {
    id: number;
    name: string;
    hardwareId: string;
    hardwareIdBase64: string;
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

export interface ProjectAttribute {
    attributeId: number;
    projectId: number;
    name: string;
    stringValue: string;
    priority: number;
}

export enum StationStatus {
    up = "up",
    down = "down",
}

export interface Station {
    id: number;
    name: string;
    owner: Owner;
    deviceId: string;
    uploads: Upload[];
    photos: Photos;
    readOnly: boolean;
    configurations: Configurations;
    updatedAt: number;
    lastReadingAt: number | null;
    battery: number | null;
    location: StationLocation | null;
    locationName: string;
    placeNameOther: string | null;
    placeNameNative: string | null;
    recordingStartedAt: Date | null;
    firmwareNumber: number | null;
    attributes: {
        attributes: ProjectAttribute[];
    };
    status: StationStatus;
}

export interface ProjectsResponse {
    projects: Project[];
}

export interface StationsResponse {
    stations: Station[];
}

export interface AssociatedStation {
    station: Station;
    manual?: { otherStationID: number; priority: number }[];
    location?: { stationID: number; distance: number }[];
    project?: { id: number };
    hidden: boolean;
}

export interface AssociatedStationsResponse {
    stations: AssociatedStation[];
}

export type SendFunction = (message: unknown) => Promise<void>;

export interface TailSensorDataRow {
    time: number;
    stationId: number;
    sensorId: number;
    moduleId: string;
    avg: number | undefined;
    min: number | undefined;
    max: number | undefined;
    last: number | undefined;
}

export interface TailSensorDataResponse {
    data: TailSensorDataRow[];
    stations: { [index: number]: { bucketSize: number } };
}

export type RecentlyAggregatedWindows = { [index: number]: TailSensorDataRow[] };

export type RecentlyAggregatedLast = { last: number | null };

export interface QueryRecentlyResponse {
    windows: RecentlyAggregatedWindows;
    stations: { [index: number]: RecentlyAggregatedLast };
}

// Intentionally keeping this synchronous since it'll get used in
// VueJS stuff quite often to make URLs that don't require custom
// headers for authentication.
export function makeAuthenticatedApiUrl(url) {
    const tokens = new TokenStorage();
    const token = tokens.getToken();
    return Config.baseUrl + url + "?token=" + token;
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
    data?: Record<string, unknown>;
    contentType?: string;
    refreshed?: boolean | null;
    blob?: boolean | null;
}

export interface SavedBookmark {
    url: string;
    bookmark: string;
}

export enum MapViewType {
    map = "map",
    list = "list",
}

class FKApi {
    private readonly baseUrl: string = Config.baseUrl;
    private readonly token: TokenStorage = new TokenStorage();
    private refreshing: Promise<any> | null = null;
    public allSensorsMemoized = _.memoize(() => this.getAllSensors());

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

                if (response.status === 403) {
                    return Promise.reject(new ForbiddenError(403));
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
            return this.logout(true).then(() => Promise.reject(new AuthenticationRequiredError()));
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
                () => this.logout(true).then(() => Promise.reject(new AuthenticationRequiredError()))
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

    public login(email: string, password: string): Promise<any> {
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

    public loginDiscourse(
        token: string | null,
        email: string | null,
        password: string | null,
        sso: string,
        sig: string
    ): Promise<{ token: string; location: string; header: string }> {
        const headers = {
            "Content-Type": "application/json",
        };
        if (token) {
            headers["Authorization"] = "Bearer " + token;
        }
        return axios({
            method: "POST",
            url: this.baseUrl + "/discourse/auth",
            headers: headers,
            data: {
                email: email,
                password: password,
                sso: sso,
                sig: sig,
            },
        })
            .then((response) => {
                return response.data as { token: string; location: string; header: string };
            })
            .then((response) => {
                this.token.setToken(response.header);
                return response;
            });
    }

    public loginOidc(
        token: string | null,
        params: {
            state: string;
            sessionState: string;
            code: string;
        }
    ): Promise<{ token: string; location: string; header: string }> {
        const headers = {
            "Content-Type": "application/json",
        };
        if (token) {
            headers["Authorization"] = "Bearer " + token;
        }
        const qp = new URLSearchParams();
        qp.append("state", params.state);
        qp.append("session_state", params.sessionState);
        qp.append("code", params.code);
        return axios({
            method: "POST",
            url: this.baseUrl + "/oidc/auth?" + qp.toString(),
            headers: headers,
        })
            .then((response) => {
                return response.data as { token: string; location: string; header: string };
            })
            .then((response) => {
                this.token.setToken(response.header);
                return response;
            });
    }

    public loginUrl(after: string | null): Promise<string> {
        const qp = new URLSearchParams();
        if (after) {
            qp.append("after", after);
        }
        return axios({
            method: "GET",
            url: this.baseUrl + "/oidc/url?" + qp.toString(),
        }).then((response) => {
            return response.data.location;
        });
    }

    public loginResume(token: string): Promise<any> {
        return axios({
            method: "POST",
            url: this.baseUrl + "/user/resume",
            headers: { "Content-Type": "application/json" },
            data: {
                token: token,
            },
        }).then((response) => this.handleLogin(response));
    }

    private handleLogin(response): Promise<string> {
        try {
            if (response.status == 204) {
                this.token.setToken(response.headers.authorization);
                if (this.wsBackoff != null) {
                    this.wsBackoff.reset();
                }
                return Promise.resolve(response.headers.authorization);
            } else {
                throw new ApiError("login failed");
            }
        } catch (err) {
            console.log("api: login error:", err, err.stack);
            throw new ApiError("login failed");
        }
    }

    public async logout(discardToken = false): Promise<void> {
        try {
            if (!this.token.authenticated()) {
                return Promise.resolve();
            }
            if (!discardToken) {
                const token = this.token.getHeader();
                const headers = { "Content-Type": "application/json", Authorization: token };
                await axios({
                    method: "POST",
                    url: this.baseUrl + "/logout",
                    headers: headers,
                });
            }
        } catch (err) {
            console.log("api: logout error:", err, err.stack);
        } finally {
            this.token.clear();
        }
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

    accept(userId) {
        return this.invoke({
            auth: Auth.Required,
            method: "PATCH",
            url: this.baseUrl + "/users/" + userId + "/accept-tnc",
            data: { accept: true },
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
            auth: Auth.Optional,
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

    getAssociatedStations(id: number): Promise<AssociatedStationsResponse> {
        return this.invoke({
            auth: Auth.Optional,
            method: "GET",
            url: this.baseUrl + `/stations/${id}/associated`,
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
            auth: Auth.Optional,
            method: "GET",
            url: this.baseUrl + "/users/project/" + projectId,
        });
    }

    editRole(data: { projectId: number; email: string; role: number }) {
        return this.invoke({
            auth: Auth.Required,
            method: "PATCH",
            url: this.baseUrl + "/projects/" + data.projectId + "/roles",
            data: { email: data.email, role: data.role },
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
            auth: Auth.Optional,
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

    public getAllSensorsMemoized(): () => Promise<SensorsResponse> {
        return this.allSensorsMemoized;
    }

    private getAllSensors(): Promise<SensorsResponse> {
        return this.invoke({
            auth: Auth.Optional,
            method: "GET",
            url: this.baseUrl + "/meta/sensors",
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

    public stationMeta(params: URLSearchParams): Promise<any> {
        return this.invoke({
            auth: Auth.Optional,
            method: "GET",
            url: this.baseUrl + "/meta/stations?" + params.toString(),
        });
    }

    public sensorData(params: URLSearchParams): Promise<any> {
        return this.invoke({
            auth: Auth.Optional,
            method: "GET",
            url: this.baseUrl + "/sensors/data?" + params.toString(),
        });
    }

    public tailSensorData(stations: number[]): Promise<TailSensorDataResponse> {
        const qp = new URLSearchParams();
        qp.append("stations", stations.join(","));
        return this.invoke({
            auth: Auth.Optional,
            method: "GET",
            url: this.baseUrl + "/sensors/data/tail?" + qp.toString(),
        });
    }

    public queryStationsRecently(stations: number[]): Promise<QueryRecentlyResponse> {
        if (stations.length == 0) {
            return Promise.resolve({
                windows: {},
                stations: {},
            });
        }
        const qp = new URLSearchParams();
        qp.append("stations", stations.join(","));
        return this.invoke({
            auth: Auth.Optional,
            method: "GET",
            url: this.baseUrl + "/sensors/data/recently?" + qp.toString(),
        });
    }

    public getQuickSensors(stations: number[]): Promise<SensorInfoResponse> {
        const qp = new URLSearchParams();
        qp.append("stations", stations.join(","));
        return this.invoke({
            auth: Auth.Optional,
            method: "GET",
            url: this.baseUrl + "/meta/stations?" + qp.toString(),
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

    public adminClearTermsAndConditions(payload) {
        return this.invoke({
            auth: Auth.Required,
            method: "DELETE",
            url: this.baseUrl + "/admin/user/tnc",
            data: payload,
        });
    }

    public getStationNotes(stationId: number): Promise<any> {
        return this.invoke({
            auth: Auth.Optional,
            method: "GET",
            url: this.baseUrl + "/stations/" + stationId + "/notes",
        });
    }

    public async getStationFieldNotes(stationId: number): Promise<any> {
        const response = await this.invoke({
            auth: Auth.Optional,
            method: "GET",
            url: this.baseUrl + "/station/" + stationId + "/station-notes",
        });

        const notes = response.notes.map((note) => {
            note.body = JSON.parse(note.body);
            return note;
        });

        return { notes: notes };
    }

    public addStationFieldNote(stationId: number, note: PortalStationFieldNotes): Promise<any> {
        note.body = JSON.stringify(note.body);
        return this.invoke({
            auth: Auth.Required,
            method: "POST",
            url: this.baseUrl + "/station/" + stationId + "/station-note",
            data: note,
        });
    }

    public patchStationFieldNote(stationId: number, note: PortalStationFieldNotes): Promise<any> {
        return this.invoke({
            auth: Auth.Required,
            method: "POST",
            url: this.baseUrl + "/station/" + stationId + "/station-note/" + note.id,
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

    setStationImage(stationId: number, photoId: number) {
        return this.invoke({
            auth: Auth.Required,
            method: "POST",
            url: this.baseUrl + "/stations/" + stationId + "/photo/" + photoId,
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

    public mentionables(query: string): Promise<{ users: MentionableUser[] }> {
        const qp = new URLSearchParams();
        qp.append("query", query);
        return this.invoke({
            auth: Auth.Required,
            method: "GET",
            url: this.baseUrl + "/mentionables?" + qp.toString(),
        });
    }

    public adminSearchUsers(query: string): Promise<{ users: SimpleUser[] }> {
        const qp = new URLSearchParams();
        qp.append("query", query);
        return this.invoke({
            auth: Auth.Required,
            method: "POST",
            url: this.baseUrl + "/admin/users/search?" + qp.toString(),
        });
    }

    public adminSearchStations(query: string): Promise<PageOfStations> {
        const qp = new URLSearchParams();
        qp.append("query", query);
        return this.invoke({
            auth: Auth.Required,
            method: "POST",
            url: this.baseUrl + "/admin/stations/search?" + qp.toString(),
        });
    }

    public adminTransferStation(stationId: number, ownerId: number): Promise<void> {
        return this.invoke({
            auth: Auth.Required,
            method: "POST",
            url: this.baseUrl + `/stations/${stationId}/transfer/${ownerId}`,
        });
    }

    public adminProcessStationData(stationId: number): Promise<void> {
        return this.invoke({
            auth: Auth.Required,
            method: "POST",
            url: this.baseUrl + `/data/stations/${stationId}/ingestions/process`,
        });
    }

    public adminProcessUpload(ingestionId: number): Promise<void> {
        const qp = new URLSearchParams();
        return this.invoke({
            auth: Auth.Required,
            method: "POST",
            url: this.baseUrl + `/data/ingestions/${ingestionId}/process?` + qp.toString(),
        });
    }

    // Think twice before you use this. Every pending ingestion_queue should have a que_job.
    protected adminProcessPending(): Promise<void> {
        return this.invoke({
            auth: Auth.Required,
            method: "POST",
            url: this.baseUrl + `/data/process`,
        });
    }

    public adminProcessStation(stationId: number, completely: boolean, skipManual: boolean): Promise<void> {
        const qp = new URLSearchParams();
        if (completely) {
            qp.append("completely", "true");
        }
        if (skipManual) {
            qp.append("skipManual", "true");
        }
        return this.invoke({
            auth: Auth.Required,
            method: "POST",
            url: this.baseUrl + `/data/stations/${stationId}/process?` + qp.toString(),
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
            auth: Auth.Optional,
            method: "GET",
            blob: true,
            url: getUrl(),
        }).then((blob) => {
            const reader = new FileReader();
            return new Promise((resolve) => {
                reader.onloadend = () => resolve(reader.result as string);
                reader.readAsDataURL(blob);
            });
        });
    }

    deleteMedia(mediaId: number) {
        return this.invoke({
            auth: Auth.Required,
            method: "DELETE",
            url: this.baseUrl + "/notes/media/" + mediaId,
        });
    }

    public getStatus(): Promise<PortalDeployStatus> {
        return this.invoke({
            auth: Auth.None,
            method: "GET",
            url: this.baseUrl + "/status",
        });
    }

    private parseBody(post: Comment): Comment {
        try {
            return _.extend(post, {
                body: JSON.parse(post.body),
                replies: this.parseBodies(post.replies),
            });
        } catch (error) {
            return _.extend(post, {
                body: {
                    type: "doc",
                    content: [{ type: "paragraph", content: [{ type: "text", text: post.body }] }],
                },
                replies: this.parseBodies(post.replies),
            });
        }
    }

    private parseBodies(posts: Comment[]): Comment[] {
        return posts.map((c) => {
            return this.parseBody(c);
        });
    }

    public async getComments(projectIDOrBookmark: number | string): Promise<{ posts: Comment[] }> {
        let apiURL;

        if (typeof projectIDOrBookmark === "number") {
            apiURL = this.baseUrl + "/discussion/projects/" + projectIDOrBookmark;
        } else {
            apiURL = this.baseUrl + "/discussion?bookmark=" + encodeURIComponent(JSON.stringify(projectIDOrBookmark));
        }

        const returned = await this.invoke({
            auth: Auth.Optional,
            method: "GET",
            url: apiURL,
        });

        const fixed = this.parseBodies(returned.posts);

        // console.log("comments", returned);

        return {
            posts: fixed,
        };
    }

    public async postComment(comment: NewComment): Promise<{ post: Comment }> {
        console.log("save-comment", comment);

        const returned = await this.invoke({
            auth: Auth.Required,
            method: "POST",
            url: this.baseUrl + "/discussion",
            data: {
                post: _.extend({}, comment, {
                    body: JSON.stringify(comment.body),
                }),
            },
        });

        // console.log("comments", returned);

        return {
            post: this.parseBody(returned.post),
        };
    }

    public async deleteComment(commentID: number): Promise<boolean> {
        return await this.invoke({
            auth: Auth.Required,
            method: "DELETE",
            url: this.baseUrl + "/discussion/" + commentID,
        });
    }

    public async editComment(commentID: number, body: Record<string, unknown>): Promise<boolean> {
        const returned = await this.invoke({
            auth: Auth.Required,
            method: "POST",
            url: this.baseUrl + "/discussion/" + commentID,
            data: { body: JSON.stringify(body) },
        });

        console.log("edit", returned);

        return returned;
    }

    public async seenNotifications(payload) {
        return this.invoke({
            auth: Auth.Required,
            method: "POST",
            url: this.baseUrl + "/notifications/seen",
            data: { ids: payload.ids },
        });
    }

    private socket: WebSocket | null = null;

    private async send(message: unknown) {
        if (!this.socket) {
            throw new Error("disconnected");
        }
        // TODO Does this need to queue before open?
        this.socket.send(JSON.stringify(message));
        return await Promise.resolve();
    }

    private async establish(callback: (message: unknown) => Promise<void>, status: (connected: boolean) => Promise<void>) {
        if (this.token.authenticated() && !this.socket) {
            const wsBase = this.baseUrl.replace("https", "wss").replace("http", "ws");
            this.socket = new WebSocket(wsBase + "/notifications");

            this.socket.addEventListener("open", () => {
                console.log("ws: connected");
                if (!this.socket) throw new Error("disconnected");
                const token = this.token.getHeader();
                this.socket.send(JSON.stringify({ token: token }));
            });

            this.socket.addEventListener("message", (event) => {
                const message = JSON.parse(event.data);
                if (message.error) {
                    console.log("ws: error", message.error);
                } else {
                    console.log("ws: message", message);
                    void callback(message);
                    this.wsBackoff.reset();
                }
            });

            this.socket.addEventListener("close", async () => {
                // console.log("ws: closed");
                void status(false);
                this.socket = null;

                this.wsBackoff.backoff();
            });
        }
    }

    private wsBackoff: Backoff | null = null;

    public async listenForNotifications(
        callback: (message: unknown) => Promise<void>,
        status: (connected: boolean) => Promise<void>
    ): Promise<SendFunction> {
        if (this.wsBackoff == null) {
            this.wsBackoff = Backoff.fibonacci({
                randomisationFactor: 0,
                initialDelay: 1000,
                maxDelay: 1000 * 60,
            });

            this.wsBackoff.on("ready", async (number: number, delay: number) => {
                // console.log("ws: ready", number, delay);

                await this.establish(callback, status);
            });

            this.wsBackoff.backoff();
        }

        return this.send.bind(this);
    }

    public async saveBookmark(bookmark: string): Promise<SavedBookmark> {
        const qp = new URLSearchParams();
        qp.append("bookmark", bookmark);
        return this.invoke({
            auth: Auth.Optional,
            method: "POST",
            url: this.baseUrl + `/bookmarks/save?${qp.toString()}`,
        });
    }

    public async resolveBookmark(token: string): Promise<SavedBookmark> {
        const qp = new URLSearchParams();
        qp.append("v", token);
        return this.invoke({
            auth: Auth.Optional,
            method: "GET",
            url: this.baseUrl + `/bookmarks/resolve?${qp.toString()}`,
        });
    }
}

export default FKApi;
