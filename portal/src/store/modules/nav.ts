import * as MutationTypes from "../mutations";

export class NavigationState {}

const getters = {};

const actions = {};

const mutations = {
    [MutationTypes.NAVIGATION]: (state: NavigationState) => {},
};

const state = () => new NavigationState();

export const nav = {
    namespaced: false,
    state,
    getters,
    actions,
    mutations,
};
