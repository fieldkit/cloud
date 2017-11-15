
import I from 'immutable'
import ViewportMercator from 'viewport-mercator-project'
import { map, constrain } from '../../utils.js'
import * as ActionTypes from '../../actions/types'

export const initialState = I.fromJS({
    project: {
        id: location.hostname.split('.')[0],
        name: location.hostname.split('.')[0]
    },
    lightboxDocumentID: null,
    previousDocumentID: null,
    nextDocumentID: null,
    expeditionPanelOpen: false,
    currentPage: 'map',
    currentExpedition: '',
    playbackMode: 'pause',
    showSensors: true,
    focus: {
        type: 'expedition',
        id: null
    },
    currentDate: new Date(),
    forceDateUpdate: false,
    mousePosition: [-1, -1],
    expeditions: {},
    viewport: {
        latitude: 0,
        longitude: 0,
        zoom: 15,
        width: window.innerWidth,
        height: window.innerHeight,
        startDragLngLat: [0, 0],
        isDragging: false,
        geoBounds: [0, 0, 0, 0]
    },
    currentDocuments: [],
    documents: {}
})

const expeditionReducer = (state = initialState, action) => {
    switch (action.type) {
        case ActionTypes.SET_PROJECT_NAME: {
            return state
                .setIn(['project', 'name'], action.name)
        }
        case ActionTypes.OPEN_LIGHTBOX: {
            const previousDocument = state.get('documents')
                .filter(d => state.get('currentDocuments').includes(d.get('id')))
                .sortBy(d => d.get('date'))
                .filter(d => d.get('date') < state.getIn(['documents', action.id, 'date']))
                .last()

            const nextDocument = state.get('documents')
                .filter(d => state.get('currentDocuments').includes(d.get('id')))
                .sortBy(d => d.get('date'))
                .filter(d => d.get('date') > state.getIn(['documents', action.id, 'date']))
                .first()

            const nextDate = state.getIn(['documents', action.id, 'date'])

            return state
                .set('lightboxDocumentID', action.id)
                .set('previousDocumentID', !!previousDocument ? previousDocument.get('id') : null)
                .set('nextDocumentID', !!nextDocument ? nextDocument.get('id') : null)
                .set('currentDate', nextDate)
                .update('viewport', viewport => updateViewport(state, nextDate, null))
        }
        case ActionTypes.CLOSE_LIGHTBOX: {
            return state
                .set('lightboxDocumentID', null)
        }
        case ActionTypes.SET_CURRENT_PAGE: {
            return state
                .set('currentPage', action.currentPage)
        }
        case ActionTypes.OPEN_EXPEDITION_PANEL: {
            return state
                .set('expeditionPanelOpen', true)
        }
        case ActionTypes.CLOSE_EXPEDITION_PANEL: {
            return state
                .set('expeditionPanelOpen', false)
        }
        case ActionTypes.REQUEST_EXPEDITION: {
            return state
                .set('currentExpedition', action.id)
                .setIn(['expeditions', action.id, 'expeditionFetching'], true)
        }
        case ActionTypes.INITIALIZE_EXPEDITIONS: {
            const currentPage = location.pathname.split('/').filter(p => !!p && p !== state.get('currentExpedition'))[0] || 'map'
            const playbackMode = currentPage === 'map' ? 'forward' : 'pause'
            return state
                .set('expeditions', action.data)
                .set('currentDate', action.data.getIn([action.id, 'startDate']))
                .set('playbackMode', playbackMode)
        }
        case ActionTypes.INITIALIZE_DOCUMENTS: {
            if (!action.data || action.data.size === 0) return state
                    .set('playbackMode', 'pause')
            const position = action.data.toList().get(0).getIn(['geometry', 'coordinates'])
            const currentDocuments = action.data.map(d => d.get('id')).toList()
            const newState = state
                .setIn(['viewport', 'latitude'], position.get(0))
                .setIn(['viewport', 'longitude'], position.get(1))
                .setIn(['expeditions', state.get('currentExpedition'), 'documentsFetching'], false)
                .set('documents', action.data)
                .set('currentDocuments', currentDocuments)
            return newState
        }
        case ActionTypes.SET_VIEWPORT: {
            const {unproject} = ViewportMercator({
                ...action.viewport
            })
            const nw = unproject([0, 0])
            const se = unproject([window.innerWidth, window.innerHeight])
            const geoBounds = [nw[0], nw[1], se[0], se[1]]
            const isDragging = action.viewport.longitude !== state.getIn(['viewport', 'longitude']) ||
                action.viewport.latitude !== state.getIn(['viewport', 'latitude']) ||
                action.viewport.zoom !== state.getIn(['viewport', 'zoom'])
            return state
                .set('playbackMode', action.viewport.isDragging ? 'pause' : state.get('playbackMode'))
                .setIn(['focus', 'type'], isDragging ? 'manual' : state.getIn(['focus', 'type']))
                .setIn(['focus', 'id'], null)
                .set('viewport', state.get('viewport')
                    .merge(action.viewport))
        }
        case ActionTypes.SET_ZOOM: {
            return state
                .setIn(['viewport', 'zoom'], action.zoom)
        }
        case ActionTypes.UPDATE_DATE: {
            const expedition = state.getIn(['expeditions', state.get('currentExpedition')])
            if (!expedition) {
                return state;
            }
            const startDate = expedition.get('startDate')
            const endDate = expedition.get('endDate')
            const nextDate = constrain(action.date, startDate, endDate)
            return state
                .set('currentDate', nextDate)
                .set('forceDateUpdate', !!action.forceUpdate)
                .update('viewport', viewport => updateViewport(state, nextDate, null))
                .update('playbackMode', playbackMode => action.playbackMode || playbackMode)

        }
        case ActionTypes.SELECT_PLAYBACK_MODE: {
            return state
                .set('playbackMode', action.mode)
        }
        case ActionTypes.SELECT_FOCUS_TYPE: {
            const focusType = action.focusType
            return state
                .setIn(['focus', 'type'], focusType)
                .set('playbackMode', 'pause')
                .update('viewport', viewport => updateViewport(state, null, focusType))
        }
        case ActionTypes.SELECT_ZOOM: {
            return state
                .setIn(['viewport', 'zoom'], action.zoom)
                .setIn(['focus', 'type'], 'manual')
                .set('playbackMode', 'pause')
        }
        case ActionTypes.TOGGLE_SENSOR_DATA: {
            let sensors_on = state.get('showSensors')
            return state
                .set('showSensors', !sensors_on)
        }
        case ActionTypes.SET_MOUSE_POSITION: {
            return state
                .setIn(['mousePosition', 0], action.x)
                .setIn(['mousePosition', 1], action.y)
        }
    }
    return state
}

const updateViewport = (state, nextDate, nextFocusType) => {
    const focusType = nextFocusType || state.getIn(['focus', 'type'])
    const date = nextDate || state.get('currentDate')

    switch (focusType) {
        case 'expedition': {
            const expedition = state.getIn(['expeditions', state.get('currentExpedition')])
            const startDate = expedition.get('startDate')
            const endDate = expedition.get('endDate')
            const documents = state.get('documents')
                .filter(d => {
                    return state.get('currentDocuments').includes(d.get('id'))
                })
                .sortBy(d => d.get('date'))

            const previousDocuments = documents
                .filter(d => {
                    return d.get('date') <= date
                })
            const previousDocument = previousDocuments.size > 0 ? previousDocuments.last() :
                date < startDate ? documents.first() :
                    date >= endDate ? documents.last() : null

            const nextDocuments = documents
                .filter(d => {
                    return d.get('date') > date
                })
            const nextDocument = nextDocuments.size > 0 ? nextDocuments.first() :
                date < startDate ? documents.first() :
                    date >= endDate ? documents.last() : null

            if (!previousDocument || !nextDocument) return state.get('viewport')

            const latitude = map(date, previousDocument.get('date'), nextDocument.get('date'), previousDocument.getIn(['geometry', 'coordinates', 0]), nextDocument.getIn(['geometry', 'coordinates', 0]))
            const longitude = map(date, previousDocument.get('date'), nextDocument.get('date'), previousDocument.getIn(['geometry', 'coordinates', 1]), nextDocument.getIn(['geometry', 'coordinates', 1]))

            return state
                .update('viewport', viewport => viewport
                    .set('longitude', longitude)
                    .set('latitude', latitude)
                    .set('zoom', 15)
            )
                .get('viewport')
        }
        case 'documents': {
            const minLng = state.get('documents').minBy(doc => {
                return doc.getIn(['geometry', 'coordinates', 1])
            }).getIn(['geometry', 'coordinates', 1])

            const maxLng = state.get('documents').maxBy(doc => {
                return doc.getIn(['geometry', 'coordinates', 1])
            }).getIn(['geometry', 'coordinates', 1])

            const minLat = state.get('documents').minBy(doc => {
                return doc.getIn(['geometry', 'coordinates', 0])
            }).getIn(['geometry', 'coordinates', 0])

            const maxLat = state.get('documents').maxBy(doc => {
                return doc.getIn(['geometry', 'coordinates', 0])
            }).getIn(['geometry', 'coordinates', 0])

            const latitude = (minLat + maxLat) / 2
            const longitude = (minLng + maxLng) / 2

            const documentWidth = maxLng - minLng
            const documentHeight = maxLat - minLat

            let zoom = state.getIn(['viewport', 'zoom'])
            if ((documentWidth * window.innerWidth) >= (documentHeight * window.innerHeight)) {
                zoom = Math.log2(360 / (documentWidth * 1.1)) + 1
            } else {
                zoom = Math.log2(180 / (documentHeight * 1.1)) + 1
            }

            return state
                .update('viewport', viewport => viewport
                    .set('longitude', longitude)
                    .set('latitude', latitude)
                    .set('zoom', zoom)
            )
                .get('viewport')
        }
        default: {
            return state.get('viewport')
        }
    }
}

export default expeditionReducer
