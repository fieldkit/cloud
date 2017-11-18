import { delay } from 'redux-saga';
import { put, take, takeLatest, takeEvery, select, all, race, call } from 'redux-saga/effects';

import * as ActionTypes from './types';
import FkApi from '../api/calls';

export function* refreshSaga() {
    while (true) {
        yield delay(1000);
    }
}

export function* loadActiveExpedition(projectSlug, expeditionSlug) {
    const [ expedition, geojson ] = yield all([
        FkApi.getExpedition(projectSlug, expeditionSlug),
        FkApi.getExpeditionGeoJson(projectSlug, expeditionSlug)
    ]);

    if (geojson.features.length > 0) {
        yield delay(2000)

        yield put({
            type: ActionTypes.FOCUS_FEATURE,
            feature: geojson.features[0]
        })
    }
}

export function* loadActiveProject() {
    const projectSlug = 'www';

    const [ project, expeditions ] = yield all([
        FkApi.getProject(projectSlug),
        FkApi.getProjectExpeditions(projectSlug)
    ]);

    if (expeditions.length > 0) {
        yield loadActiveExpedition(projectSlug, expeditions[0].slug);
    }
}

export function* rootSaga() {
    yield all([
        loadActiveProject(),
        refreshSaga()
    ]);
}
