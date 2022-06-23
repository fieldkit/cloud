import _ from "lodash";
import Vuex from "vuex";
import createLogger from "./logger";
import { clock } from "./modules/clock";
import { user } from "./modules/user";
import { stations } from "./modules/stations";
import { map } from "./modules/map";
import { progress } from "./modules/progress";
import { layout } from "./modules/layout";
import { exporting } from "./modules/exporting";
import { notifications } from "./modules/notifications";
import { Services } from "@/api";

export * from "@/api";

export * from "./modules/clock";
export * from "./modules/user";
export * from "./modules/stations";
export * from "./modules/exporting";
export * from "./modules/map";
export * from "./modules/progress";
export * from "./modules/layout";
export * from "./modules/notifications";
export * from "./modules/global";
export * from "./modules/notes";

import * as MutationTypes from "./mutations";
import * as ActionTypes from "./actions";
import { notes } from "@/store/modules/notes";

export { MutationTypes, ActionTypes };

export * from "./typed-actions";
export * from "./map-types";

function customizeLogger() {
    return createLogger({
        transformer(state) {
            if (state.user.token) {
                const copy = _.cloneDeep(state);
                copy.user.token = "YES";
                return copy;
            }
            return state;
        },
        /*
        filter(mutation, stateBefore, stateAfter) {
            return true;
        },
        actionFilter(action, state) {
            return true;
        },
        mutationTransformer(mutation) {
            return mutation;
        },
        actionTransformer(action) {
            return action;
        },
        logActions: true,
        logMutations: true,
		*/
    });
}

export default function(services: Services) {
    return new Vuex.Store({
        plugins: [
            /*customizeLogger()*/
        ],
        modules: {
            clock: clock(services),
            exporting: exporting(services),
            notifications: notifications(services),
            user: user(services),
            stations: stations(services),
            map: map(services),
            progress: progress(services),
            layout: layout(services),
            notes: notes(services),
        },
        // This was causing a call stack error (_traverse)
        strict: process.env.NODE_ENV !== "production",
    });
}
