import _ from 'lodash';

import { combineReducers } from 'redux';

import * as ActionTypes from '../actions/types';

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
    focus: { features: [] },
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

function visibleFeatures(state = visibleFeaturesInitialState, action) {
    switch (action.type) {
    case ActionTypes.API_EXPEDITION_GEOJSON_GET.SUCCESS: {
        return mergeFeatures(state, action);
    }
    case ActionTypes.API_SOURCE_GEOJSON_GET.SUCCESS: {
        return mergeFeatures(state, action);
    }
    case ActionTypes.FOCUS_FEATURE:
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
        const summary = action.response;
        const container = state.sources[summary.id] || { };
        const nextState = Object.assign({}, state);
        nextState.sources[summary.id] = {...container, ...{ summary: summary } };
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
    charts: [],
    criteria: {
        startTime: Math.trunc(start.getTime() / 1000),
        endTime: Math.trunc(now.getTime() / 1000),
    }
};

function getOrCreateChart(state, chart) {
    for (let c of state.charts) {
        if (c.id === chart.id) {
            return c;
        }
    }

    const newChart = Object.assign({}, chart);
    state.charts.push(newChart);
    return newChart;
}

function chartData(state = initialChartDataState, action) {
    switch (action.type) {
    case ActionTypes.CHART_CRITERIA_CHANGE: {
        return { ...state, ...{ criteria: action.criteria } };
    }
    case ActionTypes.CHART_DATA_LOAD: {
        const nextState = Object.assign({}, state);
        const chart = getOrCreateChart(nextState, action.chart);
        Object.assign(chart, action.chart);
        return nextState;
    }
    case ActionTypes.API_SOURCE_QUERY_GET.START: {
        const nextState = Object.assign({}, state);
        const chart = getOrCreateChart(nextState, action.chart);
        chart.loading = true;
        chart.query = {};
        return nextState;
    }
    case ActionTypes.API_SOURCE_QUERY_GET.SUCCESS: {
        const nextState = Object.assign({}, state);
        const chart = getOrCreateChart(nextState, action.chart);
        chart.loading = false;
        chart.query = action.response;
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
    chartData
});
