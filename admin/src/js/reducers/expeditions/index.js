
import * as actions from '../../actions'
import I from 'immutable'
import slug from 'slug'

export const initialState = I.fromJS({
  errors: null,
  suggestedMembers: null,
  modal: {
    type: null,
    nextAction: null,
    nextPath: null
  },
  newProject: null,
  newExpedition: null,
  // currentProjectID: null,
  // currentExpeditionID: null,
  currentTeamID: null,
  currentMemberID: [],
  currentDocumentTypeID: null,
  editedTeam: null,
  projects: null,
  expeditions: null,
  documentTypes: {
    // 'memberGeolocation': {
    //   id: 'memberGeolocation',
    //   type: 'member',
    //   name: 'Member Geolocation',
    //   description: 'Lorem ipsum dolor sit amet, consectetur adipiscing.',
    //   setupType: 'token', 
    //   inputs: ['Ambit wristband']
    // },
    // 'sighting': {
    //   id: 'sighting',
    //   type: 'member',
    //   name: 'Sighting',
    //   description: 'Lorem ipsum dolor sit amet, consectetur adipiscing.',
    //   setupType: 'token', 
    //   inputs: ['Uploader']
    // },
    // 'tweet': {
    //   id: 'tweet',
    //   type: 'social',
    //   name: 'Tweet',
    //   description: 'Lorem ipsum dolor sit amet, consectetur adipiscing.',
    //   setupType: 'token', 
    //   inputs: ['Twitter account']
    // },
    'conservifyModule': {
      id: 'conservifyModule',
      type: 'sensor',
      name: 'Conservify Module',
      description: 'Lorem ipsum dolor sit amet, consectetur adipiscing.',
      setupType: 'token', 
      inputs: ['']
    }
  },
  teams: {
  //   'o16-river-team': {
  //     id: 'o16-river-team',
  //     description: 'Lorem ipsum dolor sit amet, consectetur adipiscing.',
  //     name: 'river team',
  //     members: {
  //       steve: {
  //         id: 'steve',
  //         role: 'Expedition Leader'
  //       },
  //       jer: {
  //         id: 'jer',
  //         role: 'Team Leader'
  //       },
  //       adjany: {
  //         id: 'adjany',
  //         role: 'Team Member'
  //       }
  //     },
  //     new: false,
  //     status: 'ready',
  //     editing: false,
  //     selectedMember: null
  //   },
  //   'o16-ground-team': {
  //     id: 'o16-ground-team',
  //     description: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Proin at pellentesque ipsum, sit amet convallis lacus. Donec id dui quis ante congue placerat. Aenean sodales.',
  //     name: 'ground team',
  //     members: {
  //       steve: {
  //         id: 'steve',
  //         role: 'Expedition Leader'
  //       },
  //       jer: {
  //         id: 'jer',
  //         role: 'Team Leader'
  //       },
  //       adjany: {
  //         id: 'adjany',
  //         role: 'Team Member'
  //       }
  //     },
  //     new: false,
  //     status: 'ready',
  //     editing: false,
  //     selectedMember: null
  //   },
  //   'b16-ground-team': {
  //     id: 'b16-ground-team',  
  //     description: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit.',
  //     name: 'ground team',
  //     members: {
  //       steve: {
  //         id: 'steve',
  //         role: 'Expedition Leader'
  //       },
  //       jer: {
  //         id: 'jer',
  //         role: 'Team Leader'
  //       },
  //       adjany: {
  //         id: 'adjany',
  //         role: 'Team Member'
  //       }
  //     },
  //     new: false,
  //     status: 'ready',
  //     editing: false,
  //     selectedMember: null
  //   },
  //   'c16-river-team': {
  //     id: 'c16-river-team',
  //     description: 'Lorem ipsum dolor sit amet.',
  //     name: 'river team',
  //     members: {
  //       steve: {
  //         id: 'steve',
  //         role: 'Expedition Leader'
  //       },
  //       jer: {
  //         id: 'jer',
  //         role: 'Team Leader'
  //       },
  //       adjany: {
  //         id: 'adjany',
  //         role: 'Team Member'
  //       }
  //     },
  //     new: false,
  //     status: 'ready',
  //     editing: false,
  //     selectedMember: null
  //   },
  //   'c16-ground-team': {
  //     id: 'c16-ground-team',
  //     description: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Proin at pellentesque ipsum, sit amet convallis lacus.',
  //     name: 'ground team',
  //     members: {
  //       steve: {
  //         id: 'steve',
  //         role: 'Expedition Leader'
  //       },
  //       jer: {
  //         id: 'jer',
  //         role: 'Team Leader'
  //       },
  //       adjany: {
  //         id: 'adjany',
  //         role: 'Team Member'
  //       }
  //     },
  //     new: false,
  //     status: 'ready',
  //     editing: false,
  //     selectedMember: null
  //   },
  },
  people: {
    // 'jer': {
    //   id: 'jer',
    //   name: 'Jer Thorp',
    //   role: 'Team Leader',
    //   teams: [
    //     'o16-river-team',
    //     'c16-river-team',
    //   ],
    //   inputs: [
    //     'ambit',
    //     'sightings',
    //     'twitter',
    //   ]
    // },
    // 'steve': {
    //   id: 'steve',
    //   name: 'Steve Boyes',
    //   role: 'Expedition Leader',
    //   teams: [
    //     'o16-river-team',
    //     'b16-ground-team',
    //     'c16-river-team',
    //   ],
    //   inputs: [
    //     'ambit',
    //     'sightings',
    //     'twitter',
    //   ]
    // },
    // 'shah': {
    //   id: 'shah',
    //   name: 'Shah Selbe',
    //   role: 'Team Member',
    //   teams: [
    //     'o16-river-team',
    //     'c16-ground-team'
    //   ],
    //   inputs: [
    //     'ambit',
    //     'sightings',
    //     'twitter',
    //   ]
    // },
    // 'adjany': {
    //   id: 'adjany',
    //   name: 'Adjany Costa',
    //   role: 'Team Member',
    //   teams: [
    //     'o16-ground-team',
    //     'b16-ground-team',
    //     'c16-ground-team'
    //   ],
    //   inputs: [
    //     'ambit',
    //     'sightings',
    //     'twitter',
    //   ]
    // },
    // 'john': {
    //   id: 'john',
    //   name: 'John Hilton',
    //   role: 'Team Leader',
    //   teams: [
    //     'o16-river-team',
    //     'o16-ground-team',
    //     'b16-ground-team',
    //     'c16-river-team',
    //     'c16-ground-team'
    //   ],
    //   inputs: [
    //     'ambit',
    //     'sightings',
    //     'twitter',
    //   ]
    // }
  }
})

const expeditionReducer = (state = initialState, action) => {

  console.log('reducer:', action.type, action)
  switch (action.type) {

    case actions.SET_ERROR: {
      return state
        .set('errors', action.errors)
    }

    case actions.RECEIVE_PROJECTS: {
      return state
        .set('projects', action.projects)
    }

    case actions.RECEIVE_EXPEDITIONS: {
      const expeditions = state.get('expeditions') || I.Map()
      const newExpeditions = expeditions.merge(action.expeditions)
      return state
        .set('expeditions', newExpeditions)
        .setIn(['newProject', 'expeditions'], newExpeditions.map(e => e.get('id')).toList())
    }

    case actions.RECEIVE_TOKEN: {
      return state
        .setIn(['newExpedition', 'token'], action.token)
    }

    case actions.NEW_PROJECT: {
      const projectID = 'project-' + Date.now()
      return state
        .set(
          'newProject', 
          I.fromJS({
            id: projectID,
            name: 'New Project',
            expeditions: []
          })
        )
    }

    case actions.SET_PROJECT_PROPERTY: {
      const newState = state.setIn(
        ['newProject'].concat(action.keyPath),
        action.value
      )
      if (!!action.keyPath && action.keyPath.length === 1 && action.keyPath[0] === 'name') {
        return newState
          .setIn(['newProject', 'id'], slug(action.value))
      } else {
        return newState 
      }
    }

    case actions.SAVE_PROJECT: {
      return state
        .setIn(
          ['projects', action.id],
          state.get('newProject')
            .set('id', action.id)
        )
    }

    case actions.NEW_EXPEDITION: {
      const expeditionID = 'expedition-' + Date.now()
      return state
        .set(
          'newExpedition', 
          I.fromJS({
            id: expeditionID,
            name: 'New Expedition',
            description: 'Enter a description',
            startDate: new Date(),
            teams: [],
            selectedDocumentType: {
              member: null,
              social: null,
              sensor: null
            },
            selectedPreset: null,
            documentTypes: {},
            token: ''
          })
        )
    }

    case actions.ADD_DOCUMENT_TYPE: {
      let newState = state
        .setIn(
          ['newExpedition', 'selectedDocumentType', action.collectionType], 
          null
        )
      if (!state.getIn(['newExpedition', 'documentTypes']).has(action.id)) {
        newState = newState.setIn(
          ['newExpedition', 'documentTypes', action.id], 
          I.fromJS({
            id: action.id
          })
        )
      }
      newState = newState
        .setIn(['documentTypes', action.id, 'token'], action.token)
      return newState
    }

    case actions.REMOVE_DOCUMENT_TYPE: {
      return state
        .deleteIn(['newExpedition', 'documentTypes', action.id])
    }

    case actions.SET_EXPEDITION_PROPERTY: {
      const newState = state.setIn(
        ['newExpedition'].concat(action.keyPath),
        action.value
      )
      if (!!action.keyPath && action.keyPath.length === 1 && action.keyPath[0] === 'name') {
        return newState
          .setIn(['newExpedition', 'id'], slug(action.value))
      } else {
        return newState 
      }
    }

    case actions.SAVE_EXPEDITION: {
      const expedition = state.get('newExpedition')
      return state.setIn(['expeditions', expedition.get('id')], expedition)
        .set('newExpedition', null)
        // .set('currentExpeditionID', expedition.get('id'))
    }

    case actions.SET_CURRENT_PROJECT: {
      return state.set('newProject', state.getIn(['projects', action.projectID]))
    }

    case actions.SET_CURRENT_EXPEDITION: {
      const expedition = state.getIn(['expeditions', action.expeditionID])
      return state
        .set('newExpedition', expedition)
    }

    case actions.SET_CURRENT_TEAM:
      return state.set('currentTeamID', action.teamID)

    case actions.SET_CURRENT_MEMBER:
      return state.set('currentMemberID', action.memberID)

    case actions.START_EDITING_TEAM: {
      if (!state.get('editedTeam')) {
        const currentTeam = state.getIn(['teams', state.get('currentTeamID')])
        return state
          .set(
            'editedTeam',
            I.fromJS({
              name: currentTeam.get('name'),
              description: currentTeam.get('description'),
              members: currentTeam.get('members')
            })
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

    case actions.SET_MEMBER_PROPERTY: {
      return state.setIn(
        ['teams', state.get('currentTeamID'), 'members', action.memberID, action.key],
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
      let newTeam = state.getIn(['teams', state.get('currentTeamID')])
      state.getIn(['teams', state.get('currentTeamID')])
        .forEach((p, i) => {
          if (state.get('editedTeam').has(i)) {
            newTeam = newTeam.set(i, state.getIn(['editedTeam', i]))
          }
        })
      return state
        .set('editedTeam', null)
        .setIn(
          ['teams', state.get('currentTeamID')],
          newTeam
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

    case actions.RECEIVE_SUGGESTED_MEMBERS: {
      return state
        .set('suggestedMembers', action.members)
    }

    case actions.CLEAR_SUGGESTED_MEMBERS: {
      return state
        .set('suggestedMembers', null)
    }

    case actions.FETCH_SUGGESTED_MEMBERS: {
      return state
        .setIn(
          ['teams', state.get('currentTeamID'), 'selectedMember'], 
          null
        )
    }

    case actions.ADD_MEMBER: {
      let newState = state
        .setIn(
          ['teams', state.get('currentTeamID'), 'selectedMember'], 
          null
        )
      if (!state.getIn(['teams', state.get('currentTeamID'), 'members']).has(action.id)) {
        newState = newState.setIn(
          ['teams', state.get('currentTeamID'), 'members', action.id], 
          I.fromJS({
            id: action.id,
            role: 'Team Member'
          })
        )
      }
      return newState        
    }

    case actions.REMOVE_MEMBER: {
      return state
        .deleteIn(
          ['teams', state.get('currentTeamID'), 'members', action.id]
        )
    }

    default:
      return state
  }
}

export default expeditionReducer
