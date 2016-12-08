
import * as actions from '../../actions'
import I from 'immutable'

export const initialState = I.fromJS({
  modal: {
    type: null,
    nextAction: null,
    nextPath: null
  },
  currentExpeditionID: null,
  currentTeamID: [],
  currentMemberID: [],
  editedTeam: null,
  expeditions: {
    'okavango_16': {
      id: 'okavango_16',
      name: 'Okavango 2016',
      startDate: new Date('2016-08-17 00:00:00+02:00'),
      teams: [
        'o16-river-team',
        'o16-ground-team'
      ]
    },
    'bike_16': {
      id: 'bike_16',
      name: 'Bike Angola 2016',
      startDate: new Date('2016-06-21 00:00:00+02:00'),
      teams: [
        'b16-ground-team'
      ]
    },
    'cuito_16': {
      id: 'cuito_16',
      name: 'Cuito River 2016',
      startDate: new Date('2016-02-06 00:00:00+02:00'),
      teams: [
        'c16-river-team',
        'c16-ground-team'
      ]
    }
  },
  teams: {
    'o16-river-team': {
      id: 'o16-river-team',
      description: 'Lorem ipsum dolor sit amet, consectetur adipiscing.',
      name: 'river team',
      members: ['steve', 'jer', 'adjany'],
      new: false,
      status: 'ready',
      editing: false
    },
    'o16-ground-team': {
      id: 'o16-ground-team',
      description: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Proin at pellentesque ipsum, sit amet convallis lacus. Donec id dui quis ante congue placerat. Aenean sodales.',
      name: 'ground team',
      members: ['john', 'shah', 'jer'],
      new: false,
      status: 'ready',
      editing: false
    },
    'b16-ground-team': {
      id: 'b16-ground-team',  
      description: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit.',
      name: 'ground team',
      members: ['jer', 'john'],
      new: false,
      status: 'ready',
      editing: false
    },
    'c16-river-team': {
      id: 'c16-river-team',
      description: 'Lorem ipsum dolor sit amet.',
      name: 'river team',
      members: ['steve', 'jer', 'shah'],
      new: false,
      status: 'ready',
      editing: false
    },
    'c16-ground-team': {
      id: 'c16-ground-team',
      description: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Proin at pellentesque ipsum, sit amet convallis lacus.',
      name: 'ground team',
      members: ['john', 'adjany'],
      new: false,
      status: 'ready',
      editing: false
    },
  },
  people: {
    'jer': {
      id: 'jer',
      name: 'Jer Thorp',
      teams: [
        'o16-river-team',
        'c16-river-team',
      ],
      inputs: [
        'ambit',
        'sightings',
        'twitter',
      ]
    },
    'steve': {
      id: 'steve',
      name: 'Steve Boyes',
      teams: [
        'o16-river-team',
        'b16-ground-team',
        'c16-river-team',
      ],
      inputs: [
        'ambit',
        'sightings',
        'twitter',
      ]
    },
    'shah': {
      id: 'shah',
      name: 'Shah Selbe',
      teams: [
        'o16-river-team',
        'c16-ground-team'
      ],
      inputs: [
        'ambit',
        'sightings',
        'twitter',
      ]
    },
    'adjany': {
      id: 'adjany',
      name: 'Adjany Costa',
      teams: [
        'o16-ground-team',
        'b16-ground-team',
        'c16-ground-team'
      ],
      inputs: [
        'ambit',
        'sightings',
        'twitter',
      ]
    },
    'john': {
      id: 'john',
      name: 'John Hilton',
      teams: [
        'o16-river-team',
        'o16-ground-team',
        'b16-ground-team',
        'c16-river-team',
        'c16-ground-team'
      ],
      inputs: [
        'ambit',
        'sightings',
        'twitter',
      ]
    }
  }
})

const expeditionReducer = (state = initialState, action) => {

  // console.log('reducer:', action.type, action)
  switch (action.type) {

    case actions.SET_CURRENT_EXPEDITION: 
      return state.set('currentExpeditionID', action.expeditionID)

    case actions.SET_CURRENT_TEAM:
      return state.set('currentTeamID', action.teamID)

    case actions.SET_CURRENT_MEMBER:
      return state.set('currentMemberID', action.memberID)

    case actions.ADD_TEAM: {
      const expeditionID = state.get('currentExpeditionID')
      const teamID = 'team-' + Date.now()
      return state
        .set('currentTeamID', teamID)
        .setIn(
          ['teams', teamID], 
          I.fromJS({
            id: 'team-' + Date.now(),
            name: 'New Team',
            description: 'Enter a description',
            members: [],
            new: true,
            status: 'new'
          })
        )
        .setIn(
          ['expeditions', expeditionID, 'teams'], 
          state.getIn(['expeditions', expeditionID, 'teams']).push(teamID)
        )
    }

    case actions.REMOVE_CURRENT_TEAM: {
      const expeditionID = state.get('currentExpeditionID')
      const teamID = state.get('currentTeamID')
      return state
        .deleteIn(['teams', teamID])
        .deleteIn(['expeditions', expeditionID, 'teams', 
          state.getIn(['expeditions', expeditionID, 'teams'])
            .findIndex(function(id) {
              return id === teamID
            })
          ]
        )
        .set(
          'currentTeamID', 
          state.getIn(['expeditions', expeditionID, 'teams'])
            .filter(t => {
              return t !== teamID 
            })
            .first()
        )
        .set('editedTeam', null)
    }

    case actions.START_EDITING_TEAM: {
      if (!state.get('editedTeam')) {
        return state
          .set(
            'editedTeam',
            state.getIn(['teams', state.get('currentTeamID')])
          )
          .setIn(
            ['teams', state.get('currentTeamID'), 'status'],
            'editing'
          )
      } else return state
    }

    case actions.STOP_EDITING_TEAM: {
      return state.setIn(
        ['teams', state.get('currentTeamID'), 'status'],
        'ready'
      )
    }

    case actions.SET_TEAM_PROPERTY: {
      return state.setIn(
        ['teams', state.get('currentTeamID'), action.key],
        action.value
      )
    }

    case actions.SAVE_CHANGES_TO_TEAM: {
      return state
        .set('editedTeam', null)
        .setIn(
          ['teams', state.get('currentTeamID'), 'status'],
          'ready'
        ) 
        .setIn(
          ['teams', state.get('currentTeamID'), 'new'],
          false
        )
    }

    case actions.CLEAR_CHANGES_TO_TEAM: {
      return state
        .set('editedTeam', null)
        .setIn(
          ['teams', state.get('currentTeamID')],
          state.get('editedTeam')
        )
    }

    case actions.PROMPT_MODAL_CONFIRM_CHANGES: {
      return state
        .setIn(['modal', 'type'], 'confirm_changes')
        .setIn(['modal', 'nextAction'], I.fromJS(action.nextAction))
        .setIn(['modal', 'nextPath'], action.nextPath)
    }

    case actions.CLEAR_MODAL: {
      return state
        .setIn(['modal', 'type'], null)
        .setIn(['modal', 'nextPath'], null)
        .setIn(['modal', 'nextAction'], null)
    }

    ///////


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
