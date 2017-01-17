
import * as actions from '../../actions'
import I from 'immutable'
import slug from 'slug'
import ViewportMercator from 'viewport-mercator-project'
import { map, constrain } from '../../utils.js'

export const initialState = I.fromJS({
  currentPage: 'map',
  currentExpedition: '',
  playbackMode: 'pause',
  currentDate: new Date(),
  mousePosition: [-1, -1],
  expeditions: {},
  viewport: {
    latitude: -18.5699229,
    longitude: 22.115456,
    zoom: 10,
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
      return state
        .setIn(['expeditions', action.id], action.data)
        .set('currentDate', action.data.get('startDate'))
        .set('playbackMode', 'forward')
    }

    case actions.INITIALIZE_DOCUMENTS: {
      console.log('INIT', action.data.toJS())
      const position = action.data.toList().get(0).getIn(['geometry', 'coordinates'])
      const currentDocuments = action.data.map(d => d.get('id'))
      return state
        .setIn(['viewport', 'longitude'], position.get(0))
        .setIn(['viewport', 'latitude'], position.get(1))
        .setIn(['expeditions', state.get('currentExpedition'), 'documentsFetching'], false)
        .set('currentDocuments', currentDocuments)
    }

    case actions.SET_VIEWPORT: {
      const { unproject } = ViewportMercator({ ...action.viewport })
      const nw = unproject([0, 0])
      const se = unproject([window.innerWidth, window.innerHeight])
      const geoBounds = [nw[0], nw[1], se[0], se[1]]
      return state
        .setIn(['viewport', 'geoBounds'], geoBounds)
        .set('viewport', state.get('viewport')
          .merge(action.viewport))
    }

    case actions.UPDATE_DATE: {
      const expedition = state.getIn(['expeditions', state.get('currentExpedition')])
      const startDate = expedition.get('startDate')
      const endDate = expedition.get('endDate')
      const nextDate = constrain(action.date, startDate, endDate)

      const documents = state.get('documents')
        .filter(d => {
          return state.get('currentDocuments').includes(d.get('id'))
        })
        .sort((d1, d2) => {
          return d2.get('date') < d1.get('date')
        })

      const previousDocuments = documents
        .filter(d => {
          return d.get('date') <= nextDate
        })
      const previousDocument = previousDocuments.size > 0 ? previousDocuments.last() : 
        nextDate < startDate ? documents.first() :
        nextDate >= endDate ? documents.last() : null

      const nextDocuments = documents
        .filter(d => {
          return d.get('date') > nextDate
        })
      const nextDocument = nextDocuments.size > 0 ? nextDocuments.first() : 
        nextDate < startDate ? documents.first() :
        nextDate >= endDate ? documents.last() : null

      const longitude = map(nextDate, previousDocument.get('date'), nextDocument.get('date'), previousDocument.getIn(['geometry', 'coordinates', 0]), nextDocument.getIn(['geometry', 'coordinates', 0]))
      const latitude = map(nextDate, previousDocument.get('date'), nextDocument.get('date'), previousDocument.getIn(['geometry', 'coordinates', 1]), nextDocument.getIn(['geometry', 'coordinates', 1]))

      return state
        .set('currentDate', nextDate)
        .setIn(['viewport', 'longitude'], longitude)
        .setIn(['viewport', 'latitude'], latitude)
    }

    case actions.SELECT_PLAYBACK_MODE: {
      return state
        .set('playbackMode', action.mode)
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
