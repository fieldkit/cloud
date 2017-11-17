import { combineReducers } from 'redux'
import I from 'immutable'

import * as ActionTypes from '../actions/types'

import expeditionReducer from './expeditions'
import navigation from './navigation'
import documents from './documents'

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

function sortedDocuments(docs) {
    return I.fromJS(docs)
        .map(d => {
            return {
                id: d.get('id'),
                type: 'Feature',
                date: d.get('timestamp') * 1000,
                geometry: d.get('location')
            }
        })
        .sort((d1, d2) => {
            return d1.date - d2.date
        })
        .toList()
}

export function visibleExpedition(state = {}, action) {
    let nextState = state

    switch (action.type) {
    case ActionTypes.API_PROJECT_EXPEDITIONS_GET.SUCCESS: {
        if (action.response.length == 0) {
            return I.fromJS({})
        }
        const e = action.response[0]
        return I.fromJS({
            id: e.id,
            name: e.name,
            slug: e.slug,
            description: e.description,
            dates: {
                start: null,
                end: null,
                now: null
            }
        })
    }
    case ActionTypes.API_EXPEDITION_DOCS_GET.SUCCESS: {
        const sorted = sortedDocuments(action.response)
        if (sorted.size == 0) {
            return {}
        }

        const start = sorted.get(0).date
        const end = sorted.get(-1).date
        const now = start
        return state.set('dates', {
            start,
            end,
            now
        })
    }
    default:
        return nextState
    }
}

export function expeditionPanel(state = { open: false }, action) {
    let nextState = state

    switch (action.type) {
    case ActionTypes.OPEN_EXPEDITION_PANEL: {
        return {
            open: true
        }
    }
    case ActionTypes.CLOSE_EXPEDITION_PANEL: {
        return {
            open: false
        }
    }
    default:
        return nextState
    }
}

export default combineReducers({
    project,
    navigation,
    visibleExpedition,
    expeditionPanel,
    expeditions: expeditionReducer,
});

/* Reducer plan:

Timeline (Start, End, Current)

Document Histogram

*/

