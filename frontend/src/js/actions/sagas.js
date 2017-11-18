
import { delay } from 'redux-saga';
import {
    all, put, take, race
    // takeLatest, takeEvery, select, all, call
} from 'redux-saga/effects';

import _ from 'lodash';

import * as ActionTypes from './types';
import FkApi from '../api/calls';

import { changePlaybackMode, focusExpeditionTime } from './index';

import { PlaybackModes } from '../components/PlaybackControl';

class FkGeoJSONFeature {
    constructor(feature) {
        this.feature = feature;
    }

    time() {
        return new Date(this.feature.properties["timestamp"]);
    }

    coords() {
        return this.feature.geometry.coordinates;
    }
}

// Love this name. I'm open to another. This basically means give each pair of
// [1, 2, 3, 4] as [1, 2] [2, 3] [3, 4] as you would imagine though also include
// [null, 1] and [4, null]
// This also means you'll get called with null, null for an empty array.
function filterPairwiseIncludingSentinels(arr, func, tx) {
    const sentinel = {};
    const heh =  [ sentinel ].concat(arr).concat([ sentinel ]);
    for (let i = 0; i < heh.length - 2; i++){
        let a = heh[i] === sentinel ? null : heh[i];
        let b = heh[i + 1] === sentinel ? null : heh[i + 1];
        if (func(a, b)) {
            return [ a, b ];
        }
    }
    return null;
}

export function map(val, domain1, domain2, range1, range2) {
    if (domain1 === domain2) {
        return range1;
    }
    return (val - domain1) / (domain2 - domain1) * (range2 - range1) + range1;
}

/**
 * Class for handling GeoJSON data specifically from the fk site. This basically
 * means timestamped GeoJSON data.
 */
class FkGeoJSON {
    constructor(geojson) {
        // Usually sorted server-side, though this should be fine.
        this.features = _(geojson.features).map(f => new FkGeoJSONFeature(f)).sortBy(f => f.time()).value();
    }

    getFirst() {
        return this.features[0];
    }

    getLast() {
        return this.features[this.features.length - 1];
    }

    getFeaturesSurroundingTime(time) {
        return filterPairwiseIncludingSentinels(this.features, (a, b) => {
            const beforeTime = a != null ? a.time() : 0;
            const afterTime = b != null ? b.time() : Number.MAX_SAFE_INTEGER;
            return time >= beforeTime && time < afterTime;
        });
    }

    getCoordinatesAtTime(time) {
        const pair = _(this.getFeaturesSurroundingTime(time)).value();

        const beforeCoordinates = pair[0].coords();
        const afterCoordinates = pair[1].coords();
        const beforeTime = pair[0].time().getTime();
        const afterTime = pair[1].time().getTime();

        const latitude = map(time.getTime(), beforeTime, afterTime, beforeCoordinates[1], afterCoordinates[1])
        const longitude = map(time.getTime(), beforeTime, afterTime, beforeCoordinates[0], afterCoordinates[0])

        return [ longitude, latitude ];
    }
}

export function* walkExpedition(geojson) {
    const expedition = new FkGeoJSON(geojson);

    const start = expedition.getFirst().time();
    const end = expedition.getLast().time();

    let expeditionMinutesPerTick = 5;
    let time = start;

    while (true) {
        while (time < end) {
            const coordinates = expedition.getCoordinatesAtTime(time);

            yield put(focusExpeditionTime(time, coordinates));

            const [ activityAction, playbackAction, _ ] = yield race([
                take(ActionTypes.USER_MAP_ACTIVITY),
                take(ActionTypes.CHANGE_PLAYBACK_MODE),
                delay(100)
            ]);

            if (activityAction) {
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
                    expeditionMinutesPerTick = -5;
                    break;
                }
                case PlaybackModes.Play: {
                    expeditionMinutesPerTick = 5;
                    break;
                }
                case PlaybackModes.Forward: {
                    expeditionMinutesPerTick = 20;
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

            time.setTime(time.getTime() + (expeditionMinutesPerTick * 60 * 1000));
        }

        console.log("DONE", start, end, time);

        time = expedition.getFirst().time();

        yield put(changePlaybackMode(PlaybackModes.Pause));

        yield take(ActionTypes.CHANGE_PLAYBACK_MODE);
    }
}

export function* refreshSaga(pagedGeojson) {
    while (true) {
        yield delay(10000);

        pagedGeojson = yield FkApi.getNextExpeditionGeoJson(pagedGeojson);
    }
}

export function* loadActiveExpedition(projectSlug, expeditionSlug) {
    const [ expedition, pagedGeojson ] = yield all([
        FkApi.getExpedition(projectSlug, expeditionSlug),
        FkApi.getExpeditionGeoJson(projectSlug, expeditionSlug)
    ]);

    console.log(expedition, pagedGeojson);

    if (pagedGeojson.geo.features.length > 0) {
        yield delay(1000)

        yield all([ walkExpedition(pagedGeojson.geo), refreshSaga(pagedGeojson) ]);
    }
}

export function* loadActiveProject() {
    const projectSlug = 'www';

    const [ project, expeditions ] = yield all([
        FkApi.getProject(projectSlug),
        FkApi.getProjectExpeditions(projectSlug)
    ]);

    console.log(project);

    if (expeditions.length > 0) {
        yield loadActiveExpedition(projectSlug, expeditions[0].slug);
    }
}

export function* rootSaga() {
    yield all([
        loadActiveProject()
    ]);
}
