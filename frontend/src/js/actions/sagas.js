
import { delay } from 'redux-saga';
import {
    all, put, take, race, takeLatest, select
} from 'redux-saga/effects';

import _ from 'lodash';

import * as ActionTypes from './types';
import FkApi from '../api/calls';

import { changePlaybackMode, focusLocation, focusExpeditionTime } from './index';

import { PlaybackModes } from '../components/PlaybackControl';

import { FkGeoJSON } from '../common/geojson';

export function* getExpedition() {
    return yield select(state => new FkGeoJSON(state.visibleFeatures.geojson));
}

export function* walkGeoJson(summary) {
    const start = new Date(summary.startTime);
    const end = new Date(summary.endTime);

    const tickDuration = 100;
    const walkDurationInSeconds = 30;
    const walkLength = end - start;
    const walkLengthInSeconds = (walkLength / 1000);

    let walkSecondsPerTick = (walkLengthInSeconds / walkDurationInSeconds) / (1000 / tickDuration);

    console.log("Start", start);
    console.log("End", end);
    console.log("TickDuration", tickDuration, "WalkDuration", walkDurationInSeconds);
    console.log("SecondsPerTick", walkSecondsPerTick);

    while (true) {
        let time = new Date(start);

        while (time < end) {
            const expedition = yield getExpedition();
            const coordinates = expedition.getCoordinatesAtTime(time);

            yield put(focusExpeditionTime(time, coordinates, walkSecondsPerTick));

            const [ activityAction, focusFeatureAction, playbackAction ] = yield race([
                take(ActionTypes.USER_MAP_ACTIVITY),
                take(ActionTypes.FOCUS_FEATURE),
                take(ActionTypes.CHANGE_PLAYBACK_MODE),
                delay(tickDuration)
            ]);

            if (activityAction || focusFeatureAction) {
                yield put(changePlaybackMode(PlaybackModes.Pause));
                break;
            }

            if (playbackAction) {
                let nextAction = playbackAction;
                if (playbackAction.mode === PlaybackModes.Pause) {
                    nextAction = yield take(ActionTypes.CHANGE_PLAYBACK_MODE);
                }
                switch (nextAction.mode) {
                case PlaybackModes.Beginning: {
                    time = expedition.getFirst().time();
                    break;
                }
                case PlaybackModes.Rewind: {
                    walkSecondsPerTick = -5;
                    break;
                }
                case PlaybackModes.Play: {
                    walkSecondsPerTick = 5;
                    break;
                }
                case PlaybackModes.Forward: {
                    walkSecondsPerTick = 20;
                    break;
                }
                case PlaybackModes.End: {
                    time = expedition.getLast().time();
                    break;
                }
                default: {
                    break;
                }
                }
            }

            time.setTime(time.getTime() + (walkSecondsPerTick * 1000));
        }

        console.log("DONE", start, end, time);

        yield put(changePlaybackMode(PlaybackModes.Pause));

        yield take(ActionTypes.CHANGE_PLAYBACK_MODE);
    }
}

export function* manageMap() {
    yield takeLatest([ActionTypes.FOCUS_SOURCE], function* (focusSourceAction) {
        const source = focusSourceAction.source.source;
        yield loadAndWalkSource(source.id);
    });
}

export function* loadSources(ids) {
    const cache = {};
    const newIds = _.difference(ids, _.keys(cache));
    const sources = yield all(newIds.map(id => FkApi.getSource(id)));
    const summaries = yield all(newIds.map(id => FkApi.getSourceSummary(id)));
    console.log('Summaries', summaries, sources);
}

function getMapLocationFromQueryString() {
    const params = new URLSearchParams(window.location.search);
    const value = params.get('center') || "";
    const center = value.split(",").map(v => Number(v));
    if (center.length < 2) {
        return null;
    }
    return center;
}

function getDefaultMapLocation() {
    const fromQuery = getMapLocationFromQueryString();
    if (fromQuery != null) {
        return fromQuery;
    }

    return [-118.2688137, 34.0309388, 14];
}

export function* loadExpedition(projectSlug, expeditionSlug) {
    const [ expedition, sources ] = yield all([
        FkApi.getExpedition(projectSlug, expeditionSlug),
        FkApi.getExpeditionSources(projectSlug, expeditionSlug)
    ]);

    yield put(focusLocation(getDefaultMapLocation()));

    const sourceIds = _(sources.deviceSources).map('id').uniq().value();
    const detailedSources = yield loadSources(sourceIds);

    return [ expedition, sources, sources, detailedSources ];
}

function* loadGeojson(lastPage, call) {
    while (true) {
        lastPage = yield call(lastPage);
        if (lastPage.geo.features.length === 0) {
            break;
        }
    }
}

export function* loadAllGeojson(deviceId) {
    const [ pagedGeoJson, source ] = yield all([
        FkApi.getSourceGeoJson(deviceId),
        FkApi.getSource(deviceId),
    ]);

    yield loadGeojson(pagedGeoJson, (page) => FkApi.getNextSourceGeoJson(page));

    return [ source, yield select(state => state.visibleFeatures.geojson) ];
}

export function* loadAndWalkSource(deviceId) {
    const [ source, geojson ] = yield loadAllGeojson(deviceId);

    yield delay(1000);

    yield walkGeoJson(source, geojson.geo);
}

export function* loadActiveProject() {
    const projectSlug = 'www';

    const [ project, expeditions ] = yield all([
        FkApi.getProject(projectSlug),
        FkApi.getProjectExpeditions(projectSlug)
    ]);

    console.log('Project', project);

    if (expeditions.length > 0) {
        yield all([
            manageMap(),
            loadExpedition(projectSlug, expeditions[0].slug),
        ]);
        console.log("Done");
    }
}

export function* rootSaga() {
    yield all([
        loadActiveProject()
    ]);
}
