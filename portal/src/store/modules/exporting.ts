// import _ from "lodash";
import Promise from "bluebird";
import * as ActionTypes from "../actions";

import { ExportDataAction, ExportParams } from "../typed-actions";
import { Bookmark } from "@/views/viz/viz";
import FKApi, { OnNoReject } from "@/api/api";

const EXPORT_START = "EXPORT_START";
const EXPORT_PROGRESS = "EXPORT_PROGRESS";
const EXPORT_DONE = "EXPORT_DONE";
const EXPORT_ERROR = "EXPORT_ERROR";
const EXPORT_CHECK = "EXPORT_CHECK";

export class ActiveExport {
    constructor(public readonly bookmark: Bookmark, public readonly params: ExportParams) {}
}

export class ExportingState {
    active: { [index: string]: ActiveExport };
}

const getters = {};

type ActionParameters = { dispatch: any; commit: any; state: ExportingState };

export interface CheckPayload {
    statusUrl: string;
}

const actions = {
    [EXPORT_CHECK]: async ({ dispatch, commit }: ActionParameters, payload: CheckPayload) => {
        return new FKApi().exportStatus(payload.statusUrl).then((de) => {
            console.log(de);
            if (!de.url) {
                return Promise.delay(1000).then(() => {
                    return dispatch(EXPORT_CHECK, payload);
                });
            }
        });
    },
    [ActionTypes.BEGIN_EXPORT]: async ({ dispatch, commit }: ActionParameters, payload: ExportDataAction) => {
        commit(EXPORT_START, {});

        console.log("hello", payload);

        // TODO: TEMPORARY
        const params = new URLSearchParams();
        params.append("stations", "159");
        params.append("sensors", "2");
        return new FKApi().exportCsv(params).then((de) => {
            return dispatch("EXPORT_CHECK", { statusUrl: de.location });
        });
    },
};

const mutations = {
    [EXPORT_START]: () => {
        //
    },
    [EXPORT_PROGRESS]: () => {
        //
    },
    [EXPORT_DONE]: () => {
        //
    },
    [EXPORT_ERROR]: () => {
        //
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
