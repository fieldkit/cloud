import { delay } from 'redux-saga'
import { put, take, takeLatest, takeEvery, select, all, race, call } from 'redux-saga/effects'

export function* refreshSaga() {
    while (true) {
        yield delay(1000);
    }
}

export function* rootSaga() {
    yield all([
        refreshSaga()
    ])
}
