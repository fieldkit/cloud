
import * as actions from '../../actions'
import I from 'immutable'
import slug from 'slug'

export const initialState = I.fromJS({
  currentPage: 'map',
  currentExpedition: '',
  expeditions: {
    'okavango': {
      id: 'okavango',
      name: 'Okavango',
      startDate: new Date(0),
      endDate: new Date()
    }
  },
  viewport: {
    latitude: -18.5699229,
    longitude: 22.115456,
    zoom: 4,
    width: window.innerWidth,
    height: window.innerHeight,
    startDragLngLat: null,
    isDragging: false
  },
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
      id: 'reading-0',
      type: 'sensor-reading',
      geometry: {
        type: 'Point',
        coordinates: [125.68, 10.11]
      },
      date: 1484328818000
    },
    'reading-2': {
      id: 'reading-0',
      type: 'sensor-reading',
      geometry: {
        type: 'Point',
        coordinates: [125.7, 10.31]
      },
      date: 1484328958000
    },
    'reading-3': {
      id: 'reading-0',
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

  console.log('reducer:', action.type, action)
  switch (action.type) {
    case actions.INITIALIZE_EXPEDITION: {
      const position = state.get('documents').toList().get(0).getIn(['geometry', 'coordinates'])

      const startDate = state.get('documents').toList().minBy(d => d.date)
      const endDate = state.get('documents').toList().maxBy(d => d.date)

      console.log(new Date(startDate), new Date(endDate))

      return state
        .set('currentExpedition', action.id)
        .setIn(['viewport', 'longitude'], position.get(0))
        .setIn(['viewport', 'latitude'], position.get(1))
    }

    case actions.SET_VIEWPORT: {
      return state
        .set('viewport', state.get('viewport').merge(action.viewport))
    }
  }
  return state
}

export default expeditionReducer
