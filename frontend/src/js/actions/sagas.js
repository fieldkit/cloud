
import { delay } from 'redux-saga';
import {
    all, put, take, race, takeLatest
    // takeLatest, takeEvery, select, all, call
} from 'redux-saga/effects';

import _ from 'lodash';

import * as ActionTypes from './types';
import FkApi from '../api/calls';

import { changePlaybackMode, focusExpeditionTime } from './index';

import { PlaybackModes } from '../components/PlaybackControl';

import { FkGeoJSON } from '../common/geojson';

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

export function* loadExpeditionDetails() {
    let cache = {};

    yield takeLatest([ActionTypes.API_EXPEDITION_GEOJSON_GET.SUCCESS], function* (geojsonSuccess) {
        const geojson = geojsonSuccess.response.geo;
        const expedition = new FkGeoJSON(geojson);
        const ids = expedition.getUniqueInputIds().map(id => id.toString());
        const newIds = _.difference(ids, _.keys(cache));
        const queried = yield all(newIds.map(id => FkApi.getInput(id)));
        const indexed = _(queried).keyBy('id').value();
        cache = { ...cache, ...indexed }
        console.log(cache);
    });
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
        yield all([
            loadExpeditionDetails(),
            loadActiveExpedition(projectSlug, expeditions[0].slug)
        ])
        console.log("Done")
    }
}

export function* rootSaga() {
    yield all([
        loadActiveProject()
    ]);
}
