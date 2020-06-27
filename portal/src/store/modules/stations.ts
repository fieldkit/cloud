import Vue from "vue";
import * as ActionTypes from "../actions";
import FKApi, { Station, Project } from "../../api/api";

export const HAVE_USER_STATIONS = "HAVE_USER_STATIONS";
export const HAVE_USER_PROJECTS = "HAVE_USER_PROJECTS";
export const HAVE_COMMUNITY_PROJECTS = "HAVE_COMMUNITY_PROJECTS";
export const LOADING = "LOADING";

export class StationsState {
    stations: { user: Station[] } = { user: [] };
    projects: { user: Project[]; community: Project[] } = { user: [], community: [] };
    loading: { stations: boolean; projects: boolean };
}

const getters = {};

const actions = {
    [ActionTypes.NEED_PROJECTS]: async ({ commit, dispatch, state }: { commit: any; dispatch: any; state: StationsState }) => {
        commit(LOADING, { projects: true });
        const [communityProjects, userProjects] = await Promise.all([new FKApi().getPublicProjects(), new FKApi().getUserProjects()]);
        commit(HAVE_COMMUNITY_PROJECTS, communityProjects.projects);
        commit(HAVE_USER_PROJECTS, userProjects.projects);
        commit(LOADING, { projects: false });
    },
    [ActionTypes.NEED_STATIONS]: async ({ commit, dispatch, state }: { commit: any; dispatch: any; state: StationsState }) => {
        commit(LOADING, { stations: false });
        const [stations] = await Promise.all([new FKApi().getStations()]);
        commit(HAVE_USER_STATIONS, stations.stations);
        commit(LOADING, { stations: false });
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
    },
    [LOADING]: (state: StationsState, payload: { projects: boolean; stations: boolean }) => {
        Vue.set(state, "loading", { ...state.loading, ...payload });
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
