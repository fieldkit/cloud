// import _ from "lodash";
// import * as ActionTypes from "../actions";
import * as MutationTypes from "../mutations";

export class StationsState {}

const getters = {};

const actions = {};

const mutations = {
    [MutationTypes.INITIALIZE]: (state: StationsState) => {
        console.log("ready");
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
