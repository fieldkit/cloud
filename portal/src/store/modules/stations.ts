import _ from "lodash";
import Vue from "vue";
import * as MutationTypes from "../mutations";
import * as ActionTypes from "../actions";
import { Location, BoundingRectangle, LngLat } from "../map-types";

import {
    FKApi,
    Services,
    OnNoReject,
    Station,
    StationModule,
    ModuleSensor,
    StationRegion,
    Project,
    ProjectUser,
    ProjectFollowers,
    Activity,
    Configurations,
    Photos,
} from "@/api";

export const HAVE_USER_STATIONS = "HAVE_USER_STATIONS";
export const HAVE_USER_PROJECTS = "HAVE_USER_PROJECTS";
export const HAVE_COMMUNITY_PROJECTS = "HAVE_COMMUNITY_PROJECTS";
export const PROJECT_USERS = "PROJECT_USERS";
export const PROJECT_FOLLOWS = "PROJECT_FOLLOWS";
export const PROJECT_STATIONS = "PROJECT_STATIONS";
export const PROJECT_ACTIVITY = "PROJECT_ACTIVITY";
export const STATION_UPDATE = "STATION_UPDATE";
export const PROJECT_LOADED = "PROJECT_LOADED";
export const PROJECT_UPDATE = "PROJECT_UPDATE";
export const PROJECT_DELETED = "PROJECT_DELETED";

export class StationsState {
    stations: { [index: number]: DisplayStation } = {};
    hasNoStations: boolean;
    projects: { [index: number]: Project } = {};
    projectUsers: { [index: number]: ProjectUser[] } = {};
    projectFollowers: { [index: number]: ProjectFollowers } = {};
    projectStations: { [index: number]: DisplayStation[] } = {};
    projectActivities: { [index: number]: Activity[] } = {};
    user: {
        stations: { [index: number]: DisplayStation };
        projects: { [index: number]: Project };
    } = {
        stations: {},
        projects: {},
    };
    community: {
        projects: { [index: number]: Project };
    } = {
        projects: {},
    };
    mapped: MappedStations | null = null;
}

export function whenWasStationUpdated(station: Station): Date {
    const uploads = station.uploads.map((u) => u.time);
    if (uploads.length > 0) {
        return new Date(uploads[0]);
    }
    return new Date(station.updatedAt);
}

export class DisplaySensor {
    name: string;
    unitOfMeasure: string;
    reading: number | null;

    constructor(sensor: ModuleSensor) {
        this.name = sensor.name;
        this.unitOfMeasure = sensor.unitOfMeasure;
        this.reading = sensor.reading?.last || null;
    }
}

export class DisplayModule {
    name: string;
    sensors: DisplaySensor[];

    constructor(module: StationModule) {
        this.name = module.name;
        this.sensors = module.sensors.map((s) => new DisplaySensor(s));
    }
}

export class DisplayStation {
    public readonly id: number;
    public readonly name: string;
    public readonly configurations: Configurations;
    public readonly updatedAt: Date;
    public readonly uploadedAt: Date | null = null;
    public readonly deployedAt: Date | null = null;
    public readonly totalReadings = 0;
    public readonly location: Location | null = null;
    public readonly photos: Photos;
    public readonly modules: DisplayModule[] = [];
    public readonly placeNameOther: string | null;
    public readonly placeNameNative: string | null;
    public readonly battery: number | null;
    public readonly regions: StationRegion[] | null;

    constructor(station: Station) {
        this.id = station.id;
        this.name = station.name;
        this.configurations = station.configurations;
        this.photos = station.photos;
        this.battery = station.battery;
        this.placeNameOther = station.placeNameOther;
        this.placeNameNative = station.placeNameNative;
        this.deployedAt = station.recordingStartedAt;
        if (!station.updatedAt) throw new Error(`station missing updatedAt`);
        this.updatedAt = new Date(station.updatedAt);
        this.uploadedAt = _.first(station.uploads.filter((u) => u.type == "data").map((u) => new Date(u.time))) || null;
        this.modules =
            _(station.configurations.all)
                .map((c) => c.modules.filter((m) => !m.internal).map((m) => new DisplayModule(m)))
                .head() || [];
        if (station.location) {
            if (station.location.precise) {
                this.location = new Location(station.location.precise[1], station.location.precise[0]);
            }
            this.regions = station.location.regions;
        }
    }
}

export class ProjectModule {
    constructor(public readonly name: string, public url: string | null) {}
}

export class MapFeature {
    public readonly type = "Feature";
    public readonly geometry: { type: string; coordinates: LngLat | LngLat[][] } | null = null;
    public readonly properties: { title: string; icon: string; id: number } | null = null;

    constructor(station: DisplayStation, type: string, coordinates: any, public readonly bounds: LngLat[]) {
        this.geometry = {
            type: type,
            coordinates: coordinates,
        };
        this.properties = {
            id: station.id,
            title: station.name,
            icon: "marker",
        };
    }

    public static makeFeatures(station: DisplayStation): MapFeature[] {
        if (station.location) {
            const precise = station.location.lngLat();
            return [new MapFeature(station, "Point", precise, [precise])];
        }
        if (!station.regions || station.regions.length === 0) {
            throw new Error("unmappable");
        }
        const region = station.regions[0];
        return region.shape.map((polygon) => new MapFeature(station, "Polygon", [polygon], polygon));
    }
}

const DefaultLocation: LngLat = [-118.0730372, 34.3318104]; // TwinPeaks
const DefaultMargin = 10000;

export class MappedStations {
    public static make(stations: DisplayStation[]): MappedStations {
        const located = stations.filter((station) => station.location != null);
        const features = _.flatten(located.map((ds) => MapFeature.makeFeatures(ds)));
        const around: BoundingRectangle = features.reduce(
            (bb: BoundingRectangle, feature: MapFeature) => bb.includeAll(feature.bounds),
            new BoundingRectangle()
        );
        return new MappedStations(stations, features, around.zoomOutOrAround(DefaultLocation, DefaultMargin));
    }

    constructor(
        public readonly stations: DisplayStation[] = [],
        public readonly features: MapFeature[] = [],
        public readonly bounds: BoundingRectangle | null = null
    ) {}

    public get valid(): boolean {
        return this.bounds != null;
    }

    public boundsLngLat(): LngLat[] | null {
        if (this.bounds) {
            return this.bounds.lngLat();
        }
        return null;
    }

    public overrideBounds(bounds: [LngLat, LngLat]): MappedStations {
        return new MappedStations(this.stations, this.features, new BoundingRectangle(bounds[0], bounds[1]));
    }

    public focusOn(id: number): MappedStations {
        const station = this.stations.find((s) => s.id == id);
        if (!station || !station.location) {
            return this;
        }
        const centered = new BoundingRectangle().include(station.location.lngLat());
        return new MappedStations(this.stations, this.features, centered.zoomOutOrAround(DefaultLocation, DefaultMargin));
    }
}

export class DisplayProject {
    public readonly id: number;
    public readonly name: string;
    public readonly modules: ProjectModule[];
    public readonly duration: number | null = null;
    public readonly mapped: MappedStations;
    public readonly places: { native: string | null } = { native: null };

    constructor(public readonly project: Project, public readonly users: ProjectUser[], public readonly stations: DisplayStation[]) {
        this.id = project.id;
        this.name = project.name;
        this.modules = _(stations)
            .filter((station) => station.configurations.all.length > 0)
            .map((station) => station.configurations.all[0].modules.filter((m) => !m.internal))
            .flatten()
            .uniqBy((m) => m.name)
            .map((m) => new ProjectModule(m.name, null))
            .value();
        this.mapped = MappedStations.make(stations);
        if (project.startTime && project.endTime) {
            const start = new Date(project.startTime);
            const end = new Date(project.endTime);
            end.setDate(end.getDate() + 1); // FK-1927
            this.duration = end.getTime() - start.getTime();
        }
        this.places.native = _(stations)
            .map((s) => s.placeNameNative)
            .filter((n) => n != null && n.length > 0)
            .uniq()
            .join(", ");
    }
}

const getters = {
    projectsById(state: StationsState): { [index: number]: DisplayProject } {
        return _(Object.values(state.projects))
            .map((p) => {
                const users = state.projectUsers[p.id];
                const stations = state.projectStations[p.id];
                if (stations && users) {
                    return [p.id, new DisplayProject(p, users, stations)];
                }
                return [p.id, null];
            })
            .fromPairs()
            .value();
    },
    stationsById(state: StationsState): { [index: number]: DisplayStation } {
        return { ...state.stations, ...state.user.stations };
    },
    mapped(state: StationsState): MappedStations | null {
        return state.mapped;
    },
};

const actions = (services: Services) => {
    return {
        [ActionTypes.INITIALIZE]: async ({ dispatch }: { dispatch: any }) => {
            await dispatch(ActionTypes.NEED_COMMON);
        },
        [ActionTypes.AUTHENTICATED]: async ({ dispatch }: { dispatch: any }) => {
            await dispatch(ActionTypes.NEED_COMMON);
        },
        [ActionTypes.NEED_COMMON]: async ({ dispatch }: { dispatch: any }) => {
            await dispatch(ActionTypes.NEED_PROJECTS);
            await dispatch(ActionTypes.NEED_STATIONS);
        },
        [ActionTypes.NEED_PROJECTS]: async ({ commit, dispatch, state }: { commit: any; dispatch: any; state: StationsState }) => {
            commit(MutationTypes.LOADING, { projects: true });
            const [communityProjects, userProjects] = await Promise.all([
                services.api.getPublicProjects(),
                services.api.getUserProjects(() => {
                    return Promise.resolve({
                        projects: [],
                    });
                }),
            ]);
            commit(HAVE_COMMUNITY_PROJECTS, communityProjects.projects);
            commit(HAVE_USER_PROJECTS, userProjects.projects);
            commit(MutationTypes.LOADING, { projects: false });
        },
        [ActionTypes.NEED_STATIONS]: async ({ commit, dispatch, state }: { commit: any; dispatch: any; state: StationsState }) => {
            commit(MutationTypes.LOADING, { stations: true });
            const [stations] = await Promise.all([
                services.api.getUserStations(() => {
                    return Promise.resolve({
                        stations: [],
                    });
                }),
            ]);
            commit(HAVE_USER_STATIONS, stations.stations);
            commit(MutationTypes.LOADING, { stations: false });
        },
        [ActionTypes.NEED_PROJECT]: async (
            { commit, dispatch, state }: { commit: any; dispatch: any; state: StationsState },
            payload: { id: number }
        ) => {
            commit(MutationTypes.LOADING, { projects: true });

            const [project, users, stations] = await Promise.all([
                services.api.getProject(payload.id),
                services.api.getUsersByProject(payload.id),
                services.api.getStationsByProject(payload.id),
            ]);

            commit(PROJECT_LOADED, project);
            commit(PROJECT_USERS, { projectId: payload.id, users: users.users });
            commit(PROJECT_STATIONS, { projectId: payload.id, stations: stations.stations });

            commit(MutationTypes.LOADING, { projects: false });
        },
        [ActionTypes.NEED_STATION]: async (
            { commit, dispatch, state }: { commit: any; dispatch: any; state: StationsState },
            payload: { id: number }
        ) => {
            commit(MutationTypes.LOADING, { stations: true });

            const [station] = await Promise.all([services.api.getStation(payload.id)]);

            commit(STATION_UPDATE, station);

            commit(MutationTypes.LOADING, { stations: false });
        },
        [ActionTypes.PROJECT_FOLLOW]: async ({ commit, dispatch }: { commit: any; dispatch: any }, payload: { projectId: number }) => {
            await services.api.followProject(payload.projectId);
            commit(PROJECT_LOADED, await services.api.getProject(payload.projectId));
        },
        [ActionTypes.PROJECT_UNFOLLOW]: async ({ commit, dispatch }: { commit: any; dispatch: any }, payload: { projectId: number }) => {
            await services.api.unfollowProject(payload.projectId);
            commit(PROJECT_LOADED, await services.api.getProject(payload.projectId));
        },
        [ActionTypes.STATION_PROJECT_ADD]: async (
            { commit, dispatch }: { commit: any; dispatch: any },
            payload: { stationId: number; projectId: number }
        ) => {
            await services.api.addStationToProject(payload);

            const stations = await services.api.getStationsByProject(payload.projectId);
            commit(PROJECT_STATIONS, { projectId: payload.projectId, stations: stations.stations });
        },
        [ActionTypes.STATION_PROJECT_REMOVE]: async (
            { commit, dispatch }: { commit: any; dispatch: any },
            payload: { stationId: number; projectId: number }
        ) => {
            await services.api.removeStationFromProject(payload);

            const stations = await services.api.getStationsByProject(payload.projectId);
            commit(PROJECT_STATIONS, { projectId: payload.projectId, stations: stations.stations });
        },
        [ActionTypes.PROJECT_INVITE]: async (
            { commit, dispatch }: { commit: any; dispatch: any },
            payload: { projectId: number; email: string; role: string }
        ) => {
            await services.api.sendInvite(payload);
            const usersReply = await services.api.getUsersByProject(payload.projectId);
            commit(PROJECT_USERS, { projectId: payload.projectId, users: usersReply.users });
        },
        [ActionTypes.PROJECT_REMOVE]: async (
            { commit, dispatch }: { commit: any; dispatch: any },
            payload: { projectId: number; email: string }
        ) => {
            await services.api.removeUserFromProject(payload);
            const usersReply = await services.api.getUsersByProject(payload.projectId);
            commit(PROJECT_USERS, { projectId: payload.projectId, users: usersReply.users });
        },
        [ActionTypes.ACCEPT_PROJECT]: async ({ commit, dispatch }: { commit: any; dispatch: any }, payload: { projectId: number }) => {
            await services.api.acceptProjectInvite(payload);

            const userProjects = await services.api.getUserProjects(OnNoReject);
            commit(HAVE_USER_PROJECTS, userProjects.projects);
        },
        [ActionTypes.DECLINE_PROJECT]: async ({ commit, dispatch }: { commit: any; dispatch: any }, payload: { projectId: number }) => {
            await services.api.declineProjectInvite(payload);
        },
        [ActionTypes.ACCEPT_PROJECT_INVITE]: async (
            { commit, dispatch }: { commit: any; dispatch: any },
            payload: { id: number; token: string }
        ) => {
            await services.api.acceptInvite(payload);

            const userProjects = await services.api.getUserProjects(OnNoReject);
            commit(HAVE_USER_PROJECTS, userProjects.projects);
        },
        [ActionTypes.DECLINE_PROJECT_INVITE]: async (
            { commit, dispatch }: { commit: any; dispatch: any },
            payload: { id: number; token: string }
        ) => {
            await services.api.declineInvite(payload);
        },
        [ActionTypes.DELETE_PROJECT]: async ({ commit, dispatch }: { commit: any; dispatch: any }, payload: { projectId: number }) => {
            await services.api.deleteProject(payload);

            commit(PROJECT_DELETED, payload);

            return payload;
        },
        [ActionTypes.ADD_PROJECT]: async ({ commit, dispatch }: { commit: any; dispatch: any }, payload: any) => {
            const project = await services.api.addProject(payload);

            commit(PROJECT_UPDATE, project);

            return project;
        },
        [ActionTypes.SAVE_PROJECT]: async ({ commit, dispatch }: { commit: any; dispatch: any }, payload: any) => {
            const project = await services.api.updateProject(payload);

            commit(PROJECT_UPDATE, project);

            return project;
        },
    };
};

const mutations = {
    [HAVE_COMMUNITY_PROJECTS]: (state: StationsState, projects: Project[]) => {
        Vue.set(
            state.community,
            "projects",
            _.keyBy(projects, (p) => p.id)
        );
        Vue.set(state, "projects", { ...state.projects, ..._.keyBy(projects, (p) => p.id) });
    },
    [HAVE_USER_PROJECTS]: (state: StationsState, projects: Project[]) => {
        Vue.set(
            state.user,
            "projects",
            _.keyBy(projects, (p) => p.id)
        );
        Vue.set(state, "projects", { ...state.projects, ..._.keyBy(projects, (p) => p.id) });
    },
    [HAVE_USER_STATIONS]: (state: StationsState, payload: Station[]) => {
        const stations = payload.map((station) => new DisplayStation(station));
        Vue.set(
            state.user,
            "stations",
            _.keyBy(stations, (s) => s.id)
        );
        Vue.set(state, "stations", { ...state.stations, ..._.keyBy(stations, (s) => s.id) });
        Vue.set(state, "hasNoStations", stations.length == 0);
        Vue.set(state, "mapped", MappedStations.make(stations));
    },
    [PROJECT_USERS]: (state: StationsState, payload: { projectId: number; users: ProjectUser[] }) => {
        Vue.set(state.projectUsers, payload.projectId, payload.users);
    },
    [PROJECT_FOLLOWS]: (state: StationsState, payload: { projectId: number; followers: ProjectFollowers }) => {
        Vue.set(state.projectFollowers, payload.projectId, payload.followers);
    },
    [PROJECT_STATIONS]: (state: StationsState, payload: { projectId: number; stations: Station[] }) => {
        const projectStations = payload.stations.map((s) => new DisplayStation(s));
        state.stations = { ...state.stations, ..._.keyBy(projectStations, (s) => s.id) };
        Vue.set(state.projectStations, payload.projectId, projectStations);
    },
    [PROJECT_ACTIVITY]: (state: StationsState, payload: { projectId: number; activities: Activity[] }) => {
        Vue.set(state.projectActivities, payload.projectId, payload.activities);
    },
    [STATION_UPDATE]: (state: StationsState, payload: Station) => {
        Vue.set(state.stations, payload.id, new DisplayStation(payload));
    },
    [PROJECT_LOADED]: (state: StationsState, project: Project) => {
        Vue.set(state.projects, project.id, project);
        if (state.user.projects[project.id]) {
            Vue.set(state.user.projects, project.id, project);
        }
        if (state.community.projects[project.id]) {
            Vue.set(state.community.projects, project.id, project);
        }
    },
    [PROJECT_UPDATE]: (state: StationsState, project: Project) => {
        Vue.set(state.projects, project.id, project);
        Vue.set(state.user.projects, project.id, project);
        if (state.community.projects[project.id]) {
            Vue.set(state.community.projects, project.id, project);
        }
    },
    [PROJECT_DELETED]: (state: StationsState, payload: { projectId: number }) => {
        delete state.projects[payload.projectId];
        delete state.user.projects[payload.projectId];
        delete state.community.projects[payload.projectId];
    },
};

export const stations = (services: Services) => {
    const state = () => new StationsState();

    return {
        namespaced: false,
        state,
        getters,
        actions: actions(services),
        mutations,
    };
};
