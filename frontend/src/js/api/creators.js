import _ from 'lodash';

import * as ActionTypes from '../actions/types';

export function getProject(projectSlug) {
    return {
        types: ActionTypes.API_PROJECT_GET,
        path: '/projects/@/' + projectSlug,
        method: 'GET',
        unwrap: (r) => r
    };
}

export function getProjectExpeditions(projectSlug) {
    return {
        types: ActionTypes.API_PROJECT_EXPEDITIONS_GET,
        path: '/projects/@/' + projectSlug + '/expeditions',
        method: 'GET',
        unwrap: (r) => r.expeditions
    };
}

export function getExpedition(projectSlug, expeditionSlug) {
    return {
        types: ActionTypes.API_EXPEDITION_GET,
        path: '/projects/@/' + projectSlug + '/expeditions/@/' + expeditionSlug,
        method: 'GET',
        unwrap: (r) => r
    };
}

export function getExpeditionGeoJson(projectSlug, expeditionSlug) {
    return {
        types: ActionTypes.API_EXPEDITION_GEOJSON_GET,
        path: '/projects/@/' + projectSlug + '/expeditions/@/' + expeditionSlug + '/geojson',
        method: 'GET',
        unwrap: (r) => r
    };
}

export function getExpeditionSummary(projectSlug, expeditionSlug) {
    return {
        types: ActionTypes.API_EXPEDITION_SUMMARY_GET,
        path: '/projects/@/' + projectSlug + '/expeditions/@/' + expeditionSlug + '/summary',
        method: 'GET',
        unwrap: (r) => r
    };
}

export function getExpeditionSources(projectSlug, expeditionSlug) {
    return {
        types: ActionTypes.API_EXPEDITION_SOURCES_GET,
        path: '/projects/@/' + projectSlug + '/expeditions/@/' + expeditionSlug + '/sources',
        method: 'GET',
        unwrap: (r) => r
    };
}

export function getFeatureGeoJson(id) {
    return {
        types: ActionTypes.API_FEATURE_GEOJSON_GET,
        path: '/features/' + id + '/geojson',
        method: 'GET',
        unwrap: (r) => r
    };
}

export function getNextExpeditionGeoJson(last) {
    return {
        types: ActionTypes.API_EXPEDITION_GEOJSON_GET,
        path: last.nextUrl,
        method: 'GET',
        unwrap: (r) => r
    };
}

export function getSourceGeoJson(id) {
    return {
        types: ActionTypes.API_SOURCE_GEOJSON_GET,
        path: '/sources/' + id + '/geojson?descending=false',
        method: 'GET',
        unwrap: (r) => r
    };
}

export function getNextSourceGeoJson(last) {
    return {
        types: ActionTypes.API_SOURCE_GEOJSON_GET,
        path: last.nextUrl,
        method: 'GET',
        unwrap: (r) => r
    };
}

export function getSource(id) {
    return {
        types: ActionTypes.API_SOURCE_GET,
        path: '/sources/' + id,
        method: 'GET',
        unwrap: (r) => r
    };
}

export function getSourceSummary(id) {
    return {
        types: ActionTypes.API_SOURCE_SUMMARY_GET,
        path: '/sources/' + id + '/summary',
        method: 'GET',
        unwrap: (r) => r
    };
}

function objectToQueryString(obj) {
    return _.reduce(obj, function (result, value, key) {
        if (!_.isNull(value) && !_.isUndefined(value)) {
            if (_.isArray(value)) {
                result += _.reduce(value, function (result1, value1) {
                    if (!_.isNull(value1) && !_.isUndefined(value1)) {
                        result1 += key + '=' + value1 + '&';
                        return result1
                    } else {
                        return result1;
                    }
                }, '')
            } else {
                result += key + '=' + value + '&';
            }
            return result;
        } else {
            return result
        }
    }, '').slice(0, -1);
};

export function getQuery(chart, criteria) {
    const queryString = objectToQueryString(Object.assign({}, chart, criteria));
    return {
        types: ActionTypes.API_SOURCE_QUERY_GET,
        path: '/sources/' + chart.sourceId + '/query?' + queryString,
        method: 'GET',
        unwrap: (r) => r,
        chart: chart,
        criteria: criteria
    };
}
