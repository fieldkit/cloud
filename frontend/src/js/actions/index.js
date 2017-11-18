import * as ActionTypes from './types';

export function changePlaybackMode(mode) {
    return {
        type: ActionTypes.CHANGE_PLAYBACK_MODE,
        mode: mode
    };
}

export function focusExpeditionTime(time) {
    return {
        type: ActionTypes.FOCUS_EXPEDITION_TIME,
        time: time
    };
}
