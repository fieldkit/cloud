import { delay } from 'redux-saga'
import { put, take, takeLatest, takeEvery, select, all, race, call } from 'redux-saga/effects'
import I from 'immutable'

import * as ActionTypes from './types'
import FkApi from '../api/calls'

export function* refreshSaga() {
    while (true) {
        yield delay(1000);
    }
}

function mapBySlug(collection) {
    const map = {}
    collection.forEach(e => map[e.slug] = e)
    return map
}

function* loadExpeditionDocs(projectSlug, expeditionSlug) {
    const docs = yield FkApi.getDocs(projectSlug, expeditionSlug)
    const map = {}
    docs.forEach((d, i) => {
        d.data.type = 'Feature';
        d.data.id = d.id
        d.data.date = d.timestamp * 1000
        if (!d.data.geometry) {
            d.data.geometry = d.location
        }
        map[d.id] = d.data
    })

    return I.fromJS(map).toList()
        .sort((d1, d2) => {
            return d1.get('date') - d2.get('date')
        })
        .toList()
}

function* load() {
    const projectSlug = 'www'
    const expeditions = yield FkApi.getProjectExpeditions(projectSlug)
    const project = yield FkApi.getProject(projectSlug)
    if (!project) {
        return
    }

    yield put({
        type: ActionTypes.SET_PROJECT_NAME,
        name: project.name
    })

    const expeditionsMap = mapBySlug(expeditions);

    for (const expedition of expeditions) {
        yield put({
            type: ActionTypes.REQUEST_EXPEDITION,
            id: expedition.slug
        })

        const docs = yield loadExpeditionDocs(projectSlug, expedition.slug)
        if (docs.size > 0) {
            const startDate = docs.get(0).get('date')
            const endDate = docs.get(-1).get('date')

            const expeditions = I.fromJS(expeditionsMap)
                  .map(e => {
                      const id = e.get('slug')
                      return e
                          .set('id', id)
                          .set('name', e.get('name'))
                          .set('startDate', id === expedition.slug ? startDate : Date.now())
                          .set('endDate', id === expedition.slug ? endDate : Date.now())
                          .delete('slug')
                  })

            yield put({
                type: ActionTypes.INITIALIZE_EXPEDITIONS,
                id: expedition.slug,
                data: expeditions
            })

            yield put({
                type: ActionTypes.INITIALIZE_DOCUMENTS,
                data: docs
            })

            console.log(docs)
        }
        else {
            const expeditions = I.fromJS(expeditionsMap)
                  .map(e => {
                      const id = e.get('slug')
                      return e
                          .set('id', id)
                          .set('name', e.get('name'))
                          .set('startDate', Date.now())
                          .set('endDate', Date.now())
                          .delete('slug')
                  })

            yield put({
                type: ActionTypes.SET_ZOOM,
                zoom: 2
            })

            yield put({
                type: ActionTypes.INITIALIZE_EXPEDITIONS,
                id: expedition.slug,
                data: expeditions
            })

            yield put({
                type: ActionTypes.INITIALIZE_DOCUMENTS,
                data: docs
            })
        }
    }
}

export function* rootSaga() {
    yield all([
        refreshSaga(),
        load()
    ])
}
