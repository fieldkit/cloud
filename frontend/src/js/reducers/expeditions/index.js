
import * as actions from '../../actions'
import I from 'immutable'
import slug from 'slug'

export const initialState = I.fromJS({
  currentPage: 'map',
  currentExpedition: '',
  expeditions: {
    'okavango': {
      id: 'okavango',
      name: 'Okavango'
    }
  },
  map: {
    initialPosition: [0, 0],
    initialZoom: 7,
    initialSpeed: 100,
    position: [0, 0],
    zoom: 7,
    speed: 100,
    focus: 'sensor-reading'
  },
  documents: {
    'reading-0': {
      id: 'reading-0',
      type: 'sensor-reading',
      geometry: {
        type: 'Point',
        coordinates: [125.6, 10.1]
      },
    }
  }
})

const expeditionReducer = (state = initialState, action) => {

  console.log('reducer:', action.type, action)
  switch (action.type) {
    case actions.INITIALIZE_EXPEDITION: {
      const position = state.get('documents').toList().get(0).getIn(['geometry', 'position'])
      return state
        .set('currentExpedition', action.id)
        .set('initialPosition', position)
        .set('position', position)
    }
  }
  return state
}

export default expeditionReducer
