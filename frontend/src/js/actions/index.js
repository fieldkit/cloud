
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
export const SET_MOUSE_POSITION = 'SET_MOUSE_POSITION'

// export function requestExpedition (id) {
//   return function (dispatch, getState) {
//     dispatch({
//       type: REQUEST_EXPEDITION,
//       id
//     })
//     setTimeout(() => {
//       const res = {
//         id: 'okavango',
//         name: 'Okavango',
//         focusType: 'sensor-reading',
//         startDate: 1484328718000,
//         endDate: 1484329258000
//       }
//       dispatch(initializeExpedition(id, res))
//       dispatch(requestDocuments(id))
//     }, 500)
//   }
// }


export function requestExpedition (expeditionID) {
  return function (dispatch, getState) {

    dispatch({
      type: REQUEST_EXPEDITION,
      id: expeditionID
    })

    let projectID = location.hostname.split('.')[0]
    if (projectID === 'localhost') projectID = 'demo'
    console.log('getting expedition')
    // FKApiClient.get().getExpedition(projectID, expeditionID)
      // .then(resExpedition => {
      const resExpedition = {"name":"demoExpedition","slug":"demoexpedition"}
        console.log('expedition received:', resExpedition)
        if (!resExpedition) {
          console.log('error getting expedition')
        } else {
          console.log('expedition properly received')

          // {"name":"ian test","slug":"ian-test"}


          // FKApiClient.get().getDocuments(projectID, expeditionID)
            // .then(resDocuments => {
                const resDocuments = [{"id":"GIUJ2WATF2HRCWBEKI34FUAVFZ4S7OHE","message_id":"PGMOCCWTIHR527UKISQ4EV4JYG55DJDU","request_id":"DZKYM26OJL2ENEQ5DG6VLDHMIXMTXQN2","input_id":"SCX6732QLCM7VVKUR3EBFPCF2MDJ5JPV","data":{"Geometry":{"coordinates":[38.90714645385742,-77.04476165771484],"type":"Point"},"date":1.484660606e+09,"type":"Feature"}},{"id":"J65LGHQM22QMVAZMYA2S4IEVZAQSZR4W","message_id":"VLGTLUGTCF7QVT6TOOEN3AFJHATXQVAH","request_id":"KDYJEKGJS3ZH5NDX5IFCR7LVHQNLIFNP","input_id":"SCX6732QLCM7VVKUR3EBFPCF2MDJ5JPV","data":{"Geometry":{"coordinates":[38.99214553833008,-76.80714416503906],"type":"Point"},"date":1.484611955e+09,"type":"Feature"}},{"id":"PI5KOJMYQPMA6VPUEPYGNUB2K3MULR45","message_id":"XGYLMLC5HHQE7LR3RG43LV5M4ER75C25","request_id":"53OTLO72IIXCTMN55PKUH3ZNCMTOD5KM","input_id":"SCX6732QLCM7VVKUR3EBFPCF2MDJ5JPV","data":{"Geometry":{"coordinates":[38.90659713745117,-77.04178619384766],"type":"Point"},"date":1.484660868e+09,"type":"Feature"}},{"id":"UDTRVNZGAYTWEVJLJWV6NYG4SRFVGT4Z","message_id":"NLAJGOVGZQMCADGQEVD5JTXV3HAPBHR2","request_id":"SVYSPQGB37YPFNHI4U2COZEENV4O4ZEG","input_id":"SCX6732QLCM7VVKUR3EBFPCF2MDJ5JPV","data":{"Geometry":{"coordinates":[38.904422760009766,-77.04060363769531],"type":"Point"},"date":1.484661089e+09,"type":"Feature"}},{"id":"W6CTJG6AGPNFLJATQRSISBWUIDE7GEQR","message_id":"QUL6V5QSLHJ2Y5CGKYUSC3NM3APHQCBD","request_id":"3KGKE44IVXRQOJKZERTSU5EHO44GTT4U","input_id":"SCX6732QLCM7VVKUR3EBFPCF2MDJ5JPV","data":{"date":1.484665769e+09,"geometry":{"coordinates":[38.903324127197266,-77.03483581542969],"type":"Point"},"type":"Feature"}},{"id":"4HNYNXJVQYGCB6OEQVGH3AK6TKUQGZPS","message_id":"2TB3BCSU5Y6K2J3MTZ6NPZFFWTGLJOUK","request_id":"RXYNJ6IUP5FEAGXH5EAHUZMSOTPLPYI6","input_id":"SCX6732QLCM7VVKUR3EBFPCF2MDJ5JPV","data":{"Geometry":{"coordinates":[38.90373611450195,-77.03819274902344],"type":"Point"},"date":1.484661273e+09,"type":"Feature"}}]
              console.log('documents received:', resDocuments)
              if (!resDocuments) {
                console.log('error getting documents')
              } else {
                console.log('documents properly received')
                // const documents = resDocuments.map(d => {
                //   d.data = d.data * 1000
                //   return d
                //     .set('date', d.data.date * 1000)
                // })

                const documentMap = {}
                resDocuments.forEach(d => {
                  d.data.id = d.id
                  d.data.date = d.data.date * 1000
                  d.data.geometry = d.data.Geometry
                  documentMap[d.id] = d.data
                })
                const documents = I.fromJS(documentMap)
                console.log('WOPPP', documents.toJS())

                const startDate = documents.toList().get(0).get('date')
                const endDate = documents.toList().get(resDocuments.length - 1).get('date')

                const expeditionData = I.fromJS({
                  id: expeditionID,
                  name: expeditionID,
                  focusType: 'sensor-reading',
                  startDate,
                  endDate
                })

                dispatch([
                  {
                    type: INITIALIZE_EXPEDITION,
                    id: expeditionID,
                    data: expeditionData
                  },
                  {
                    type: INITIALIZE_DOCUMENTS,
                    data: documents
                  }
                ])
              }
            // })

          // const expeditionMap = {}
          // res.forEach(e => {
          //   expeditionMap[e.slug] = e
          // })
          // const expeditions = I.fromJS(expeditionMap)
          //   .map(e => {
          //     return e.merge(I.fromJS({
          //       id: e.get('slug'),
          //       token: '',
          //       selectedDocumentType: {},
          //       documentTypes: {},              
          //     }))
          //   })
          // dispatch(receiveExpeditions(projectID, expeditions, false))
        }
      // }) 
  }
}


export function requestExpeditions () {
  return function (dispatch, getState) {
    const projectID = location.hostname.split('.')[0]
    const expeditionID = getState().expeditions.get('currentExpeditionID')
    FKApiClient.get().getExpedition(projectID, expeditionID)
      .then(res => {
        console.log('expeditions received:', res)
        if (!res) {
          console.log('error getting expedition')
        } else {
          console.log('expedition properly received')

          // {"name":"ian test","slug":"ian-test"}

          

          // const expeditionMap = {}
          // res.forEach(e => {
          //   expeditionMap[e.slug] = e
          // })
          // const expeditions = I.fromJS(expeditionMap)
          //   .map(e => {
          //     return e.merge(I.fromJS({
          //       id: e.get('slug'),
          //       token: '',
          //       selectedDocumentType: {},
          //       documentTypes: {},              
          //     }))
          //   })
          // dispatch(receiveExpeditions(projectID, expeditions, false))
        }
      })
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

export function setMousePosition (x, y) {
  return function (dispatch, getState) {
    dispatch({
      type: SET_MOUSE_POSITION,
      x,
      y
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
