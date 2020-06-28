import Vue from "vue";
import * as ActionTypes from "../actions";
import * as MutationTypes from "../mutations";
import TokenStorage from "@/api/tokens";
import FkApi, { LoginPayload, LoginResponse, CurrentUser } from "@/api/api";

export const UPDATE_TOKEN = "UPDATE_TOKEN";
export const CURRENT_USER = "CURRENT_USER";

export const REFRESH_CURRENT_USER = "REFRESH_CURRENT_USER";

export class UserState {
    token: string | null = null;
    user: CurrentUser | null = null;
}

const getters = {
    isAuthenticated: (state: UserState) => {
        return state.token !== null;
    },
};

const actions = {
    [ActionTypes.INITIALIZE]: ({ commit, dispatch }: { commit: any; dispatch: any }) => {
        return dispatch(REFRESH_CURRENT_USER);
    },
    [ActionTypes.AUTHENTICATE]: ({ commit, dispatch, state }: { commit: any; dispatch: any; state: UserState }, payload: LoginPayload) => {
        return new FkApi().login(payload.email, payload.password).then(token => {
            commit(UPDATE_TOKEN, token);
            return dispatch(REFRESH_CURRENT_USER).then(() => {
                return new LoginResponse(token);
            });
        });
    },
    [ActionTypes.LOGOUT]: ({ commit }: { commit: any }) => {
        commit(UPDATE_TOKEN, null);
        commit(CURRENT_USER, null);
    },
    [REFRESH_CURRENT_USER]: ({ commit, dispatch, state }: { commit: any; dispatch: any; state: UserState }) => {
        return new FkApi().getCurrentUser().then(user => {
            commit(CURRENT_USER, user);
            return user;
        });
    },
};

const mutations = {
    [MutationTypes.INITIALIZE]: (state: UserState, token: string) => {
        Vue.set(state, "token", new TokenStorage().getToken());
    },
    [UPDATE_TOKEN]: (state: UserState, token: string) => {
        Vue.set(state, "token", token);
    },
    [CURRENT_USER]: (state: UserState, user: CurrentUser) => {
        Vue.set(state, "user", user);
    },
};

const state = () => new UserState();

export const user = {
    namespaced: false,
    state,
    getters,
    actions,
    mutations,
};
