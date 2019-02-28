import _ from 'lodash';
import Promise from 'bluebird';

import * as ActionTypes from './types';
import { GeoRectSet, GeoRect } from '../common/geo';
import { CALL_WEB_API } from '../api/middleware';
import { getQuery } from '../api/creators';
import { FkApi } from '../api/calls';

export function changePlaybackMode(mode) {
    return {
        type: ActionTypes.CHANGE_PLAYBACK_MODE,
        mode: mode
    };
}

export function focusLocation(center) {
    return {
        type: ActionTypes.FOCUS_LOCATION,
        center: center,
        altitude: center.length === 3 ? center[2] : 10,
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

export function loadMapFeatures(criteria) {
    return (dispatch, getState) => {
        const { map } = getState();

        const desired = new GeoRect(criteria);
        const loaded = new GeoRectSet(map.loaded);
        if (loaded.contains(desired)) {
            return { };
        }

        const loading = desired.enlarge(2);
        loaded.add(loading);

        const api = new FkApi(dispatch);
        return api.queryMapFeatures({
            ne: loading.ne,
            sw: loading.sw,
        }).then(data => {
            const { sources } = getState();
            const sourceIds = _.union(_(data.spatial).map(s => s.sourceId).value(), _(data.temporal).map(s => s.sourceId).value());
            return Promise.all(sourceIds.map(id => {
                if (sources[id]) {
                    return Promise.resolve(sources[id]);
                }
                return api.getSource(id).then(source => {
                    return api.getSourceSummary(id);
                });
            }));
        });
    };
}

export function loadDevices() {
    return (dispatch, getState) => {
        const api = new FkApi(dispatch);
        return api.queryDevices().then(data => {
            return data;
        });
    };
}

export function loadDeviceFiles(deviceId) {
    return (dispatch, getState) => {
        const api = new FkApi(dispatch);
        return api.queryDeviceFilesData(deviceId).then(data => {
            return api.queryDeviceFilesLogs(deviceId).then(logs => {
                return { data, logs };
            });
        });
    };
}
