"use strict";

/**
 * Get the first item that pass the test
 * by second argument function
 *
 * @param {Array} list
 * @param {Function} f
 * @return {*}
 */
function find(list, f) {
    return list.filter(f)[0];
}

/**
 * Deep copy the given object considering circular structure.
 * This function caches all nested objects and its copies.
 * If it detects circular structure, use cached copy to avoid infinite loop.
 *
 * @param {*} obj
 * @param {Array<Object>} cache
 * @return {*}
 */
function deepCopy(obj, cache) {
    if (cache === void 0) cache = [];

    // just return if obj is immutable value
    if (obj === null || typeof obj !== "object") {
        return obj;
    }

    // if obj is hit, it is in circular structure
    const hit = find(cache, function(c) {
        return c.original === obj;
    });
    if (hit) {
        return hit.copy;
    }

    const copy = Array.isArray(obj) ? [] : {};
    // put the copy into cache at first
    // because we want to refer it in recursive deepCopy
    cache.push({
        original: obj,
        copy: copy,
    });

    Object.keys(obj).forEach(function(key) {
        copy[key] = deepCopy(obj[key], cache);
    });

    return copy;
}

function startMessage(logger, message, collapsed) {
    const startMessage = collapsed ? logger.groupCollapsed : logger.group;

    // render
    try {
        startMessage.call(logger, message);
    } catch (e) {
        logger.log(message);
    }
}

function endMessage(logger) {
    try {
        logger.groupEnd();
    } catch (e) {
        logger.log("store: end");
    }
}

// Credits: borrowed code from fcomb/redux-logger

function createLogger(ref) {
    if (ref === void 0) ref = {};
    let collapsed = ref.collapsed;
    if (collapsed === void 0) collapsed = true;
    let filter = ref.filter;
    if (filter === void 0)
        filter = function(_mutation, _stateBefore, _stateAfter) {
            return true;
        };
    let transformer = ref.transformer;
    if (transformer === void 0)
        transformer = function(state) {
            return state;
        };
    let mutationTransformer = ref.mutationTransformer;
    if (mutationTransformer === void 0)
        mutationTransformer = function(mut) {
            return mut;
        };
    let actionFilter = ref.actionFilter;
    if (actionFilter === void 0)
        actionFilter = function(action, state) {
            return true;
        };
    let actionTransformer = ref.actionTransformer;
    if (actionTransformer === void 0)
        actionTransformer = function(act) {
            return act;
        };
    let logMutations = ref.logMutations;
    if (logMutations === void 0) logMutations = true;
    let logActions = ref.logActions;
    if (logActions === void 0) logActions = true;
    let logger = ref.logger;
    if (logger === void 0) logger = console;

    return function(store) {
        let prevState = deepCopy(store.state);

        if (typeof logger === "undefined") {
            return;
        }

        if (logMutations) {
            store.subscribe(function(mutation, state) {
                const nextState = deepCopy(state);

                if (filter(mutation, prevState, nextState)) {
                    const formattedMutation = mutationTransformer(mutation);
                    const message = "store: mutation " + mutation.type;

                    startMessage(logger, message, collapsed);
                    logger.log("store: prev state", transformer(prevState));
                    logger.log("store: mutation", formattedMutation);
                    logger.log("store: next state", transformer(nextState));
                    endMessage(logger);
                }

                prevState = nextState;
            });
        }

        if (logActions) {
            store.subscribeAction(function(action, state) {
                if (actionFilter(action, state)) {
                    const formattedAction = actionTransformer(action);
                    const message = "store: action " + action.type;

                    startMessage(logger, message, collapsed);
                    logger.log("store: action", formattedAction);
                    endMessage(logger);
                }
            });
        }
    };
}

export default createLogger;
