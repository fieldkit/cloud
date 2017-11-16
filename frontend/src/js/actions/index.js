
import fetch from 'whatwg-fetch'
import * as d3 from 'd3'
import { browserHistory } from 'react-router'
import FKApiClient from '../api/api.js'
import { getSampleData, updateDeepLinking } from '../utils'

import I from 'immutable'

import * as Types from './types'

export function setViewport(viewport, manual) {
    return (dispatch, getState) => {
        dispatch({
            type: Types.SET_VIEWPORT,
            viewport,
            manual
        })
    }
}

export function updateDate(date, playbackMode, forceUpdate) {
    return (dispatch, getState) => {
        dispatch({
            type: Types.UPDATE_DATE,
            date,
            playbackMode,
            forceUpdate
        })
    }
}

export function selectPlaybackMode(mode) {
    return (dispatch, getState) => {
        dispatch({
            type: Types.SELECT_PLAYBACK_MODE,
            mode
        })
    }
}

export function selectFocusType(focusType) {
    return (dispatch, getState) => {
        dispatch({
            type: Types.SELECT_FOCUS_TYPE,
            focusType
        })
    }
}

export function selectZoom(zoom) {
    return (dispatch, getState) => {
        dispatch({
            type: Types.SELECT_ZOOM,
            zoom
        })
    }
}

export function toggleSensorData() {
    return (dispatch, getState) => {
        dispatch({
            type: Types.TOGGLE_SENSOR_DATA
        })
    }
}

export function setMousePosition(x, y) {
    return (dispatch, getState) => {
        dispatch({
            type: Types.SET_MOUSE_POSITION,
            x,
            y
        })
    }
}

export function openExpeditionPanel() {
    return (dispatch, getState) => {
        dispatch({
            type: Types.OPEN_EXPEDITION_PANEL
        })
    }
}

export function closeExpeditionPanel() {
    return (dispatch, getState) => {
        dispatch({
            type: Types.CLOSE_EXPEDITION_PANEL
        })
    }
}

export function openLightbox(id) {
    return (dispatch, getState) => {
        dispatch({
            type: Types.OPEN_LIGHTBOX,
            id
        })
    }
}

export function closeLightbox() {
    return (dispatch, getState) => {
        dispatch({
            type: Types.CLOSE_LIGHTBOX
        })
    }
}
