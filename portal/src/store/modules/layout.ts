import _ from "lodash";
import * as MutationTypes from "../mutations";
import { DisplayStation } from "./stations";
import { GlobalState, GlobalGetters } from "./global";

import { Services, CurrentUser, Project } from "@/api";

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

const actions = (services: Services) => {
    return {};
};

const mutations = {};

export const layout = (services: Services) => {
    const state = () => new LayoutState();

    return {
        namespaced: false,
        state,
        getters,
        actions: actions(services),
        mutations,
    };
};
