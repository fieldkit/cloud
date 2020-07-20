import _ from "lodash";
import Vue from "vue";
import * as MutationTypes from "../mutations";
import * as ActionTypes from "../actions";
import { Location, BoundingRectangle } from "../map-types";
import FKApi, {
    OnNoReject,
    Station,
    StationModule,
    ModuleSensor,
    Project,
    ProjectUser,
    ProjectFollowers,
    Activity,
    Configurations,
    Photos,
} from "../../api/api";

export const HAVE_USER_STATIONS = "HAVE_USER_STATIONS";
export const HAVE_USER_PROJECTS = "HAVE_USER_PROJECTS";
export const HAVE_COMMUNITY_PROJECTS = "HAVE_COMMUNITY_PROJECTS";
export const PROJECT_USERS = "PROJECT_USERS";
export const PROJECT_FOLLOWS = "PROJECT_FOLLOWS";
export const PROJECT_STATIONS = "PROJECT_STATIONS";
export const PROJECT_ACTIVITY = "PROJECT_ACTIVITY";
export const STATION_UPDATE = "STATION_UPDATE";
export const PROJECT_UPDATE = "PROJECT_UPDATE";

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
}

export function whenWasStationUpdated(station: Station): Date {
    const uploads = station.uploads.map((u) => u.time);
    if (uploads.length > 0) {
        return new Date(uploads[0]);
    }
    return new Date(station.updated);
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
    id: number;
    name: string;
    updated: Date;
    configurations: Configurations;
    receivedAt: Date | null = null;
    deployedAt: Date | null = null;
    totalReadings = 0;
    location: Location | null = null;
    photos: Photos;
    modules: DisplayModule[] = [];
    placeNameOther: string | null;
    placeNameNative: string | null;

    constructor(station: Station) {
        this.id = station.id;
        this.name = station.name;
        this.updated = whenWasStationUpdated(station);
        this.configurations = station.configurations;
        this.photos = station.photos;
        this.placeNameOther = station.placeNameOther;
        this.placeNameNative = station.placeNameNative;
        this.deployedAt = station.recordingStartedAt;
        this.modules =
            _(station.configurations.all)
                .map((c) => c.modules.filter((m) => !m.internal).map((m) => new DisplayModule(m)))
                .head() || [];
        if (station.location) {
            this.location = new Location(station.location.latitude, station.location.longitude);
        }
    }
}

export class ProjectModule {
    constructor(public readonly name: string, public url: string | null) {}
}

export class MapFeature {
    type = "Feature";
    geometry: { type: string; coordinates: number[] } | null = null;
    properties: { title: string; icon: string; id: number } | null = null;

    constructor(station: DisplayStation) {
        this.geometry = {
            type: "Point",
            coordinates: station.location.lngLat(),
        };
        this.properties = {
            id: station.id,
            title: station.name,
            icon: "marker",
        };
    }
}

export class MappedStations {
    bounds: BoundingRectangle | null = null;
    features: MapFeature[] = [];

    constructor(stations: DisplayStation[]) {
        const DefaultLocation = new Location(34.3318104, -118.0730372); // TwinPeaks

        const located = stations.filter((station) => station.location != null);
        const around = located.reduce((bb, station) => bb.include(station.location), new BoundingRectangle());
        const feetAround = located.length > 0 ? 1000 : 100000;

        this.bounds = around.expandIfSingleCoordinate(DefaultLocation, feetAround);
        this.features = located.map((ds) => new MapFeature(ds));
    }

    get valid(): boolean {
        return this.bounds != null;
    }

    boundsLngLat(): number[][] {
        console.log("map: returning bounds", this.bounds, this.features);
        return [this.bounds.min.lngLat(), this.bounds.max.lngLat()];
    }
}

export class DisplayProject {
    id: number;
    name: string;
    modules: ProjectModule[];
    duration: number | null = null;
    mapped: MappedStations;
    places: { native: string | null } = { native: null };

    constructor(
        public readonly project: Project,
        public readonly users: ProjectUser[],
        public readonly stations: DisplayStation[],
        public readonly activities: Activity[]
    ) {
        this.id = project.id;
        this.name = project.name;
        this.modules = _(stations)
            .filter((station) => station.configurations.all.length > 0)
            .map((station) => station.configurations.all[0].modules.filter((m) => !m.internal))
            .flatten()
            .uniqBy((m) => m.name)
            .map((m) => new ProjectModule(m.name, null))
            .value();
        this.mapped = new MappedStations(stations);
        if (project.startTime && project.endTime) {
            const start = new Date(project.startTime);
            const end = new Date(project.endTime);
            end.setDate(end.getDate() + 1); // FK-1927
            this.duration = end.getTime() - start.getTime();
        }
        this.places.native = _(stations)
            .map((s) => s.placeNameNative)
            .uniq()
            .filter((n) => n && n.length > 0)
            .join(", ");
    }
}

const getters = {
    projectsById(state: StationsState): { [index: number]: DisplayProject } {
        return _(Object.values(state.projects))
            .map((p) => {
                const users = state.projectUsers[p.id] || [];
                const stations = state.projectStations[p.id] || [];
                const activities = state.projectActivities[p.id] || [];
                return new DisplayProject(p, users, stations, activities);
            })
            .keyBy((p) => p.project.id)
            .value();
    },
    stationsById(state: StationsState): { [index: number]: Station } {
        return _.keyBy(state.stations, (p) => p.id);
    },
    mapped(state: StationsState): MappedStations {
        return new MappedStations(Object.values(state.stations));
    },
};

const actions = {
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
            new FKApi().getPublicProjects(),
            new FKApi().getUserProjects(() => {
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
            new FKApi().getUserStations(() => {
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

        const api = new FKApi();
        const [project, users, stations] = await Promise.all([
            api.getProject(payload.id),
            api.getUsersByProject(payload.id),
            api.getStationsByProject(payload.id),
        ]);

        commit(PROJECT_UPDATE, project);
        commit(PROJECT_USERS, { projectId: payload.id, users: users.users });
        commit(PROJECT_STATIONS, { projectId: payload.id, stations: stations.stations });

        commit(MutationTypes.LOADING, { projects: false });
    },
    [ActionTypes.NEED_STATION]: async (
        { commit, dispatch, state }: { commit: any; dispatch: any; state: StationsState },
        payload: { id: number }
    ) => {
        commit(MutationTypes.LOADING, { stations: true });

        const [station] = await Promise.all([new FKApi().getStationFromVuex(payload.id)]);

        commit(STATION_UPDATE, station);

        commit(MutationTypes.LOADING, { stations: false });
    },
    [ActionTypes.PROJECT_FOLLOW]: async ({ commit, dispatch }: { commit: any; dispatch: any }, payload: { projectId: number }) => {
        await new FKApi().followProject(payload.projectId);
    },
    [ActionTypes.PROJECT_UNFOLLOW]: async ({ commit, dispatch }: { commit: any; dispatch: any }, payload: { projectId: number }) => {
        await new FKApi().unfollowProject(payload.projectId);
    },
    [ActionTypes.STATION_PROJECT_ADD]: async (
        { commit, dispatch }: { commit: any; dispatch: any },
        payload: { stationId: number; projectId: number }
    ) => {
        await new FKApi().addStationToProject(payload);

        const stations = await new FKApi().getStationsByProject(payload.projectId);
        commit(PROJECT_STATIONS, { projectId: payload.projectId, stations: stations.stations });
    },
    [ActionTypes.STATION_PROJECT_REMOVE]: async (
        { commit, dispatch }: { commit: any; dispatch: any },
        payload: { stationId: number; projectId: number }
    ) => {
        await new FKApi().removeStationFromProject(payload);

        const stations = await new FKApi().getStationsByProject(payload.projectId);
        commit(PROJECT_STATIONS, { projectId: payload.projectId, stations: stations.stations });
    },
    [ActionTypes.PROJECT_INVITE]: async (
        { commit, dispatch }: { commit: any; dispatch: any },
        payload: { projectId: number; email: string; role: string }
    ) => {
        await new FKApi().sendInvite(payload);
        const usersReply = await new FKApi().getUsersByProject(payload.projectId);
        commit(PROJECT_USERS, { projectId: payload.projectId, users: usersReply.users });
    },
    [ActionTypes.PROJECT_REMOVE]: async (
        { commit, dispatch }: { commit: any; dispatch: any },
        payload: { projectId: number; email: string }
    ) => {
        await new FKApi().removeUserFromProject(payload);
        const usersReply = await new FKApi().getUsersByProject(payload.projectId);
        commit(PROJECT_USERS, { projectId: payload.projectId, users: usersReply.users });
    },
    [ActionTypes.ACCEPT_PROJECT_INVITE]: async (
        { commit, dispatch }: { commit: any; dispatch: any },
        payload: { id: number; token: string }
    ) => {
        await new FKApi().acceptInvite(payload);

        const userProjects = await new FKApi().getUserProjects(OnNoReject);
        commit(HAVE_USER_PROJECTS, userProjects.projects);
    },
    [ActionTypes.DECLINE_PROJECT_INVITE]: async (
        { commit, dispatch }: { commit: any; dispatch: any },
        payload: { id: number; token: string }
    ) => {
        await new FKApi().declineInvite(payload);
    },
    [ActionTypes.SAVE_PROJECT]: async ({ commit, dispatch }: { commit: any; dispatch: any }, payload: any) => {
        const project = await new FKApi().updateProject(payload);

        commit(PROJECT_UPDATE, project);
    },
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
    },
    [PROJECT_USERS]: (state: StationsState, payload: { projectId: number; users: ProjectUser[] }) => {
        Vue.set(state.projectUsers, payload.projectId, payload.users);
    },
    [PROJECT_FOLLOWS]: (state: StationsState, payload: { projectId: number; followers: ProjectFollowers }) => {
        Vue.set(state.projectFollowers, payload.projectId, payload.followers);
    },
    [PROJECT_STATIONS]: (state: StationsState, payload: { projectId: number; stations: Station[] }) => {
        Vue.set(
            state.projectStations,
            payload.projectId,
            payload.stations.map((s) => new DisplayStation(s))
        );
    },
    [PROJECT_ACTIVITY]: (state: StationsState, payload: { projectId: number; activities: Activity[] }) => {
        Vue.set(state.projectActivities, payload.projectId, payload.activities);
    },
    [STATION_UPDATE]: (state: StationsState, payload: Station) => {
        Vue.set(state.stations, payload.id, new DisplayStation(payload));
    },
    [PROJECT_UPDATE]: (state: StationsState, project: Project) => {
        Vue.set(state.projects, project.id, project);
        Vue.set(state.user.projects, project.id, project);
    },
};

const state = () => new StationsState();

export const stations = {
    namespaced: false,
    state,
    getters,
    actions,
    mutations,
};
