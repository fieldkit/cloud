import { put, call } from 'redux-saga/effects';

import FKApiClient from './api';

import * as Creators from './creators';
import apiMiddleware, { CALL_WEB_API } from './middleware';

export class FkApi {
    constructor(dispatch) {
        this.dispatch = dispatch;
        for (let name of Object.keys(Creators)) {
            const fn = Creators[name];
            this[name] = function() {
                const args = Array.prototype.slice.call(arguments);
                return this.invoke(fn.apply(null, args));
            };
        }

    }

    invoke(callApi) {
        const action = {};
        action[CALL_WEB_API] = callApi;
        return apiMiddleware(null)(this.dispatch.bind(this))(action);
    }

};

function* api(callApi) {
    yield put(Object.assign({}, callApi, {
        type: callApi.types.START,
    }));

    try {
        const data = yield FKApiClient.getJSON(callApi.path);
        const unwrapped = callApi.unwrap != null ? callApi.unwrap(data) : data;
        yield put(Object.assign({}, callApi, {
            type: callApi.types.SUCCESS,
            response: unwrapped,
        }));

        return unwrapped;
    }
    catch (err) {
        yield put(Object.assign({}, callApi, {
            type: callApi.types.FAILURE,
            error: err,
        }));
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
