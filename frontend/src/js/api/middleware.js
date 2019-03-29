import FKApiClient from './api';

export const CALL_WEB_API = Symbol('Call Web API');

export default store => dispatch => action => {
    const callApi = action[CALL_WEB_API];
    if (typeof callApi === 'undefined') {
        return dispatch(action);
    }

    dispatch(Object.assign({}, callApi, {
        type: callApi.types.START
    }));


    return FKApiClient.getJSON(callApi.path).then(response => {
        dispatch(Object.assign({}, callApi, {
            type: callApi.types.SUCCESS,
            response: response,
        }));

        return response;
    }, err => {
        dispatch(Object.assign({}, callApi, {
            type: callApi.types.FAILURE,
            error: err,
        }));

        return err;
    });
};
