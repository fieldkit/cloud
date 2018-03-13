import * as ActionTypes from './types';
import { CALL_WEB_API } from '../api/middleware';
import { getQuery } from '../api/creators';

export function changePlaybackMode(mode) {
    return {
        type: ActionTypes.CHANGE_PLAYBACK_MODE,
        mode: mode
    };
}

export function focusLocation(center, altitude) {
    return {
        type: ActionTypes.FOCUS_LOCATION,
        center: center,
        altitude: altitude
    };
}

export function focusExpeditionTime(time, center, expeditionSecondsPerTick) {
    return {
        type: ActionTypes.FOCUS_TIME,
        time: time,
        center: center,
        expeditionSecondsPerTick: expeditionSecondsPerTick
    };
}

export function notifyOfUserMapActivity() {
    return {
        type: ActionTypes.USER_MAP_ACTIVITY
    };
}

export function focusSource(source) {
    return {
        type: ActionTypes.FOCUS_SOURCE,
        source: source
    };
}

export function focusFeature(feature) {
    return {
        type: ActionTypes.FOCUS_FEATURE,
        feature: feature
    };
}

export function loadChartData(chartDef) {
    return (dispatch, getState) => {
        const criteria = getState().chartData.criteria;

        dispatch({
            [CALL_WEB_API]: getQuery(chartDef, criteria)
        });
    };
}

export function changeCriteria(criteria) {
    return (dispatch, getState) => {
        const queries = getState().chartData.queries;

        dispatch({
            type: ActionTypes.CHART_CRITERIA_CHANGE,
            criteria: criteria
        });

        for (let id of Object.keys(queries)) {
            const query = queries[id];
            dispatch({
                [CALL_WEB_API]: getQuery(query.chartDef, criteria)
            });
        }
    };
}
