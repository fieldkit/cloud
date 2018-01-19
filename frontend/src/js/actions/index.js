import * as ActionTypes from './types';

export function changePlaybackMode(mode) {
    return {
        type: ActionTypes.CHANGE_PLAYBACK_MODE,
        mode: mode
    };
}

export function focusExpeditionTime(time, center) {
    return {
        type: ActionTypes.FOCUS_TIME,
        time: time,
        center: center
    };
}

export function notifyOfUserMapActivity() {
    return {
        type: ActionTypes.USER_MAP_ACTIVITY
    };
}

export function focusFeature(feature) {
    return {
        type: ActionTypes.FOCUS_FEATURE,
        feature: feature
    }
}
