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
        path: '/projects/@/' + projectSlug + '/expeditions/@/' + expeditionSlug + '/inputs',
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
        path: '/inputs/' + id + '/geojson?descending=false',
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
        path: '/inputs/' + id,
        method: 'GET',
        unwrap: (r) => r
    };
}

export function getSourceSummary(id) {
    return {
        types: ActionTypes.API_SOURCE_SUMMARY_GET,
        path: '/inputs/' + id + '/summary',
        method: 'GET',
        unwrap: (r) => r
    };
}
