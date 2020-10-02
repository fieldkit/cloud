import _ from "lodash";
import Vue from "vue";
import * as MutationTypes from "../mutations";
import { Services } from "@/api";

type LoadingPayload = { [index: string]: boolean };

export class ProgressState {
    loading: { [index: string]: boolean } = {};
}

const getters = {
    isBusy(state: ProgressState) {
        return _.some(Object.values(state.loading).filter((l) => l));
    },
};

const actions = (services: Services) => {
    return {};
};

const mutations = {
    [MutationTypes.LOADING]: (state: ProgressState, payload: LoadingPayload) => {
        Vue.set(state, "loading", { ...state.loading, ...payload });
    },
};

export const progress = (services: Services) => {
    const state = () => new ProgressState();

    return {
        namespaced: false,
        state,
        getters,
        actions: actions(services),
        mutations,
    };
};
