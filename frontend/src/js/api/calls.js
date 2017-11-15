import { put, call } from 'redux-saga/effects'

import FKApiClient from '../api/api'

import * as Creators from './creators'

function* api(raw) {
    yield put({
        type: raw.types[0],
        message: raw.message
    })

    try {
        const data = yield FKApiClient.getJSON(raw.path)
        const unwrapped = raw.unwrap != null ? raw.unwrap(data) : data
        yield put({
            type: raw.types[1],
            response: unwrapped,
        })

        return unwrapped
    }
    catch (err) {
        yield put({
            type: raw.types[2],
            error: err,
        })
        throw err
    }
}

const Calls = {}

for (let name of Object.keys(Creators)) {
    const fn = Creators[name]
    Calls[name] = function() {
        const args = Array.prototype.slice.call(arguments);
        return call(api, fn.apply(null, args))
    }
}

export default Calls
