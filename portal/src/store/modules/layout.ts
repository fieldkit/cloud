import _ from "lodash";
import * as MutationTypes from "../mutations";
import FKApi, { CurrentUser, Project } from "../../api/api";
import { DisplayStation } from "./stations";
import { GlobalState, GlobalGetters } from "./global";

export class LayoutState {
    user: CurrentUser | null;
    users: {
        stations: DisplayStation[];
        projects: Project[];
    } = {
        stations: [],
        projects: [],
    };
    community: {
        projects: Project[];
    } = {
        projects: [],
    };
}

const getters = {
    layout: (state: LayoutState, getters: any, rootState: GlobalState, rootGetters: GlobalGetters) => {
        return {};
    },
};

const actions = {};

const mutations = {};

const state = () => new LayoutState();

export const layout = {
    namespaced: false,
    state,
    getters,
    actions,
    mutations,
};
