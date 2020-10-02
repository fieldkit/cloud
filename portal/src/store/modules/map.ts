import { Services } from "@/api";

export class MapState {}

const getters = {};

const actions = (services: Services) => {
    return {};
};

const mutations = {};

export const map = (services: Services) => {
    const state = () => new MapState();

    return {
        namespaced: false,
        state,
        getters,
        actioins: actions(services),
        mutations,
    };
};
