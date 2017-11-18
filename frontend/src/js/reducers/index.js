import { combineReducers } from 'redux';

import * as ActionTypes from '../actions/types';

function activeExpedition(state = { project: null, expedition: null }, action) {
    let nextState = state;
    switch (action.type) {
    case ActionTypes.API_PROJECT_GET.SUCCESS:
        return Object.assign({ }, state, { project: action.response });
    case ActionTypes.API_EXPEDITION_GET.SUCCESS:
        return Object.assign({ }, state, { expedition: action.response });
    default:
        return nextState;
    }
}

function visibleFeatures(state = { }, action) {
    let nextState = state;
    switch (action.type) {
    case ActionTypes.API_EXPEDITION_GEOJSON_GET.SUCCESS: {
        let newGeojson = Object.assign({}, action.response.geo);
        if (state.geojson) {
            newGeojson.features = [ ...state.geojson.features, ...newGeojson.features ]
        }
        return Object.assign({ }, state, { geojson: newGeojson });
    }
    case ActionTypes.FOCUS_FEATURE:
        return Object.assign({ }, state, { focus: { feature: action.feature } });
    case ActionTypes.FOCUS_TIME:
        return Object.assign({ }, state, { focus: { time: action.time, center: action.center } });
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
    visibleFeatures,
    playbackMode,
});
