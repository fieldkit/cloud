import Vue from "vue";
import { Services } from "@/api";
import * as ActionTypes from "@/store/actions";
import * as MutationTypes from "@/store/mutations";

export class DiscussionState {
    dataEvents
}

const getters = {
    dataEvents(state) {
        return state.dataEvents;
    },
};

const actions = (services: Services) => {
    return {
        [ActionTypes.NEED_DATA_EVENTS]: async (
            { commit, dispatch, state }: { commit: any; dispatch: any; state: DiscussionState },
            payload: { bookmark: string }
        ) => {
            const dataEvents = await services.api.getDataEvents(payload.bookmark);
            commit(MutationTypes.DATA_EVENTS_UPDATE, dataEvents);
        },
    };
};

const mutations = {
    [MutationTypes.DATA_EVENTS_UPDATE]: (
        state: DiscussionState,
        payload: {
            events
        }
    ) => {
        console.log("DE MUTATION", payload)
        Vue.set(state, "dataEvents", payload.events);
    },
};

export const dataEvents = (services: Services) => {
    const state = () => new DiscussionState();

    return {
        namespaced: false,
        state,
        getters,
        actions: actions(services),
        mutations,
    };
};
