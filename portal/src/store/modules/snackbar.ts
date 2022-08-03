import Vue from "vue";
import * as MutationTypes from "../mutations";
import * as ActionTypes from "@/store/actions";

export enum SnackbarStyle {
    neutral = "netural",
    success = "success",
    fail = "fail",
}

export class SnackbarState {
    visible: boolean = false;
    message: string = "";
    style: SnackbarStyle;
}

const getters = {
    visible(state: SnackbarState) {
        return state.visible;
    },
    message(state: SnackbarState) {
        return state.message;
    },
};

const mutations = {
    [MutationTypes.SHOW_SNACKBAR]: (state: SnackbarState, payload: { message: string; type: SnackbarStyle }) => {
        Vue.set(state, "message", payload.message);
        Vue.set(state, "visible", true);
        Vue.set(state, "style", payload.type);
    },
    [MutationTypes.HIDE_SNACKBAR]: (state: SnackbarState) => {
        Vue.set(state, "message", "");
        Vue.set(state, "visible", false);
        Vue.set(state, "style", null);
    },
};

const actions = () => {
    return {
        [ActionTypes.SHOW_SNACKBAR]: async (
            { commit, dispatch, state }: { commit: any; dispatch: any; state: SnackbarState },
            message: string
        ) => {
            commit(MutationTypes.SHOW_SNACKBAR, message);
            setTimeout(() => {
                dispatch(ActionTypes.HIDE_SNACKBAR);
            }, 5000); // persitence time is actually done via the animation from snackbar comp, this is just to clear the state
        },
        [ActionTypes.HIDE_SNACKBAR]: async ({ commit, dispatch, state }: { commit: any; dispatch: any; state: SnackbarState }) => {
            commit(MutationTypes.HIDE_SNACKBAR);
        },
    };
};

export const snackbar = () => {
    const state = () => new SnackbarState();

    return {
        namespaced: false,
        state,
        actions: actions(),
        getters,
        mutations,
    };
};
