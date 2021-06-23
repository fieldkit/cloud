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
import { Services } from "@/api";

export * from "@/api";

export * from "./modules/clock";
export * from "./modules/user";
export * from "./modules/stations";
export * from "./modules/exporting";
export * from "./modules/map";
export * from "./modules/progress";
export * from "./modules/layout";
export * from "./modules/global";

import * as MutationTypes from "./mutations";
import * as ActionTypes from "./actions";

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
        plugins: [customizeLogger()],
        modules: {
            clock: clock(services),
            exporting: exporting(services),
            user: user(services),
            stations: stations(services),
            map: map(services),
            progress: progress(services),
            layout: layout(services),
        },
        // This was causing a call stack error (_traverse)
        strict: process.env.NODE_ENV !== "production",
    });
}
