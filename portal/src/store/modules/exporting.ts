import _ from "lodash";
import Vue from "vue";
import Promise from "bluebird";
import * as ActionTypes from "../actions";

import { ExportDataAction, ExportParams } from "../typed-actions";
import { Bookmark } from "@/views/viz/viz";
import FKApi, { OnNoReject, UserExports, ExportStatus } from "@/api/api";

const EXPORT_START = "EXPORT_START";
const EXPORT_PROGRESS = "EXPORT_PROGRESS";
const EXPORT_DONE = "EXPORT_DONE";
const EXPORT_ERROR = "EXPORT_ERROR";
const EXPORT_CHECK = "EXPORT_CHECK";
const USER_EXPORTS = "USER_EXPORTS";

export class ExportingState {
    history: { [index: string]: ExportStatus };
}

const getters = {};

export interface CheckPayload {
    statusUrl: string;
}

type ActionParameters = { dispatch: any; commit: any; state: ExportingState };

function makeExportParams(bookmark: Bookmark): URLSearchParams {
    const range = bookmark.allTimeRange;
    const params = new URLSearchParams();
    params.append("stations", bookmark.allStations.join(","));
    params.append("sensors", bookmark.allSensors.join(","));
    params.append("start", range.start.toString());
    params.append("end", range.end.toString());
    return params;
}

const actions = {
    [EXPORT_CHECK]: async ({ dispatch, commit }: ActionParameters, payload: CheckPayload) => {
        return new FKApi().exportStatus(payload.statusUrl).then((de) => {
            console.log("exporting:status", de);
            if (!de.downloadUrl) {
                commit(EXPORT_PROGRESS, de);
                return Promise.delay(1000).then(() => {
                    return dispatch(EXPORT_CHECK, payload);
                });
            } else {
                commit(EXPORT_DONE, de);
                return dispatch(ActionTypes.NEED_EXPORTS);
            }
        });
    },
    [ActionTypes.NEED_EXPORTS]: async ({ dispatch, commit }: ActionParameters, payload: ExportDataAction) => {
        return new FKApi().getUserExports().then((exports) => {
            commit(USER_EXPORTS, exports);

            return Promise.all(
                exports.exports.map((de) => {
                    if (!de.completedAt) {
                        return dispatch(EXPORT_CHECK, { statusUrl: de.statusUrl });
                    }
                    return de;
                })
            );
        });
    },
    [ActionTypes.BEGIN_EXPORT]: async ({ dispatch, commit }: ActionParameters, payload: ExportDataAction) => {
        const params = makeExportParams(payload.bookmark);
        console.log("exporting:begin", payload, params);
        return new FKApi()
            .exportData(params, payload.params)
            .then((de) => {
                commit(EXPORT_START, de);
                return dispatch(EXPORT_CHECK, { statusUrl: de.statusUrl });
            })
            .catch((error) => {
                commit(EXPORT_ERROR, error);
            });
    },
};

const mutations = {
    [USER_EXPORTS]: (state: ExportingState, payload: UserExports) => {
        Vue.set(state, "history", _.fromPairs(payload.exports.map((de) => [de.token, de])));
    },
    [EXPORT_START]: (state: ExportingState, de: ExportStatus) => {
        Vue.set(state.history, de.token, de);
    },
    [EXPORT_PROGRESS]: (state: ExportingState, de: ExportStatus) => {
        Vue.set(state.history, de.token, de);
    },
    [EXPORT_DONE]: (state: ExportingState, de: ExportStatus) => {
        Vue.set(state.history, de.token, de);
    },
    [EXPORT_ERROR]: (state: ExportingState, error: any) => {
        console.log("exporting:error", error);
    },
};

const state = () => new ExportingState();

export const exporting = {
    namespaced: false,
    state,
    getters,
    actions,
    mutations,
};
