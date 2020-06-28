import _ from "lodash";
import Vue from "vue";
import * as MutationTypes from "../mutations";
import * as ActionTypes from "../actions";
import FKApi, {
    Station,
    StationModule,
    ModuleSensor,
    Project,
    ProjectUser,
    ProjectFollowers,
    Activity,
    Configurations,
    HasLocation,
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

export class StationsState {
    stations: { all: { [index: number]: Station }; user: Station[] } = { all: {}, user: [] };
    projects: { user: Project[]; community: Project[] } = { user: [], community: [] };
    projectUsers: { [index: number]: ProjectUser[] } = {};
    projectFollowers: { [index: number]: ProjectFollowers } = {};
    projectStations: { [index: number]: DisplayStation[] } = {};
    projectActivities: { [index: number]: Activity[] } = {};
    hasNoStations: boolean;
}

export function whenWasStationUpdated(station: Station): Date {
    const uploads = station.uploads.map(u => u.time);
    if (uploads.length > 0) {
        return new Date(uploads[0]);
    }
    return new Date(station.updated);
}

export class Location implements HasLocation {
    constructor(public readonly latitude: number, public readonly longitude) {}

    clone(): Location {
        return new Location(this.latitude, this.longitude);
    }
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
        this.sensors = module.sensors.map(s => new DisplaySensor(s));
    }
}

export class DisplayStation {
    id: number;
    name: string;
    updated: Date;
    configurations: Configurations;
    receivedAt: Date | null = null;
    totalReadings = 0;
    location: Location | null = null;
    photos: Photos;
    modules: DisplayModule[] = [];

    constructor(station: Station) {
        this.id = station.id;
        this.name = station.name;
        this.updated = whenWasStationUpdated(station);
        this.configurations = station.configurations;
        this.photos = station.photos;
        this.modules =
            _(station.configurations.all)
                .map(c => c.modules.filter(m => !m.internal).map(m => new DisplayModule(m)))
                .head() || [];
        if (station.location) {
            this.location = new Location(station.location.latitude, station.location.longitude);
        }
    }
}

export class ProjectModule {
    constructor(public readonly name: string, public url: string | null) {}
}

export class DisplayProject {
    id: number;
    modules: ProjectModule[];
    duration: number | null = null;

    constructor(
        public readonly project: Project,
        public readonly users: ProjectUser[],
        public readonly stations: DisplayStation[],
        public readonly activities: Activity[]
    ) {
        this.id = project.id;
        this.modules = _(stations)
            .filter(station => station.configurations.all.length > 0)
            .map(station => station.configurations.all[0].modules.filter(m => !m.internal))
            .flatten()
            .uniqBy(m => m.name)
            .map(m => new ProjectModule(m.name, null))
            .value();
    }
}

const getters = {
    projectsById(state: StationsState): { [index: number]: DisplayProject } {
        return _(state.projects.user)
            .map(p => {
                const users = state.projectUsers[p.id] || [];
                const stations = state.projectStations[p.id] || [];
                const activities = state.projectActivities[p.id] || [];
                return new DisplayProject(p, users, stations, activities);
            })
            .keyBy(p => p.project.id)
            .value();
    },
    stationsById(state: StationsState): { [index: number]: Station } {
        return _.keyBy(state.stations, p => p.id);
    },
};

const actions = {
    [ActionTypes.NEED_COMMON]: async ({ dispatch }: { dispatch: any }) => {
        await dispatch(ActionTypes.NEED_PROJECTS);
        await dispatch(ActionTypes.NEED_STATIONS);
    },
    [ActionTypes.NEED_PROJECTS]: async ({ commit, dispatch, state }: { commit: any; dispatch: any; state: StationsState }) => {
        commit(MutationTypes.LOADING, { projects: true });
        const [communityProjects, userProjects] = await Promise.all([new FKApi().getPublicProjects(), new FKApi().getUserProjects()]);
        commit(HAVE_COMMUNITY_PROJECTS, communityProjects.projects);
        commit(HAVE_USER_PROJECTS, userProjects.projects);
        commit(MutationTypes.LOADING, { projects: false });
    },
    [ActionTypes.NEED_STATIONS]: async ({ commit, dispatch, state }: { commit: any; dispatch: any; state: StationsState }) => {
        commit(MutationTypes.LOADING, { stations: true });
        const [stations] = await Promise.all([new FKApi().getStations()]);
        commit(HAVE_USER_STATIONS, stations.stations);
        commit(MutationTypes.LOADING, { stations: false });
    },
    [ActionTypes.NEED_PROJECT]: async (
        { commit, dispatch, state }: { commit: any; dispatch: any; state: StationsState },
        payload: { id: number }
    ) => {
        commit(MutationTypes.LOADING, { projects: true });

        const api = new FKApi();
        const [users, followers, stations, activity] = await Promise.all([
            api.getUsersByProject(payload.id),
            api.getProjectFollows(payload.id),
            api.getStationsByProject(payload.id),
            api.getProjectActivity(payload.id),
        ]);

        commit(PROJECT_USERS, { projectId: payload.id, users: users.users });
        commit(PROJECT_FOLLOWS, { projectId: payload.id, followers: followers });
        commit(PROJECT_STATIONS, { projectId: payload.id, stations: stations.stations });
        commit(PROJECT_ACTIVITY, { projectId: payload.id, activities: activity.activities });

        commit(MutationTypes.LOADING, { projects: false });
    },
    [ActionTypes.NEED_STATION]: async (
        { commit, dispatch, state }: { commit: any; dispatch: any; state: StationsState },
        payload: { id: number }
    ) => {
        commit(MutationTypes.LOADING, { stations: true });

        const api = new FKApi();
        const [station, _needStations, _needProjects] = await Promise.all([
            api.getStation(payload.id),
            dispatch(ActionTypes.NEED_STATIONS),
            dispatch(ActionTypes.NEED_PROJECTS),
        ]);

        commit(STATION_UPDATE, station);

        if (station.location) {
            const [nativeLands, placeNames] = await Promise.all([api.getNativeLand(station.location), api.getPlaceName(station.location)]);
            console.log(nativeLands, placeNames);
        }

        commit(MutationTypes.LOADING, { stations: false });
    },
    [ActionTypes.PROJECT_FOLLOW]: async ({ commit, dispatch }: { commit: any; dispatch: any }, payload: { projectId: number }) => {
        await new FKApi().followProject(payload.projectId);
    },
    [ActionTypes.PROJECT_UNFOLLOW]: async ({ commit, dispatch }: { commit: any; dispatch: any }, payload: { projectId: number }) => {
        await new FKApi().unfollowProject(payload.projectId);
    },
    [ActionTypes.STATION_PROJECT_ADD]: async ({ commit, dispatch }: { commit: any; dispatch: any }, payload: { projectId: number }) => {
        // await new FKApi().unfollowProject(payload.projectId);
    },
    [ActionTypes.STATION_PROJECT_REMOVE]: async ({ commit, dispatch }: { commit: any; dispatch: any }, payload: { projectId: number }) => {
        // await new FKApi().unfollowProject(payload.projectId);
    },
};

const mutations = {
    [HAVE_COMMUNITY_PROJECTS]: (state: StationsState, projects: Project[]) => {
        state.projects.community = projects;
    },
    [HAVE_USER_PROJECTS]: (state: StationsState, projects: Project[]) => {
        state.projects.user = projects;
    },
    [HAVE_USER_STATIONS]: (state: StationsState, stations: Station[]) => {
        state.stations.user = stations;
        state.hasNoStations = stations.length == 0;
        stations.forEach(station => {
            Vue.set(state.stations.all, station.id, new DisplayStation(station));
        });
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
            payload.stations.map(s => new DisplayStation(s))
        );
    },
    [PROJECT_ACTIVITY]: (state: StationsState, payload: { projectId: number; activities: Activity[] }) => {
        Vue.set(state.projectActivities, payload.projectId, payload.activities);
    },
    [STATION_UPDATE]: (state: StationsState, payload: Station) => {
        Vue.set(state.stations.all, payload.id, new DisplayStation(payload));
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
