import { combineReducers } from 'redux';

import * as ActionTypes from '../actions/types';

function activeExpedition(state = { project: null, expedition: null }, action) {
    let nextState = state;
    switch (action.type) {
    case ActionTypes.API_PROJECT_GET.SUCCESS:
        return Object.assign({ }, state, { project: action.response });
    case ActionTypes.API_EXPEDITION_GET.SUCCESS:
        return Object.assign({ }, state, { expedition: action.expedition });
    default:
        return nextState;
    }
}

function visibleGeoJson(state = { }, action) {
    let nextState = state;
    switch (action.type) {
    case ActionTypes.API_EXPEDITION_GEOJSON_GET.SUCCESS:
        return Object.assign({ }, state, action.response);
    case ActionTypes.FEATURES_FOCUS:
        return Object.assign({ }, state, { focus: action.feature });
    default:
        return nextState;
    }
}

function playbackMode(state = { }, action) {
    let nextState = state;
    switch (action.type) {
    case ActionTypes.CHANGE_PLAYBACK_MODE:
        return action.mode; // TODO: Treat this object as immutable. Needs better protection.
    default:
        return nextState;
    }
}

export default combineReducers({
    activeExpedition,
    visibleGeoJson,
    playbackMode,
});
