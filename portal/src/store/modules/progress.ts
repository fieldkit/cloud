import _ from "lodash";
import Vue from "vue";
import * as MutationTypes from "../mutations";

type LoadingPayload = { [index: string]: boolean };

export class ProgressState {
    loading: { [index: string]: boolean };
}

const getters = {
    isBusy(state: ProgressState) {
        return _.some(Object.values(state.loading).filter(l => l));
    },
};

const actions = {};

const mutations = {
    [MutationTypes.LOADING]: (state: ProgressState, payload: LoadingPayload) => {
        Vue.set(state, "loading", { ...state.loading, ...payload });
    },
};

const state = () => new ProgressState();

export const progress = {
    namespaced: false,
    state,
    getters,
    actions,
    mutations,
};
