
import fetch from 'whatwg-fetch'
// import { fetch } from "../vendor_modules/redux-auth"
import * as d3 from 'd3'
import { browserHistory } from 'react-router'
import FKApiClient from '../api/api.js'
import { getSampleData } from '../utils'

import I from 'immutable'

export const REQUEST_EXPEDITION = 'REQUEST_EXPEDITION'
export const INITIALIZE_EXPEDITION = 'INITIALIZE_EXPEDITION'
export const REQUEST_DOCUMENTS = 'REQUEST_DOCUMENTS'
export const INITIALIZE_DOCUMENTS = 'INITIALIZE_DOCUMENTS'
export const SET_VIEWPORT = 'SET_VIEWPORT'
export const UPDATE_DATE = 'UPDATE_DATE'
export const SELECT_PLAYBACK_MODE = 'SELECT_PLAYBACK_MODE'
export const SELECT_FOCUS_TYPE = 'SELECT_FOCUS_TYPE'
export const SELECT_ZOOM = 'SELECT_ZOOM'
export const JUMP_TO = 'JUMP_TO'
export const SET_MOUSE_POSITION = 'SET_MOUSE_POSITION'
export const SET_ZOOM = 'SET_ZOOM'

export function requestExpedition (expeditionID) {
  return function (dispatch, getState) {

    dispatch({
      type: REQUEST_EXPEDITION,
      id: expeditionID
    })

    let projectID = location.hostname.split('.')[0]
    if (projectID === 'localhost') projectID = 'eric'
    console.log('getting expedition')
    FKApiClient.getExpedition(projectID, expeditionID)
      .then(resExpedition => {
      // const resExpedition = {"name":"demoExpedition","slug":"demoexpedition"}
        console.log('expedition received:', resExpedition)
        if (!resExpedition) {
          console.log('error getting expedition')
        } else {
          console.log('expedition properly received')

          // {"name":"ian test","slug":"ian-test"}

          FKApiClient.getDocuments(projectID, expeditionID)
            .then(resDocuments => {
              if (!resDocuments) resDocuments = []
              // const resDocuments = getSampleData()

              console.log('documents received:', resDocuments)
              if (!resDocuments) {
                console.log('error getting documents')
              } else {
                console.log('documents properly received')
                const documentMap = {}
                resDocuments
                  .forEach((d, i) => {
                    d.data.id = d.id
                    d.data.date = d.data.date * 1000
                    if (!d.data.geometry) d.data.geometry = d.data.Geometry
                    documentMap[d.id] = d.data
                  })
                const documents = I.fromJS(documentMap)

                if (documents.size > 0) {
                  const startDate = documents.toList()
                    .sort((d1, d2) => {
                      return d1.get('date') - d2.get('date')
                    })
                    .get(0).get('date')
                  const endDate = documents.toList()
                    .sort((d1, d2) => {
                      return d1.get('date') - d2.get('date')
                    })
                    .get(resDocuments.length - 1).get('date')

                  console.log(new Date(startDate), new Date(endDate))

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
                } else {
                  const startDate = new Date()
                  const endDate = new Date()
                  const expeditionData = I.fromJS({
                    id: expeditionID,
                    name: expeditionID,
                    focusType: 'sensor-reading',
                    startDate,
                    endDate
                  })

                  dispatch([
                    {
                      type: SET_ZOOM,
                      zoom: 2
                    },
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
              }
            })
        }
      }) 
  }
}

export function setViewport(viewport, manual) {
  return function (dispatch, getState) {
    dispatch({
      type: SET_VIEWPORT,
      viewport,
      manual
    })
  }
}

export function updateDate (date, playbackMode, forceUpdate) {
  return function (dispatch, getState) {
    dispatch({
      type: UPDATE_DATE,
      date,
      playbackMode,
      forceUpdate
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

export function selectFocusType (focusType) {
  return function (dispatch, getState) {
    dispatch({
      type: SELECT_FOCUS_TYPE,
      focusType
    })
  }
}

export function selectZoom (zoom) {
  return function (dispatch, getState) {
    dispatch({
      type: SELECT_ZOOM,
      zoom
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
