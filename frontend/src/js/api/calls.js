import { put, call } from 'redux-saga/effects';

import FKApiClient from './api';

import * as Creators from './creators';

function* api(raw) {
    yield put({
        type: raw.types.START,
        message: raw.message
    });

    try {
        const data = yield FKApiClient.getJSON(raw.path);
        const unwrapped = raw.unwrap != null ? raw.unwrap(data) : data;
        yield put({
            type: raw.types.SUCCESS,
            response: unwrapped,
        });

        return unwrapped;
    }
    catch (err) {
        yield put({
            type: raw.types.FAILURE,
            error: err,
        });
        throw err;
    }
}

const Calls = {};

for (let name of Object.keys(Creators)) {
    const fn = Creators[name];
    Calls[name] = function() {
        const args = Array.prototype.slice.call(arguments);
        return call(api, fn.apply(null, args));
    };
}

export default Calls;
