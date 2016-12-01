
import * as actions from '../../actions'
import I from 'immutable'

export const initialState = I.fromJS({
  expeditions: [
    {
      id: 'okavango_16',
      name: 'Okavango 2016',
      updating: false,
      startDate: new Date('2016-08-17 00:00:00+02:00'),
      teams: [
        {
          id: 'river',
          name: 'River Team',
          activity: [5,7,2,4,3,4,5,7,7,9],
          description: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Suspendisse porttitor consequat orci, in sagittis est consequat sodales. Integer euismod arcu non nibh laoreet, imperdiet efficitur quam aliquet. Integer felis felis, euismod ac purus a, aliquet scelerisque tortor.',
          members: [
            {
              id: 'steve',
              name: 'Steve Boyes',
              role: 'Expedition Leader',
              activity: [8,6,7,7,6,7,8,5,3,2],
              inputs: ['Ambit Wristband', 'Sighting App']
            },
            {
              id: 'chris',
              name: 'Chris Boyes',
              role: 'Team Leader',
              activity: [5,7,2,4,3,4,5,7,7,9],
              inputs: ['Sighting App']
            },
            {
              id: 'shah',
              name: 'Shah Selbe',
              role: 'Team Member',
              activity: [1,2,5,7,8,9,8,9,0,0],
              inputs: ['Ambit Wristband', 'Sighting App']
            },
            {
              id: 'jer',
              name: 'Jer Thorp',
              role: 'Team Member',
              activity: [5,7,2,4,3,4,5,7,7,9],
              inputs: ['Ambit Wristband', 'Sighting App', 'Twitter', 'Medium', 'SoundCloud']
            }
          ],
        },
        {
          id: 'ground',
          name: 'Ground Team',
          activity: [5,7,2,4,3,4,5,7,7,9],
          description: 'Integer euismod arcu non nibh laoreet, imperdiet efficitur quam aliquet. Integer felis felis, euismod ac purus a, aliquet scelerisque tortor.',
          members: [
            {
              id: 'john',
              name: 'John Hilton',
              role: 'Team Leader',
              activity: [1,2,5,7,8,9,8,9,0,0],
              inputs: ['Ambit Wristband', 'Sighting App', 'Twitter', 'Medium', 'SoundCloud']
            },
            {
              id: 'adjany',
              name: 'Adjany Costa',
              role: 'Team Member',
              activity: [2,1,4,6,2,4,2,1,4,6],
              inputs: ['Ambit Wristband', 'Sighting App']
            }
          ],
        },
        {
          id: 'river',
          name: 'River Team',
          activity: [5,7,2,4,3,4,5,7,7,9],
          description: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Suspendisse porttitor consequat orci, in sagittis est consequat sodales. Integer euismod arcu non nibh laoreet, imperdiet efficitur quam aliquet. Integer felis felis, euismod ac purus a, aliquet scelerisque tortor.',
          members: [
            {
              id: 'steve',
              name: 'Steve Boyes',
              role: 'Expedition Leader',
              activity: [8,6,7,7,6,7,8,5,3,2],
              inputs: ['Ambit Wristband', 'Sighting App']
            },
            {
              id: 'chris',
              name: 'Chris Boyes',
              role: 'Team Leader',
              activity: [5,7,2,4,3,4,5,7,7,9],
              inputs: ['Sighting App']
            },
            {
              id: 'shah',
              name: 'Shah Selbe',
              role: 'Team Member',
              activity: [1,2,5,7,8,9,8,9,0,0],
              inputs: ['Ambit Wristband', 'Sighting App']
            },
            {
              id: 'jer',
              name: 'Jer Thorp',
              role: 'Team Member',
              activity: [5,7,2,4,3,4,5,7,7,9],
              inputs: ['Ambit Wristband', 'Sighting App', 'Twitter', 'Medium', 'SoundCloud']
            }
          ],
        }
      ]
    },
    {
      id: 'bike_angola_16',
      name: 'Bike Angola 16',
      updating: false,
      startDate: new Date('2016-07-06 00:00:00+02:00'),
      teams: [
        {
          id: 'river',
          name: 'River Team',
          activity: [5,7,2,4,3,4,5,7,7,9],
          description: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Suspendisse porttitor consequat orci, in sagittis est consequat sodales. Integer euismod arcu non nibh laoreet, imperdiet efficitur quam aliquet. Integer felis felis, euismod ac purus a, aliquet scelerisque tortor.',
          members: [
            {
              id: 'steve',
              name: 'Steve Boyes',
              role: 'Expedition Leader',
              activity: [8,6,7,7,6,7,8,5,3,2],
              inputs: ['Ambit Wristband', 'Sighting App']
            },
            {
              id: 'chris',
              name: 'Chris Boyes',
              role: 'Team Leader',
              activity: [5,7,2,4,3,4,5,7,7,9],
              inputs: ['Sighting App']
            },
            {
              id: 'shah',
              name: 'Shah Selbe',
              role: 'Team Member',
              activity: [1,2,5,7,8,9,8,9,0,0],
              inputs: ['Ambit Wristband', 'Sighting App']
            },
            {
              id: 'jer',
              name: 'Jer Thorp',
              role: 'Team Member',
              activity: [5,7,2,4,3,4,5,7,7,9],
              inputs: ['Ambit Wristband', 'Sighting App', 'Twitter', 'Medium', 'SoundCloud']
            }
          ],
        },
        {
          id: 'ground',
          name: 'Ground Team',
          activity: [5,7,2,4,3,4,5,7,7,9],
          description: 'Integer euismod arcu non nibh laoreet, imperdiet efficitur quam aliquet. Integer felis felis, euismod ac purus a, aliquet scelerisque tortor.',
          members: [
            {
              id: 'john',
              name: 'John Hilton',
              role: 'Team Leader',
              activity: [1,2,5,7,8,9,8,9,0,0],
              inputs: ['Ambit Wristband', 'Sighting App', 'Twitter', 'Medium', 'SoundCloud']
            },
            {
              id: 'adjany',
              name: 'Adjany Costa',
              role: 'Team Member',
              activity: [2,1,4,6,2,4,2,1,4,6],
              inputs: ['Ambit Wristband', 'Sighting App']
            }
          ],
        }
      ]  
    },
    {
      id: 'cuando_16',
      name: 'Cuando 16',
      updating: false,
      startDate: new Date('2016-10-01 00:00:00+02:00'),
      teams: [
        {
          id: 'river',
          name: 'River Team',
          activity: [5,7,2,4,3,4,5,7,7,9],
          description: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Suspendisse porttitor consequat orci, in sagittis est consequat sodales. Integer euismod arcu non nibh laoreet, imperdiet efficitur quam aliquet. Integer felis felis, euismod ac purus a, aliquet scelerisque tortor.',
          members: [
            {
              id: 'steve',
              name: 'Steve Boyes',
              role: 'Expedition Leader',
              activity: [8,6,7,7,6,7,8,5,3,2],
              inputs: ['Ambit Wristband', 'Sighting App']
            },
            {
              id: 'chris',
              name: 'Chris Boyes',
              role: 'Team Leader',
              activity: [5,7,2,4,3,4,5,7,7,9],
              inputs: ['Sighting App']
            },
            {
              id: 'shah',
              name: 'Shah Selbe',
              role: 'Team Member',
              activity: [1,2,5,7,8,9,8,9,0,0],
              inputs: ['Ambit Wristband', 'Sighting App']
            },
            {
              id: 'jer',
              name: 'Jer Thorp',
              role: 'Team Member',
              activity: [5,7,2,4,3,4,5,7,7,9],
              inputs: ['Ambit Wristband', 'Sighting App', 'Twitter', 'Medium', 'SoundCloud']
            }
          ],
        },
        {
          id: 'ground',
          name: 'Ground Team',
          activity: [5,7,2,4,3,4,5,7,7,9],
          description: 'Integer euismod arcu non nibh laoreet, imperdiet efficitur quam aliquet. Integer felis felis, euismod ac purus a, aliquet scelerisque tortor.',
          members: [
            {
              id: 'john',
              name: 'John Hilton',
              role: 'Team Leader',
              activity: [1,2,5,7,8,9,8,9,0,0],
              inputs: ['Ambit Wristband', 'Sighting App', 'Twitter', 'Medium', 'SoundCloud']
            },
            {
              id: 'adjany',
              name: 'Adjany Costa',
              role: 'Team Member',
              activity: [2,1,4,6,2,4,2,1,4,6],
              inputs: ['Ambit Wristband', 'Sighting App']
            }
          ],
        }
      ]
    }
  ]
})

const expeditionReducer = (state = initialState, action) => {
  switch (action.type) {
    case actions.UPDATE_EXPEDITION:
      return state
        .setIn([
          'expeditions',
          state.get('expeditions').findIndex(function(e) {
            return e.get('id') === action.expedition.get('id')
          })
        ],
          action.expedition.set('updating', true)
        )
    case actions.EXPEDITION_UPDATED:
      return state
        .setIn([
          'expeditions',
          state.get('expeditions').findIndex(function(e) {
            return e.get('id') === action.expedition.get('id')
          })
        ],
          action.expedition.set('updating', false)
        )
    default:
      return state
  }
}

export default expeditionReducer
