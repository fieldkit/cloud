function expandApiTypes(base) {
    return {
        'START': base + "_START",
        'SUCCESS': base + "_SUCCESS",
        'FAILURE': base + "_FAILURE",
    };
}

export const API_EXPEDITION_GEOJSON_GET = expandApiTypes("API_EXPEDITION_GEOJSON_GET");
export const API_EXPEDITION_SOURCES_GET = expandApiTypes("API_EXPEDITION_SOURCES_GET");
export const API_EXPEDITION_SUMMARY_GET = expandApiTypes("API_EXPEDITION_SUMMARY_GET");
export const API_PROJECT_EXPEDITIONS_GET = expandApiTypes("API_PROJECT_EXPEDITIONS_GET");
export const API_EXPEDITION_GET = expandApiTypes("API_EXPEDITION_GET");
export const API_PROJECT_GET = expandApiTypes("API_PROJECT_GET");
export const API_SOURCE_GET = expandApiTypes("API_SOURCE_GET");
export const API_SOURCE_SUMMARY_GET = expandApiTypes("API_SOURCE_SUMMARY_GET");
export const API_CLUSTER_GEOMETRY_GET = expandApiTypes("API_CLUSTER_GEOMETRY_GET");
export const API_FEATURE_GEOJSON_GET = expandApiTypes("API_FEATURE_GEOJSON_GET");
export const API_SOURCE_GEOJSON_GET = expandApiTypes("API_SOURCE_GEOJSON_GET");
export const API_SOURCE_QUERY_GET = expandApiTypes("API_SOURCE_QUERY_GET");

export const CHANGE_PLAYBACK_MODE = 'CHANGE_PLAYBACK_MODE';

export const FOCUS_SOURCE = 'FOCUS_SOURCE';
export const FOCUS_FEATURE = 'FOCUS_FEATURE';
export const FOCUS_TIME = 'FOCUS_TIME';
export const FOCUS_LOCATION = 'FOCUS_LOCATION';

export const USER_MAP_ACTIVITY = 'USER_MAP_ACTIVITY';

export const CHART_SOURCE_LOAD = 'CHART_SOURCE_LOAD';
export const CHART_DATA_LOAD = 'CHART_DATA_LOAD';
export const CHART_CRITERIA_CHANGE = 'CHART_CRITERIA_CHANGE';

export function isVerboseActionType(type) {
    return type === FOCUS_TIME;
}
