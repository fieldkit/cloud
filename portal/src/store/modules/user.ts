import Vue from "vue";
import * as ActionTypes from "../actions";
import * as MutationTypes from "../mutations";

import { ResumeAction, DiscourseLoginAction } from "@/store";

import { FKApi, TokenStorage, Services, LoginPayload, LoginResponse, CurrentUser } from "@/api";

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

const actions = (services: Services) => {
    return {
        [ActionTypes.INITIALIZE]: async ({ commit, dispatch, state }: ActionParameters) => {
            if (state.token) {
                await dispatch(ActionTypes.REFRESH_CURRENT_USER);
            }
        },
        [ActionTypes.LOGIN]: async ({ commit, dispatch, state }: ActionParameters, payload: LoginPayload) => {
            await services.api.login(payload.email, payload.password).then((token) => {
                commit(UPDATE_TOKEN, token);
                return dispatch(ActionTypes.REFRESH_CURRENT_USER).then(() => {
                    return dispatch(ActionTypes.AUTHENTICATED).then(() => {
                        return new LoginResponse(token);
                    });
                });
            });
        },
        [ActionTypes.LOGIN_RESUME]: async ({ commit, dispatch, state }: ActionParameters, payload: ResumeAction) => {
            await services.api.resume(payload.token).then((token) => {
                commit(UPDATE_TOKEN, token);
                return dispatch(ActionTypes.REFRESH_CURRENT_USER).then(() => {
                    return dispatch(ActionTypes.AUTHENTICATED).then(() => {
                        return new LoginResponse(token);
                    });
                });
            });
        },
        [ActionTypes.LOGIN_DISCOURSE]: async ({ commit, dispatch, state }: ActionParameters, payload: DiscourseLoginAction) => {
            await services.api
                .loginDiscourse(
                    payload.token,
                    payload.login?.email ?? null,
                    payload.login?.password ?? null,
                    payload.discourse.sso,
                    payload.discourse.sig
                )
                .then((response) => {
                    commit(UPDATE_TOKEN, response.token);
                    return dispatch(ActionTypes.REFRESH_CURRENT_USER).then(() => {
                        return dispatch(ActionTypes.AUTHENTICATED).then(() => {
                            // return new LoginResponse(response.token);
                        });
                    });
                });
        },
        [ActionTypes.LOGOUT]: async ({ commit }: { commit: any }) => {
            try {
                await services.api.logout();
            } finally {
                commit(UPDATE_TOKEN, null);
                commit(CURRENT_USER, null);
            }
        },
        [ActionTypes.REFRESH_CURRENT_USER]: ({ commit, dispatch, state }: ActionParameters) => {
            return services.api.getCurrentUser().then((user) => {
                commit(CURRENT_USER, user);
                return user;
            });
        },
        [ActionTypes.UPLOAD_USER_PHOTO]: ({ commit, dispatch, state }: ActionParameters, payload: any) => {
            return services.api.uploadUserImage({ type: payload.type, file: payload.file }).then(() => {
                return dispatch(ActionTypes.REFRESH_CURRENT_USER);
            });
        },
        [ActionTypes.UPDATE_USER_PROFILE]: ({ commit, dispatch, state }: ActionParameters, payload: any) => {
            return services.api.updateUser(payload.user).then(() => {
                return dispatch(ActionTypes.REFRESH_CURRENT_USER);
            });
        },
    };
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

export const user = (services: Services) => {
    const state = () => new UserState();

    return {
        namespaced: false,
        state,
        getters,
        actions: actions(services),
        mutations,
    };
};
