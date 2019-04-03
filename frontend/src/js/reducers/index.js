import _ from 'lodash';

import { combineReducers } from 'redux';

import * as ActionTypes from '../actions/types';

import { DataColors } from '../common/colors';
import { FkGeoJSON } from '../common/geojson';

function activeExpedition(state = { project: null, expedition: null }, action) {
    switch (action.type) {
    case ActionTypes.API_PROJECT_GET.SUCCESS:
        return Object.assign({ }, state, { project: action.response });
    case ActionTypes.API_EXPEDITION_GET.SUCCESS:
        return Object.assign({ }, state, { expedition: action.response });
    default:
        return state;
    }
}

const visibleFeaturesInitialState = {
    focus: {
        center: null,
    },
    sources: { },
    geojson: { features: []}
};

function mergeFeatures(state, action) {
    let newGeojson = Object.assign({}, action.response.geo);
    if (state.geojson) {
        newGeojson.features = [ ...state.geojson.features, ...newGeojson.features ];
    }
    return Object.assign({ }, state, {
        geojson: newGeojson
    });
}

function createSummary(summary) {
    _.each(summary.readings, reading => {
        reading.chartDef = {
            id: ["reading", summary.id, reading.name].join("-"),
            sourceId: summary.id,
            keys: [reading.name],
        };
    });
    return summary;
}

function visibleFeatures(state = visibleFeaturesInitialState, action) {
    switch (action.type) {
    case ActionTypes.API_EXPEDITION_GEOJSON_GET.SUCCESS: {
        return mergeFeatures(state, action);
    }
    case ActionTypes.API_SOURCE_GEOJSON_GET.SUCCESS: {
        return mergeFeatures(state, action);
    }
    case ActionTypes.FOCUS_FEATURE: {
        return Object.assign({ }, state, {
            focus: {
                expeditionSecondsPerTick: 0,
                feature: action.feature,
                time: null,
                center: action.feature.geometry.coordinates,
                altitude: 0,
                features: []
            }
        });
    }
    case ActionTypes.API_MAP_FEATURES_GET.SUCCESS: {
        const nextState = Object.assign({ }, state);
        const geometries = action.response.geometries;
        geometries.forEach(g => {
            nextState.sources[g.sourceId] = nextState.sources[g.sourceId] || {};
            nextState.sources[g.sourceId].geometries = nextState.sources[g.sourceId].geometries || {};
            nextState.sources[g.sourceId].geometries[g.id] = g;
        });
        return nextState;
    }
    case ActionTypes.API_FEATURE_GEOJSON_GET.SUCCESS: {
        const feature = action.response.geo.features[0];
        const nextState = Object.assign({}, state);
        _.each(nextState.sources, (container, id) => {
            if (container.source.lastFeatureId === feature.properties.id) {
                container.lastFeature = feature;
            }
        });
        return nextState;
    }
    case ActionTypes.API_SOURCE_SUMMARY_GET.SUCCESS: {
        const summary = createSummary(action.response);
        const container = state.sources[summary.id] || { };
        const nextState = Object.assign({}, state);
        nextState.sources[summary.id] = {...container, ...{ summary: summary } };
        return nextState;
    }
    case ActionTypes.API_CLUSTER_GEOMETRY_GET.SUCCESS: {
        const geometry = action.response;
        const nextState = Object.assign({}, state);
        nextState.sources[geometry.sourceId].geometries[geometry.id] = geometry;
        return nextState;
    }
    case ActionTypes.API_SOURCE_GET.SUCCESS: {
        const source = action.response;
        const container = state.sources[source.id] || { };
        const nextState = Object.assign({}, state);
        nextState.sources[source.id] = {...container, ...{ source: source } };
        return nextState;
    }
    case ActionTypes.FOCUS_LOCATION: {
        return Object.assign({ }, state, {
            focus: {
                expeditionSecondsPerTick: 0,
                time: null,
                center: action.center,
                altitude: action.altitude,
                features: []
            }
        });
    }
    case ActionTypes.FOCUS_TIME:
        const expeditionSecondsPerTick = action.expeditionSecondsPerTick;
        const expedition = new FkGeoJSON(state.geojson);
        const features = expedition.getFeaturesWithinTime(action.time, 1 * 60 * 1000);

        return Object.assign({ }, state, {
            focus: {
                expeditionSecondsPerTick: expeditionSecondsPerTick,
                time: action.time,
                center: action.center,
                altitude: 0,
                features: features
            }
        });
    default:
        return state;
    }
}

function playbackMode(state = { }, action) {
    switch (action.type) {
    case ActionTypes.CHANGE_PLAYBACK_MODE:
        return action.mode; // TODO: Treat this object as immutable. Needs better protection.
    default:
        return state;
    }
}

const now = new Date();
const start = new Date();
start.setDate(start.getDate() - 1);

const initialChartDataState = {
    queries: {},
    criteria: {
        startTime: Math.trunc(start.getTime() / 1000),
        endTime: Math.trunc(now.getTime() / 1000),
    }
};

function chartData(state = initialChartDataState, action) {
    switch (action.type) {
    case ActionTypes.CHART_CRITERIA_CHANGE: {
        return { ...state, ...{ criteria: action.criteria } };
    }
    case ActionTypes.API_SOURCE_QUERY_GET.START: {
        const nextState = Object.assign({}, state);
        const pending = nextState.queries[action.chartDef.id] || {};
        pending.chartDef = action.chartDef;
        pending.loading = true;
        pending.query = {};
        nextState.queries[action.chartDef.id] = pending;
        return nextState;
    }
    case ActionTypes.API_SOURCE_QUERY_GET.SUCCESS: {
        const nextState = Object.assign({}, state);
        const loaded = nextState.queries[action.chartDef.id];
        loaded.chartDef = action.chartDef;
        loaded.loading = false;
        loaded.query = action.response;
        nextState.queries[action.chartDef.id] = loaded;
        return nextState;
    }
    default:
        return state;
    }
}

const initialMapState = {
    loaded: [],
};

function map(state = initialMapState, action) {
    switch (action.type) {
    case ActionTypes.API_MAP_FEATURES_GET.START: {
        const nextState = Object.assign({}, state);
        nextState.loaded.push(action.criteria);
        return nextState;
    }
    case ActionTypes.API_MAP_FEATURES_GET.SUCCESS: {
        return state;
    }
    default:
        return state;
    }
}

const initialSourcesState = {
    sources: {},
    summaries: {},
    geometries: {}
};

function sources(state = initialSourcesState, action) {
    switch (action.type) {
    case ActionTypes.API_SOURCE_GET.SUCCESS: {
        const nextState = Object.assign({}, state);
        const source = action.response;
        nextState[source.id] = source;
        return nextState;
    }
    default:
        return state;
    }
}

const LosAngelesCenter = [-118.26928432026534, 34.03118429949713];
const initialFilesState = {
    devices: [],
    filesByDevice: { },
    center: LosAngelesCenter
};

function files(state = initialFilesState, action) {
    switch (action.type) {
    case ActionTypes.FOCUS_LOCATION: {
        return { ...state, ...{ center: action.center } };
    }
    case ActionTypes.API_LOAD_DEVICES.SUCCESS: {
        const colors = new DataColors();
        const files = action.response;
        for (let i = 0; i < files.devices.length; ++i) {
            files.devices[i].color = colors.get();
        }
        console.log(files);
        return { ...state, ...files };
    }
    case ActionTypes.API_LOAD_DEVICE_LOGS_FILES.START: {
        if (action.criteria.page === 0) {
            const { deviceId } = action.criteria;
            if (_.isObject(state.filesByDevice[deviceId])) {
                const nextState = _.cloneDeep(state);
                nextState.filesByDevice[deviceId].logs = [];
                return nextState;
            }
        }
        return state;
    }
    case ActionTypes.API_LOAD_DEVICE_DATA_FILES.START: {
        if (action.criteria.page === 0) {
            const { deviceId } = action.criteria;
            if (_.isObject(state.filesByDevice[deviceId])) {
                const nextState = _.cloneDeep(state);
                nextState.filesByDevice[deviceId].data = [];
                return nextState;
            }
            return state;
        }
        return state;
    }
    case ActionTypes.API_LOAD_DEVICE_LOGS_FILES.SUCCESS: {
        const nextState = _.cloneDeep(state);
        const { deviceId } = action.criteria;

        let files = nextState.filesByDevice[deviceId];
        if (!_.isObject(files)) {
            nextState.filesByDevice[deviceId] = files = { logs: [], data: [], all: [] };
        }

        files.logs = [ ...files.logs, ...action.response.files ];
        files.all = _([ ...files.data, ...files.logs ]).sortBy(i => i.time).reverse().value();

        return nextState;
    }
    case ActionTypes.API_LOAD_DEVICE_DATA_FILES.SUCCESS: {
        const nextState = _.cloneDeep(state);
        const { deviceId } = action.criteria;

        let files = nextState.filesByDevice[deviceId];
        if (!_.isObject(files)) {
            nextState.filesByDevice[deviceId] = files = { logs: [], data: [], all: [] };
        }

        files.data = [ ...files.data, ...action.response.files ];
        files.all = _([ ...files.data, ...files.logs ]).sortBy(i => i.time).reverse().value();

        return nextState;
    }
    default:
        return state;
    }
}

export default combineReducers({
    activeExpedition,
    visibleFeatures,
    playbackMode,
    chartData,
    map,
    sources,
    files,
});
