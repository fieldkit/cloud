import { combineReducers } from 'redux'
import I from 'immutable'

import * as ActionTypes from '../actions/types'

import expeditionReducer from './expeditions'
import viewportReducer from './expeditions/viewport'

export function project(state = {}, action) {
    let nextState = state;

    switch (action.type) {
    case ActionTypes.API_PROJECT_GET.SUCCESS: {
        return I.fromJS(action.response)
    }
    default:
        return nextState
    }
}

export default combineReducers({
    project,
    viewport: viewportReducer,
    expeditions: expeditionReducer
});
