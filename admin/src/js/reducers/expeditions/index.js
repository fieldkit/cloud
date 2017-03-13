
import * as actions from '../../actions'
import I from 'immutable'
import slug from 'slug'

export const initialState = I.fromJS({
  user: null,
  breadcrumbs: [null, null, null],
  errors: null,
  suggestedMembers: null,
  modal: {
    type: null,
    nextAction: null,
    nextPath: null
  },
  currentProject: null,
  currentExpedition: null,
  editedTeam: null,
  projects: null,
  expeditions: null,
  inputs: {
    'conservifymodule': {
      id: 'conservifymodule',
      type: 'sensor',
      name: 'Conservify Module',
      description: 'Lorem ipsum dolor sit amet, consectetur adipiscing.',
      setupType: 'token'
    }
  }
})

const expeditionReducer = (state = initialState, action) => {

  console.log('reducer:', action.type, action)
  try {
    switch (action.type) {
      case actions.RECEIVE_USER: {
        return state.set('user', action.user)
      }

      case actions.SET_BREADCRUMBS: {
        let newState = state
          .setIn(['breadcrumbs', action.level], action.value)
        for (let i = action.level + 1; i < 3; i ++) {
          newState = newState.setIn(['breadcrumbs', i], null)
        }
        return newState
      }

      case actions.SET_ERROR: {
        return state
          .set('errors', action.errors)
      }

      case actions.RECEIVE_PROJECTS: {
        return state
          .set('projects', action.projects)
      }

      case actions.RECEIVE_EXPEDITIONS: {
        return state
          .set('expeditions', action.expeditions)
          .setIn(['currentProject', 'expeditions'], action.expeditions.map(e => e.get('id')).toList())
      }

      case actions.RECEIVE_TOKEN: {
        return state
          .setIn(['currentExpedition', 'token'], action.token)
      }

      case actions.NEW_PROJECT: {
        const projectID = 'project-' + Date.now()
        return state
          .set(
            'currentProject', 
            I.fromJS({
              id: projectID,
              name: 'Project Name',
              expeditions: []
            })
          )
      }

      case actions.SET_PROJECT_PROPERTY: {
        const newState = state.setIn(
          ['currentProject'].concat(action.keyPath),
          action.value
        )
        if (!!action.keyPath && action.keyPath.length === 1 && action.keyPath[0] === 'name') {
          return newState
            .setIn(['currentProject', 'id'], slug(action.value))
        } else {
          return newState 
        }
      }

      case actions.SAVE_PROJECT: {
        return state
          .setIn(
            ['projects', action.id],
            state.get('currentProject')
              .set('id', action.id)
          )
      }

      case actions.NEW_EXPEDITION: {
        const expeditionID = 'expedition-' + Date.now()
        return state
          .set(
            'currentExpedition', 
            I.fromJS({
              id: expeditionID,
              name: 'Expedition Name',
              description: 'Enter a description',
              startDate: new Date(),
              teams: [],
              inputs: [],
              selectedPreset: null,
              token: ''
            })
          )
      }

      case actions.ADD_INPUT: {
        if (!state.getIn(['currentExpedition', 'inputs']).includes(action.id)) {
          return state
            .setIn(
              ['currentExpedition', 'inputs'],
              state.getIn(['currentExpedition', 'inputs']).push(action.id)
            )
            .setIn(
              ['inputs', action.id, 'token'],
              action.token
            )
        } else {
          return state
            .setIn(
              ['inputs', action.id, 'token'],
              action.token
            )
        }
      }

      case actions.REMOVE_INPUT: {
        return state
          .deleteIn(['currentExpedition', 'inputs', action.id])
      }

      case actions.RECEIVE_INPUTS: {
        return state
          .setIn(['currentExpedition', 'inputs'], action.inputs)
      }

      case actions.SET_EXPEDITION_PROPERTY: {
        const newState = state.setIn(
          ['currentExpedition'].concat(action.keyPath),
          action.value
        )
        if (!!action.keyPath && action.keyPath.length === 1 && action.keyPath[0] === 'name') {
          return newState
            .setIn(['currentExpedition', 'id'], slug(action.value))
        } else {
          return newState 
        }
      }

      case actions.SAVE_EXPEDITION: {
        const expedition = state.get('currentExpedition')
        return state.setIn(['expeditions', expedition.get('id')], expedition)
          .set('currentExpedition', null)
      }

      case actions.SET_CURRENT_PROJECT: {
        return state.set('currentProject', state.getIn(['projects', action.projectID]))
      }

      case actions.SET_CURRENT_EXPEDITION: {
        const expedition = state.getIn(['expeditions', action.expeditionID])
        return state
          .set('currentExpedition', expedition)
      }

      // case actions.START_EDITING_TEAM: {
      //   if (!state.get('editedTeam')) {
      //     const currentTeam = state.getIn(['teams', state.get('currentTeamID')])
      //     return state
      //       .set(
      //         'editedTeam',
      //         I.fromJS({
      //           name: currentTeam.get('name'),
      //           description: currentTeam.get('description'),
      //           members: currentTeam.get('members')
      //         })
      //       )
      //       .setIn(
      //         ['teams', state.get('currentTeamID'), 'status'],
      //         'editing'
      //       )
      //   } else return state
      // }

      // case actions.STOP_EDITING_TEAM: {
      //   return state.setIn(
      //     ['teams', state.get('currentTeamID'), 'status'],
      //     'ready'
      //   )
      // }

      // case actions.SET_TEAM_PROPERTY: {
      //   return state.setIn(
      //     ['teams', state.get('currentTeamID'), action.key],
      //     action.value
      //   )
      // }

      // case actions.SET_MEMBER_PROPERTY: {
      //   return state.setIn(
      //     ['teams', state.get('currentTeamID'), 'members', action.memberID, action.key],
      //     action.value
      //   )
      // }

      // case actions.SAVE_CHANGES_TO_TEAM: {
      //   return state
      //     .set('editedTeam', null)
      //     .setIn(
      //       ['teams', state.get('currentTeamID'), 'status'],
      //       'ready'
      //     ) 
      //     .setIn(
      //       ['teams', state.get('currentTeamID'), 'new'],
      //       false
      //     )
      // }

      // case actions.CLEAR_CHANGES_TO_TEAM: {
      //   let newTeam = state.getIn(['teams', state.get('currentTeamID')])
      //   state.getIn(['teams', state.get('currentTeamID')])
      //     .forEach((p, i) => {
      //       if (state.get('editedTeam').has(i)) {
      //         newTeam = newTeam.set(i, state.getIn(['editedTeam', i]))
      //       }
      //     })
      //   return state
      //     .set('editedTeam', null)
      //     .setIn(
      //       ['teams', state.get('currentTeamID')],
      //       newTeam
      //     )
      // }

      // case actions.PROMPT_MODAL_CONFIRM_CHANGES: {
      //   return state
      //     .setIn(['modal', 'type'], 'confirm_changes')
      //     .setIn(['modal', 'nextAction'], I.fromJS(action.nextAction))
      //     .setIn(['modal', 'nextPath'], action.nextPath)
      // }

      // case actions.CLEAR_MODAL: {
      //   return state
      //     .setIn(['modal', 'type'], null)
      //     .setIn(['modal', 'nextPath'], null)
      //     .setIn(['modal', 'nextAction'], null)
      // }

      // case actions.RECEIVE_SUGGESTED_MEMBERS: {
      //   return state
      //     .set('suggestedMembers', action.members)
      // }

      // case actions.CLEAR_SUGGESTED_MEMBERS: {
      //   return state
      //     .set('suggestedMembers', null)
      // }

      // case actions.FETCH_SUGGESTED_MEMBERS: {
      //   return state
      //     .setIn(
      //       ['teams', state.get('currentTeamID'), 'selectedMember'], 
      //       null
      //     )
      // }

      // case actions.ADD_MEMBER: {
      //   let newState = state
      //     .setIn(
      //       ['teams', state.get('currentTeamID'), 'selectedMember'], 
      //       null
      //     )
      //   if (!state.getIn(['teams', state.get('currentTeamID'), 'members']).has(action.id)) {
      //     newState = newState.setIn(
      //       ['teams', state.get('currentTeamID'), 'members', action.id], 
      //       I.fromJS({
      //         id: action.id,
      //         role: 'Team Member'
      //       })
      //     )
      //   }
      //   return newState        
      // }

      // case actions.REMOVE_MEMBER: {
      //   return state
      //     .deleteIn(
      //       ['teams', state.get('currentTeamID'), 'members', action.id]
      //     )
      // }

      default:
        return state
    }
  } catch (err) {
    console.log('error in reducer:', err)
  }
}

export default expeditionReducer
