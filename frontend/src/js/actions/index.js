
import fetch from 'whatwg-fetch'
// import { fetch } from "../vendor_modules/redux-auth"
import * as d3 from 'd3'
import { browserHistory } from 'react-router'
import {FKApiClient} from '../api/api.js'

import I from 'immutable'

export const REQUEST_EXPEDITION = 'REQUEST_EXPEDITION'
export const INITIALIZE_EXPEDITION = 'INITIALIZE_EXPEDITION'
export const REQUEST_DOCUMENTS = 'REQUEST_DOCUMENTS'
export const INITIALIZE_DOCUMENTS = 'INITIALIZE_DOCUMENTS'
export const SET_VIEWPORT = 'SET_VIEWPORT'
export const UPDATE_DATE = 'UPDATE_DATE'
export const SELECT_PLAYBACK_MODE = 'SELECT_PLAYBACK_MODE'
export const JUMP_TO = 'JUMP_TO'

export function requestExpedition (id) {
  return function (dispatch, getState) {
    dispatch({
      type: REQUEST_EXPEDITION,
      id
    })
    setTimeout(() => {
      const res = {
        id: 'okavango',
        name: 'Okavango',
        focusType: 'sensor-reading',
        startDate: 1484328718000,
        endDate: 1484329258000
      }
      dispatch(initializeExpedition(id, res))
      dispatch(requestDocuments(id))
    }, 500)
  }
}

export function initializeExpedition (id, data) {
  return function (dispatch, getState) {
    dispatch({
      type: INITIALIZE_EXPEDITION,
      id,
      data: I.fromJS({
        ...data,
        expeditionFetching: false,
        documentsFetching: true
      })
    })
  }
}

export function requestDocuments (id) {
  return function (dispatch, getState) {
    dispatch({
      type: REQUEST_DOCUMENTS
    })
    setTimeout(() => {
      const res = {
        'reading-0': {
          id: 'reading-0',
          type: 'sensor-reading',
          geometry: {
            type: 'Point',
            coordinates: [125.6, 10.1]
          },
          date: 1484328718000
        },
        'reading-1': {
          id: 'reading-1',
          type: 'sensor-reading',
          geometry: {
            type: 'Point',
            coordinates: [125.68, 10.11]
          },
          date: 1484328818000
        },
        'reading-2': {
          id: 'reading-2',
          type: 'sensor-reading',
          geometry: {
            type: 'Point',
            coordinates: [125.7, 10.31]
          },
          date: 1484328958000
        },
        'reading-3': {
          id: 'reading-3',
          type: 'sensor-reading',
          geometry: {
            type: 'Point',
            coordinates: [125.3, 10.35]
          },
          date: 1484329258000
        }
      }
      dispatch(initializeDocuments(id, I.fromJS(res)))
    }, 500)
  }
}

export function initializeDocuments (id, data) {
  return function (dispatch, getState) {
    dispatch({
      type: INITIALIZE_DOCUMENTS,
      id,
      data
    })
  }
}

export function setViewport(viewport) {
  return function (dispatch, getState) {
    dispatch({
      type: SET_VIEWPORT,
      viewport
    })
  }
}

export function updateDate (date) {
  return function (dispatch, getState) {
    dispatch({
      type: UPDATE_DATE,
      date
    })
  }
}

export function selectPlaybackMode (mode) {
  return function (dispatch, getState) {
    dispatch({
      type: SELECT_PLAYBACK_MODE,
      mode
    })
  }
}

// export const LOGIN_REQUEST = 'LOGIN_REQUEST'

// export function setExpeditionPreset (presetType) {
//   return function (dispatch, getState) {
//     dispatch ({
//       type: SET_EXPEDITION_PRESET,
//       presetType
//     })
//   }
// }

// export function fetchSuggestedDocumentTypes (input, type, callback) {

//   return function (dispatch, getState) {
//     window.setTimeout(() => {
//       const documentTypes = getState().expeditions
//         .get('documentTypes')
//         .filter((d) => {
//           const nameCheck = d.get('name').toLowerCase().indexOf(input.toLowerCase()) > -1
//           const typeCheck = d.get('type') === type
//           const membershipCheck = getState().expeditions
//             .getIn(['expeditions', getState().expeditions.get('currentExpeditionID'), 'documentTypes'])
//             .has(d.get('id'))
//           return (nameCheck) && !membershipCheck && typeCheck
//         })
//         .map((m) => {
//           return { value: m.get('id'), label: m.get('name')}
//         })
//         .toArray()

//       callback(null, {
//         options: documentTypes,
//         complete: true
//       })
//     }, 500)
//   }
// }
