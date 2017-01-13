
import fetch from 'whatwg-fetch'
// import { fetch } from "../vendor_modules/redux-auth"
import * as d3 from 'd3'
import { browserHistory } from 'react-router'
import {FKApiClient} from '../api/api.js'

export const INITIALIZE_EXPEDITION = 'INITIALIZE_EXPEDITION'

export function initializeExpedition (id) {
  return function (dispatch, getState) {
    dispatch ({
      type: INITIALIZE_EXPEDITION,
      id
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
