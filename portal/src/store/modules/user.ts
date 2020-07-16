import Vue from "vue";
import * as ActionTypes from "../actions";
import * as MutationTypes from "../mutations";
import TokenStorage from "@/api/tokens";
import FkApi, { LoginPayload, LoginResponse, CurrentUser } from "@/api/api";

export const UPDATE_TOKEN = "UPDATE_TOKEN";
export const CURRENT_USER = "CURRENT_USER";

export class UserState {
    token: string | null = null;
    user: CurrentUser | null = null;
}

const getters = {
    isAuthenticated: (state: UserState) => {
        return state.token !== null;
    },
};

type ActionParameters = { commit: any; dispatch: any; state: UserState };

const actions = {
    [ActionTypes.INITIALIZE]: ({ commit, dispatch, state }: ActionParameters) => {
        if (state.token) {
            return dispatch(ActionTypes.REFRESH_CURRENT_USER);
        }
        return null;
    },
    [ActionTypes.LOGIN]: ({ commit, dispatch, state }: ActionParameters, payload: LoginPayload) => {
        return new FkApi().login(payload.email, payload.password).then((token) => {
            commit(UPDATE_TOKEN, token);
            return dispatch(ActionTypes.REFRESH_CURRENT_USER).then(() => {
                return dispatch(ActionTypes.AUTHENTICATED).then(() => {
                    return new LoginResponse(token);
                });
            });
        });
    },
    [ActionTypes.LOGOUT]: ({ commit }: { commit: any }) => {
        commit(UPDATE_TOKEN, null);
        commit(CURRENT_USER, null);
    },
    [ActionTypes.REFRESH_CURRENT_USER]: ({ commit, dispatch, state }: ActionParameters) => {
        return new FkApi().getCurrentUser().then((user) => {
            commit(CURRENT_USER, user);
            return user;
        });
    },
    [ActionTypes.UPLOAD_USER_PHOTO]: ({ commit, dispatch, state }: ActionParameters, payload: any) => {
        return new FkApi().uploadUserImage({ type: payload.type, file: payload.file }).then(() => {
            return dispatch(ActionTypes.REFRESH_CURRENT_USER);
        });
    },
    [ActionTypes.UPDATE_USER_PROFILE]: ({ commit, dispatch, state }: ActionParameters, payload: any) => {
        return new FkApi().updateUser(payload.user).then(() => {
            return dispatch(ActionTypes.REFRESH_CURRENT_USER);
        });
    },
};

const mutations = {
    [MutationTypes.INITIALIZE]: (state: UserState, token: string) => {
        Vue.set(state, "token", new TokenStorage().getToken());
    },
    [UPDATE_TOKEN]: (state: UserState, token: string) => {
        new TokenStorage().setToken(token);
        Vue.set(state, "token", token);
        if (token === null) {
            Vue.set(state, "user", null);
        }
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
