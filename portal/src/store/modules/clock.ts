import Vue from "vue";
import * as MutationTypes from "../mutations";

export class Clock {
    public now: Date = new Date();
    public unix: number = Math.round(new Date().getTime() / 1000);
}

export class ClockState {
    public wall: Clock = new Clock();
}

const getters = {};

const actions = {};

const mutations = {
    [MutationTypes.TICK]: (state: ClockState) => {
        Vue.set(state, "wall", new Clock());
    },
};

const state = () => new ClockState();

export const clock = {
    namespaced: false,
    state,
    getters,
    actions,
    mutations,
};
