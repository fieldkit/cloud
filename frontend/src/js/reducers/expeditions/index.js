
import * as actions from '../../actions'
import I from 'immutable'
import slug from 'slug'
import ViewportMercator from 'viewport-mercator-project'
import { map, constrain } from '../../utils.js'

export const initialState = I.fromJS({
  currentPage: 'map',
  currentExpedition: '',
  playbackMode: 'pause',
  focus: {
    type: 'expedition',
    id: null
  },
  currentDate: new Date(),
  mousePosition: [-1, -1],
  expeditions: {},
  viewport: {
    latitude: 0,
    longitude: 0,
    zoom: 15,
    width: window.innerWidth,
    height: window.innerHeight,
    startDragLngLat: null,
    isDragging: false,
    geoBounds: [0, 0, 0, 0]
  },
  currentDocuments: [],
  documents: {
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
})

const updateViewport = (state, nextDate, nextFocusType) => {
  const focusType = nextFocusType || state.getIn(['focus', 'type'])
  const date = nextDate || state.get('currentDate')

  switch (focusType) {
    case 'expedition' : {
      const expedition = state.getIn(['expeditions', state.get('currentExpedition')])
      const startDate = expedition.get('startDate')
      const endDate = expedition.get('endDate')

      const documents = state.get('documents')
        .filter(d => {
          return state.get('currentDocuments').includes(d.get('id'))
        })
        .sort((d1, d2) => {
          return d1.get('date') - d2.get('date')
        })

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

      const longitude = map(date, previousDocument.get('date'), nextDocument.get('date'), previousDocument.getIn(['geometry', 'coordinates', 1]), nextDocument.getIn(['geometry', 'coordinates', 1]))
      const latitude = map(date, previousDocument.get('date'), nextDocument.get('date'), previousDocument.getIn(['geometry', 'coordinates', 0]), nextDocument.getIn(['geometry', 'coordinates', 0]))

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

const expeditionReducer = (state = initialState, action) => {

  if (action.type !== actions.UPDATE_DATE && action.type !== actions.SET_MOUSE_POSITION ) {
    console.log('reducer:', action.type, action)
  }
  switch (action.type) {

    case actions.REQUEST_EXPEDITION: {
      return state
        .set('currentExpedition', action.id)
        .setIn(['expeditions', action.id, 'expeditionFetching'], true)
    }

    case actions.INITIALIZE_EXPEDITION: {
      const currentPage = location.pathname.split('/').filter(p => !!p && p !== state.get('currentExpedition'))[0] || 'map'
      const playbackMode = currentPage === 'map' ? 'forward' : 'pause'
      return state
        .setIn(['expeditions', action.id], action.data)
        .set('currentDate', action.data.get('startDate'))
        .set('playbackMode', playbackMode)
    }

    case actions.INITIALIZE_DOCUMENTS: {
      const position = action.data.toList().get(0).getIn(['geometry', 'coordinates'])
      const currentDocuments = action.data.map(d => d.get('id')).toList()
      const newState = state
        .setIn(['viewport', 'longitude'], position.get(0))
        .setIn(['viewport', 'latitude'], position.get(1))
        .setIn(['expeditions', state.get('currentExpedition'), 'documentsFetching'], false)
        .set('documents', action.data)
        .set('currentDocuments', currentDocuments)
      return newState
    }

    case actions.SET_VIEWPORT: {
      const { unproject } = ViewportMercator({ ...action.viewport })
      const nw = unproject([0, 0])
      const se = unproject([window.innerWidth, window.innerHeight])
      const geoBounds = [nw[0], nw[1], se[0], se[1]]

      return state
        .set('viewport', I.fromJS(action.viewport))
        .setIn(['viewport', 'geoBounds'], geoBounds)
        .set('playbackMode', action.viewport.isDragging ? 'pause' : state.get('playbackMode'))
        .setIn(['focus', 'type'], action.viewport.isDragging ? 'manual' : state.getIn(['focus', 'type']))
        .setIn(['focus', 'id'], null)
        .set('viewport', state.get('viewport')
          .merge(action.viewport))
    }

    case actions.SET_ZOOM: {
      return state
        .setIn(['viewport', 'zoom'], action.zoom)
    }

    case actions.UPDATE_DATE: {
      const expedition = state.getIn(['expeditions', state.get('currentExpedition')])
      const startDate = expedition.get('startDate')
      const endDate = expedition.get('endDate')
      const nextDate = constrain(action.date, startDate, endDate)
      return state
        .set('currentDate', nextDate)
        .update('viewport', viewport => updateViewport(state, nextDate, null))
        .update('playbackMode', playbackMode => action.playbackMode || playbackMode)
    }

    case actions.SELECT_PLAYBACK_MODE: {
      return state
        .set('playbackMode', action.mode)
    }

    case actions.SELECT_FOCUS_TYPE: {
      const focusType = action.focusType
      return state
        .setIn(['focus', 'type'], focusType)
        .set('playbackMode', 'pause')
        .update('viewport', viewport => updateViewport(state, null, focusType))
    }

    case actions.SELECT_ZOOM: {
      return state
        .setIn(['viewport', 'zoom'], action.zoom)
        .setIn(['focus', 'type'], 'manual')
        .set('playbackMode', 'pause')
    }

    case actions.SET_MOUSE_POSITION: {
      return state
        .setIn(['mousePosition', 0], action.x)
        .setIn(['mousePosition', 1], action.y)
    }
    
  }
  return state
}

export default expeditionReducer
