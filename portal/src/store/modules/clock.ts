import Vue from "vue";
import * as MutationTypes from "../mutations";
import { Services } from "@/api";

export class Clock {
    public now: Date = new Date();
    public unix: number = Math.round(new Date().getTime() / 1000);
}

export class ClockState {
    public wall: Clock = new Clock();
}

const getters = {};

const actions = (services: Services) => {
    return {};
};

const mutations = {
    [MutationTypes.TICK]: (state: ClockState) => {
        Vue.set(state, "wall", new Clock());
    },
};

export const clock = (services: Services) => {
    const state = () => new ClockState();

    return {
        namespaced: false,
        state,
        getters,
        actions: actions(services),
        mutations,
    };
};
