
import * as actions from '../../actions'
import I from 'immutable'
import slug from 'slug'

export const initialState = I.fromJS({
  suggestedMembers: null,
  modal: {
    type: null,
    nextAction: null,
    nextPath: null
  },
  currentProjectID: null,
  currentExpeditionID: null,
  currentTeamID: null,
  currentMemberID: [],
  currentDocumentTypeID: null,
  editedTeam: null,
  projects: {},
  expeditions: {
    // 'kayak_adventure': {
    //   id: 'kayak_adventure',
    //   name: 'Kayak Adventure',
    //   startDate: new Date('2016-08-17 00:00:00+02:00'),
    //   teams: [],
    //   selectedDocumentType: {},
    //   selectedPreset: null,
    //   documentTypes: {},
    //   token: 'd0sid0239ud29h2ijbe109eudsoijdo2109u2wdlkn'
    // },
    // 'okavango_16': {
    //   id: 'okavango_16',
    //   name: 'Okavango 2016',
    //   startDate: new Date('2016-08-17 00:00:00+02:00'),
    //   teams: [
    //     'o16-river-team',
    //     'o16-ground-team'
    //   ],
    //   selectedDocumentType: {
    //     member: null,
    //     social: null,
    //     sensor: null
    //   },
    //   selectedPreset: null,
    //   documentTypes: {},
    //   token: 'd0sid0239ud29h2ijbe109eudsoijdo2109u2wdlkn'
    // },
    // 'bike_16': {
    //   id: 'bike_16',
    //   name: 'Bike Angola 2016',
    //   startDate: new Date('2016-06-21 00:00:00+02:00'),
    //   teams: [
    //     'b16-ground-team'
    //   ],
    //   selectedDocumentType: {
    //     member: null,
    //     social: null,
    //     sensor: null
    //   },
    //   selectedPreset: null,
    //   documentTypes: {},
    //   token: 'd0sid0239ud29h2ijbe109eudsoijdo2109u2wdlkn'
    // },
    // 'cuito_16': {
    //   id: 'cuito_16',
    //   name: 'Cuito River 2016',
    //   startDate: new Date('2016-02-06 00:00:00+02:00'),
    //   teams: [
    //     'c16-river-team',
    //     'c16-ground-team'
    //   ],
    //   selectedDocumentType: {
    //     member: null,
    //     social: null,
    //     sensor: null
    //   },
    //   selectedPreset: null,
    //   documentTypes: {},
    //   token: 'd0sid0239ud29h2ijbe109eudsoijdo2109u2wdlkn'
    // }
  },
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

    case actions.RECEIVE_PROJECTS: {
      return state
        .set('projects', action.projects)
    }

    case actions.RECEIVE_EXPEDITIONS: {
      const newExpeditions = state.get('expeditions').merge(action.expeditions)
      return state
        .set('expeditions', newExpeditions)
        .setIn(['projects', state.get('currentProjectID'), 'expeditions'], newExpeditions.map(e => e.get('id')).toList())
    }

    case actions.RECEIVE_TOKEN: {
      return state
        .setIn(['expeditions', action.expeditionID, 'token'], action.token)
    }

    case actions.ADD_PROJECT: {
      const projectID = 'project-' + Date.now()
      return state
        .set('currentProjectID', projectID)
        .setIn(
          ['projects', projectID], 
          I.fromJS({
            id: projectID,
            name: 'New Project',
            expeditions: []
          })
        )
    }

    case actions.SET_PROJECT_PROPERTY: {
      let newState = state.setIn(
        ['projects', state.get('currentProjectID')].concat(action.keyPath),
        action.value
      )
      if (action.keyPath.length === 1 && action.keyPath[0] === 'name' && !!action.value) {
        const lastID = state.get('currentProjectID')
        const id = slug(action.value)
        newState = newState
          .setIn(
            ['projects', id],
            newState.getIn(['projects', state.get('currentProjectID')])
          )
          .setIn(['projects', id, 'id'], id)
          .set('currentProjectID', id)
        if (lastID !== id) {
          newState = newState.deleteIn(['projects', state.get('currentProjectID')])
        }
      }
      return newState
    }

    case actions.ADD_EXPEDITION: {
      const expeditionID = 'expedition-' + Date.now()
      return state
        .set('currentExpeditionID', expeditionID)
        .setIn(
          ['projects', state.get('currentProjectID'), 'expeditions'],
          state.getIn(['projects', state.get('currentProjectID'), 'expeditions']).push(expeditionID)
        )
        .setIn(
          ['expeditions', expeditionID], 
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
            token: 'd0sid0239ud29h2ijbe109eudsoijdo2109u2wdlkn'
          })
        )
    }

    case actions.SET_EXPEDITION_PRESET: {
      switch (action.presetType) {
        case 'rookie': {
          return state
            .setIn(
              ['expeditions', state.get('currentExpeditionID'), 'selectedPreset'],
              action.presetType
            )
            .setIn(
              ['expeditions', state.get('currentExpeditionID'), 'documentTypes'],
              I.fromJS({
                memberGeolocation: {id: 'memberGeolocation'},
                sighting: {id: 'sighting'},
              })
            )
            .setIn(
              ['teams', 'main-team'],
              I.fromJS({
                id: 'main-team',
                description: 'Lorem ipsum dolor sit amet, consectetur adipiscing.',
                name: 'Main Team',
                members: {
                  adjany: {
                    id: 'adjany',
                    role: 'Team Leader'
                  }
                },
                new: false,
                status: 'ready',
                editing: false,
                selectedMember: null
              })
            )
            .setIn(
              ['expeditions', state.get('currentExpeditionID'), 'teams'],
              I.fromJS(['main-team'])
            )
            .set('currentTeamID', 'main-team')
        }

        case 'advanced': {
          return state
            .setIn(
              ['expeditions', state.get('currentExpeditionID'), 'selectedPreset'],
              action.presetType
            )
            .setIn(
              ['expeditions', state.get('currentExpeditionID'), 'documentTypes'],
              I.fromJS({
                memberGeolocation: {id: 'memberGeolocation'},
                tweet: {id: 'tweet'},
                sensorReading: {id: 'sensorReading'},
                sighting: {id: 'sighting'},
              })
            )
            .setIn(
              ['teams', 'main-team'],
              I.fromJS({
                id: 'o16-river-team',
                description: 'Lorem ipsum dolor sit amet, consectetur adipiscing.',
                name: 'Main Team',
                members: {
                  adjany: {
                    id: 'adjany',
                    role: 'Team Leader'
                  }
                },
                new: false,
                status: 'ready',
                editing: false,
                selectedMember: null
              })
            )
        }

        case 'pro': {
          return state
        }
      }
    }

    case actions.ADD_DOCUMENT_TYPE: {
      let newState = state
        .setIn(
          ['expeditions', state.get('currentExpeditionID'), 'selectedDocumentType', action.collectionType], 
          null
        )
      if (!state.getIn(['expeditions', state.get('currentExpeditionID'), 'documentTypes']).has(action.id)) {
        newState = newState.setIn(
          ['expeditions', state.get('currentExpeditionID'), 'documentTypes', action.id], 
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
        .deleteIn(['expeditions', state.get('currentExpeditionID'), 'documentTypes', action.id])
    }

    case actions.SET_EXPEDITION_PROPERTY: {
      let newState = state.setIn(
        ['expeditions', state.get('currentExpeditionID')].concat(action.keyPath),
        action.value
      )
      if (action.keyPath.length === 1 && action.keyPath[0] === 'name' && !!action.value) {
        const lastID = state.get('currentExpeditionID')
        const id = slug(action.value)
        newState = newState
          .setIn(
            ['expeditions', id],
            newState.getIn(['expeditions', state.get('currentExpeditionID')])
          )
          .setIn(['expeditions', id, 'id'], id)
          .set('currentExpeditionID', id)
          .setIn(
            ['projects', state.get('currentProjectID'), 'expeditions'],
            state.getIn(['projects', state.get('currentProjectID'), 'expeditions'])
              .pop()
              .push(id)
          )
        if (lastID !== id) {
          newState = newState.deleteIn(['expeditions', state.get('currentExpeditionID')])
        }
      }
      return newState
    }

    case actions.SET_CURRENT_PROJECT: 
      return state.set('currentProjectID', action.projectID)

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
            members: {},
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
