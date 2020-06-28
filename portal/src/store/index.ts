// import _ from "lodash";
import Vuex from "vuex";
import createLogger from "./logger";
import { clock } from "./modules/clock";
import { user } from "./modules/user";
import { stations } from "./modules/stations";
import { map } from "./modules/map";
import { progress } from "./modules/progress";
// import * as MutationTypes from "./mutations";
// import * as ActionTypes from "./actions";

function customizeLogger() {
    return createLogger({
        /*
        filter(mutation, stateBefore, stateAfter) {
            return true;
        },
        actionFilter(action, state) {
            return true;
        },
        transformer(state) {
            return state;
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

export default function() {
    return new Vuex.Store({
        plugins: [customizeLogger()],
        modules: {
            clock,
            user,
            stations,
            map,
            progress,
        },
        // This was causing a call stack error (_traverse)
        strict: process.env.NODE_ENV !== "production",
    });
}
