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

function* asyncApi(callApi) {
    yield put(Object.assign({}, callApi, {
        type: callApi.types.START,
    }));

    try {
        let data;
        if (callApi.method === "GET") {
            data = yield FKApiClient.getJSON(callApi.path);
        }
        else {
            data = yield FKApiClient.postJSON(callApi.path, callApi.body);
        }
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

const AsyncCalls = {};
const PromisedCalls = {};

for (let name of Object.keys(Creators)) {
    const fn = Creators[name];
    PromisedCalls[name] = function() {
        const args = Array.prototype.slice.call(arguments);
        const callApi = fn.apply(null, args);
        if (callApi.method === "GET") {
            return FKApiClient.getJSON(callApi.path).then((data) => {
                return callApi.unwrap != null ? callApi.unwrap(data) : data;
            });
        }
        else {
            return FKApiClient.postJSON(callApi.path, callApi.body).then((data) => {
                return callApi.unwrap != null ? callApi.unwrap(data) : data;
            });
        }
    };
    AsyncCalls[name] = function() {
        const args = Array.prototype.slice.call(arguments);
        return call(asyncApi, fn.apply(null, args));
    };
}

export const FkPromisedApi = PromisedCalls;

export default AsyncCalls;
