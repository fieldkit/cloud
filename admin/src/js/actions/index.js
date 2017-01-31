
import fetch from 'whatwg-fetch'
// import { fetch } from "../vendor_modules/redux-auth"
import * as d3 from 'd3'
import { browserHistory } from 'react-router'
import FKApiClient from '../api/api.js'
import I from 'immutable'
import slug from 'slug'

export const LOGIN_REQUEST = 'LOGIN_REQUEST'
export const LOGIN_SUCCESS = 'LOGIN_SUCCESS'
export const LOGIN_ERROR = 'LOGIN_ERROR'
export const SIGNUP_REQUEST = 'SIGNUP_REQUEST'
export const SIGNUP_SUCCESS = 'SIGNUP_SUCCESS'
export const SIGNUP_ERROR = 'SIGNUP_ERROR'

export const CONNECT = 'CONNECT'
export const DISCONNECT = 'DISCONNECT'

export const UPDATE_EXPEDITION = 'UPDATE_EXPEDITION'
export const EXPEDITION_UPDATED = 'EXPEDITION_UPDATED'


/*

EXPEDITION ACTIONS

*/

export const NEW_PROJECT = 'NEW_PROJECT'
export const NEW_EXPEDITION = 'NEW_EXPEDITION'
export const SET_PROJECT_PROPERTY = 'SET_PROJECT_PROPERTY'
export const SET_EXPEDITION_PROPERTY = 'SET_EXPEDITION_PROPERTY'
export const ADD_DOCUMENT_TYPE = 'ADD_DOCUMENT_TYPE'
export const REMOVE_DOCUMENT_TYPE = 'REMOVE_DOCUMENT_TYPE'
export const SET_EXPEDITION_PRESET = 'SET_EXPEDITION_PRESET'

export function setExpeditionPreset (presetType) {
  return function (dispatch, getState) {
    dispatch ({
      type: SET_EXPEDITION_PRESET,
      presetType
    })
  }
}

export function fetchSuggestedDocumentTypes (input, type, callback) {

  return function (dispatch, getState) {
    window.setTimeout(() => {
      const documentTypes = getState().expeditions
        .get('documentTypes')
        .filter((d) => {
          const nameCheck = d.get('name').toLowerCase().indexOf(input.toLowerCase()) > -1
          const typeCheck = d.get('type') === type
          const membershipCheck = getState().expeditions
            .getIn(['newExpedition', 'documentTypes'])
            .has(d.get('id'))
          return (nameCheck) && !membershipCheck && typeCheck
        })
        .map((m) => {
          return { value: m.get('id'), label: m.get('name')}
        })
        .toArray()

      callback(null, {
        options: documentTypes,
        complete: true
      })
    }, 500)
  }
}

export function addDocumentType (id, collectionType) {
  return function (dispatch, getState) {
    const projectID = getState().expeditions.get('currentProjectID')
    const expedition = getState().expeditions.get('newExpedition')
    const expeditionID = expedition.get('id')
    FKApiClient.addInput(projectID, expeditionID, id)
      .then(res => {
        console.log('server response:', res)
        if (!res) {
          console.log('adding input, error')
        } else {
          console.log('input successfully added')
          // {"id":"LIHQVRPTV7UXRDXMG7F36IRFQ7AIEZBD","expedition_id":"JF55GNWOT4GC3FDPWWX5NT5RWFFJPKRZ","name":"sensor","slug":"sensor"}
          dispatch({
            type: ADD_DOCUMENT_TYPE,
            id,
            collectionType,
            token: res.id
          })
        }
      })
  }
}

export function removeDocumentType (id) {
  return function (dispatch, getState) {
    dispatch({
      type: REMOVE_DOCUMENT_TYPE,
      id
    })
  }
}

export function newProject () {
  return function (dispatch, getState) {
    dispatch(
      {
        type: NEW_PROJECT
      }
    )
  }
}

export function newExpedition () {
  return function (dispatch, getState) {
    dispatch(
      {
        type: NEW_EXPEDITION
      }
    )
  }
}

export function setExpeditionProperty (keyPath, value) {
  return function (dispatch, getState) {
    dispatch({
      type: SET_EXPEDITION_PROPERTY,
      keyPath,
      value
    })
  }
}

export function setProjectProperty (keyPath, value) {
  return function (dispatch, getState) {
    dispatch({
      type: SET_PROJECT_PROPERTY,
      keyPath,
      value
    })
  }
}

export function initNewTeamsSection () {
  return function (dispatch, getState) {
    const actions = []
    if (!getState().expeditions.get('currentTeamID')) {
      dispatch([
        {
          type: ADD_TEAM
        },
        {
          type: START_EDITING_TEAM
        }
      ])
    }
  }
}


/*

TEAM ACTIONS

*/

export const SET_CURRENT_PROJECT = 'SET_CURRENT_PROJECT'
export const SET_CURRENT_EXPEDITION = 'SET_CURRENT_EXPEDITION'
export const SET_CURRENT_TEAM = 'SET_CURRENT_TEAM'
export const SET_CURRENT_MEMBER = 'SET_CURRENT_MEMBER'
export const ADD_TEAM = 'ADD_TEAM'
export const REMOVE_CURRENT_TEAM = 'REMOVE_CURRENT_TEAM'
export const START_EDITING_TEAM = 'START_EDITING_TEAM'
export const STOP_EDITING_TEAM = 'STOP_EDITING_TEAM'
export const SET_TEAM_PROPERTY = 'SET_TEAM_PROPERTY'
export const SET_MEMBER_PROPERTY = 'SET_MEMBER_PROPERTY'
export const CLEAR_CHANGES_TO_TEAM = 'CLEAR_CHANGES_TO_TEAM'
export const SAVE_CHANGES_TO_TEAM = 'SAVE_CHANGES_TO_TEAM'
export const PROMPT_MODAL_CONFIRM_CHANGES = 'PROMPT_MODAL_CONFIRM_CHANGES'
export const CANCEL_ACTION = 'CANCEL_ACTION'
export const CLEAR_MODAL = 'CLEAR_MODAL'

// export const FETCH_SUGGESTED_MEMBERS = 'FETCH_SUGGESTED_MEMBERS'
export const RECEIVE_SUGGESTED_MEMBERS = 'RECEIVE_SUGGESTED_MEMBERS'
export const SUGGESTED_MEMBERS_ERROR = 'SUGGESTED_MEMBERS_ERROR'
export const ADD_MEMBER = 'ADD_MEMBER'
export const REMOVE_MEMBER = 'REMOVE_MEMBER'

export function fetchSuggestedMembers (input, callback) {

  return function (dispatch, getState) {
    window.setTimeout(() => {
      const members = getState().expeditions
        .get('people')
        .filter((m) => {
          const nameCheck = m.get('name').toLowerCase().indexOf(input.toLowerCase()) > -1
          const usernameCheck = m.get('id').toLowerCase().indexOf(input.toLowerCase()) > -1
          const membershipCheck = getState().expeditions
            .getIn(['teams', getState().expeditions.get('currentTeamID'), 'members'])
            .has(m.get('id'))
          return (nameCheck || usernameCheck) && !membershipCheck
        })
        .map((m) => {
          return { value: m.get('id'), label: m.get('name') + ' — ' + m.get('id')}
        })
        .toArray()

      callback(null, {
        options: members,
        complete: true
      })
    }, 500)
  }
}

export function initTeamSection () {
  return function (dispatch, getState) {
    const expeditionID = location.pathname.split('/')[2]
    const teamID = getState().expeditions.getIn(['expeditions', expeditionID,'teams',0])
    dispatch([
      {
        type: SET_CURRENT_EXPEDITION,
        expeditionID
      },
      {
        type: SET_CURRENT_TEAM,
        teamID
      }
    ])
  }
}

export function setCurrentProject (id) {
  return function (dispatch, getState) {
    dispatch({
      type: SET_CURRENT_PROJECT,
      projectID: id
    })
  }
}

export function setCurrentExpedition (id) {
  return function (dispatch, getState) {
    dispatch({
      type: SET_CURRENT_EXPEDITION,
      expeditionID: id
    })
  }
}

export function setCurrentTeam (id) {
  return function (dispatch, getState) {
    const action = {
      type: SET_CURRENT_TEAM,
      teamID: id
    }
    if (!!getState().expeditions.get('editedTeam')) {
      dispatch(promptModalConfirmChanges(null, action))
    } else {
      dispatch(action)
    }
  }
}

export function setCurrentMember (id) {
  return function (dispatch, getState) {
    dispatch({
      type: SET_CURRENT_MEMBER,
      memberID: id
    })
  }
}

export function addTeam () {
  return function (dispatch, getState) {
    const actions = [
      {
        type: ADD_TEAM
      },
      {
        type: START_EDITING_TEAM
      }
    ]
    if (!!getState().expeditions.get('editedTeam')) {
      dispatch(promptModalConfirmChanges(
        null, 
        actions
      ))
    } else {
      dispatch(actions)
    }
  }
}

export function removeCurrentTeam () {
  return function (dispatch, getState) {
    dispatch({
      type: REMOVE_CURRENT_TEAM
    })
  }
}

export function startEditingTeam () {
  return function (dispatch, getState) {
    dispatch({
      type: START_EDITING_TEAM
    })
  }
}

export function stopEditingTeam () {
  return function (dispatch, getState) {
    dispatch({
      type: STOP_EDITING_TEAM
    })
  }
}

export function setTeamProperty (key, value) {
  return function (dispatch, getState) {
    dispatch({
      type: SET_TEAM_PROPERTY,
      key,
      value
    })
  }
}

export function setMemberProperty (memberID, key, value) {
  return function (dispatch, getState) {
    dispatch({
      type: SET_MEMBER_PROPERTY,
      memberID,
      key,
      value
    })
  }
}

export function saveChangesToTeam () {
  return function (dispatch, getState) {
    dispatch({
      type: SAVE_CHANGES_TO_TEAM
    })
  }
}

export function saveChangesAndResume () {
  return function (dispatch, getState) {
    const modal = getState().expeditions.get('modal')
    dispatch([
      {
        type: SAVE_CHANGES_TO_TEAM
      }
    ])
    // console.log('modal', modal.toJS())
    if (!!modal.get('nextAction')) {
      dispatch([
        modal.get('nextAction').toJS(),
        {
          type: CLEAR_MODAL
        }
      ])
    } else if (!!modal.get('nextPath')) {
      // console.log('aga', modal.get('nextPath'))
      dispatch(
        {
          type: CLEAR_MODAL
        }
      )
      browserHistory.push(modal.get('nextPath'))
    }
  }
}

export function clearChangesToTeam () {
  return function (dispatch, getState) {
    dispatch({
      type: CLEAR_CHANGES_TO_TEAM
    })
  }
}

export function promptModalConfirmChanges (nextPath, nextAction) {
  return function (dispatch, getState) {
    dispatch({
      type: PROMPT_MODAL_CONFIRM_CHANGES,
      nextPath,
      nextAction
    })
  }
}

export function cancelAction () {
  return function (dispatch, getState) {
    dispatch({
      type: CLEAR_MODAL
    })
  }
}

export function addMember (id) {
  return function (dispatch, getState) {
    dispatch({
      type: ADD_MEMBER,
      id
    })
  }
}

export function removeMember (id) {
  return function (dispatch, getState) {
    dispatch({
      type: REMOVE_MEMBER,
      id
    })
  }
}



/*

DATA ACTIONS

*/


export const SET_ERROR = 'SET_ERROR'
export const REQUEST_PROJECTS = 'REQUEST_PROJECTS'
export const RECEIVE_PROJECTS = 'RECEIVE_PROJECTS'
export const SAVE_PROJECT = 'SAVE_PROJECT'
export const REQUEST_EXPEDITIONS = 'REQUEST_EXPEDITIONS'
export const RECEIVE_EXPEDITIONS = 'RECEIVE_EXPEDITIONS'
export const SUBMIT_GENERAL_SETTINGS = 'SUBMIT_GENERAL_SETTINGS'
export const RECEIVE_TOKEN = 'RECEIVE_TOKEN'
export const RECEIVE_INPUT = 'RECEIVE_INPUT'
export const SAVE_EXPEDITION = 'SAVE_EXPEDITION'

export function requestProjects (callback) {
  return function (dispatch, getState) {
    FKApiClient.getProjects()
      .then(res => {
        console.log('projects received:', res)
        if (!res) {
          dispatch({
            type: RECEIVE_PROJECTS,
            projects: I.Map()
          })
          if (callback) callback()
        } else {
          const projectMap = {}
          res.forEach(p => {
            projectMap[p.slug] = p
          })
          const projects = I.fromJS(projectMap)
            .map(p => {
              return p
                .set('id', p.get('slug'))
                .merge(I.fromJS({expeditions: []}))
            })
          dispatch(receiveProjects(projects))
        }
      })
  }
}

export function saveProject (id, name) {
  return function (dispatch, getState) {
    FKApiClient.createProjects(name)
      .then(res => {
        console.log('project created', res)
        const projectMap = {};
        [res].forEach(p => {
          projectMap[p.slug] = p
        })
        const projects = I.fromJS(projectMap)
          .map(p => {
            return p.merge(I.fromJS({expeditions: []}))
          })
        dispatch({
          type: SAVE_PROJECT,
          id
        })
        browserHistory.push('/admin/' + id)
      })
  }
}

export function receiveProjects (projects) {
  return function (dispatch, getState) {
    const projectID = projects.toList().getIn([0, 'slug'])
    dispatch([
      {
        type: RECEIVE_PROJECTS,
        projects
      },
      {
        type: SET_CURRENT_PROJECT,
        projectID: projectID
      }
    ])
  }
}

export function requestExpeditions (callback) {
  return function (dispatch, getState) {
    const projectID = getState().expeditions.get('currentProjectID')
    FKApiClient.getExpeditions(projectID)
      .then(res => {
        console.log('expeditions received:', res)
        if (!res) {
          dispatch({
            type: RECEIVE_EXPEDITIONS,
            projectID,
            expeditions: I.Map()
          })
          if (callback) callback()
        } else {
          const expeditionMap = {}
          res.forEach(e => {
            expeditionMap[e.slug] = e
          })
          const expeditions = I.fromJS(expeditionMap)
            .map(e => {
              return e.merge(I.fromJS({
                id: e.get('slug'),
                token: '',
                selectedDocumentType: {},
                documentTypes: {},              
              }))
            })
          dispatch({
            type: RECEIVE_EXPEDITIONS,
            projectID,
            expeditions
          })
        }
      })
  }
}

export function saveGeneralSettings (callback) {
  return function (dispatch, getState) {
    const projectID = getState().expeditions.get('currentProjectID')
    const expedition = getState().expeditions.get('newExpedition')
    const expeditionName = expedition.get('name')
    const expeditionID = expedition.get('id')
    FKApiClient.postGeneralSettings(projectID, expeditionName)
      .then(res => {
        console.log('expeditions successfully saved:', res)
        if (!res) {
          // error
          console.log('error with expedition creation')
        } else {
          FKApiClient.addExpeditionToken(projectID, expeditionID)
            .then(res => {
              console.log('server response:', res)
              if(!res) {
                console.log('no response')
              } else {
                console.log('storing token')
                dispatch([
                  {
                    type: RECEIVE_TOKEN,
                    token: res.ID
                  },
                  {
                    type: SET_ERROR,
                    error: null
                  },
                ])
                if (!!callback) callback()
              }
            })
        }
      })
      .catch(error => {
        dispatch({
          type: SET_ERROR,
          error: I.fromJS(error.message)
        })
      })
  }
}

export function submitInputs () {
  return function (dispatch, getState) {

    const projectID = getState().expeditions.get('currentProjectID')
    const expedition = getState().expeditions.get('newExpedition')
    const expeditionID = expedition.get('id')
    const inputName = expedition.get('documentTypes').toList().get(0).get('id')
    console.log('sending input', projectID, expeditionID, inputName)
    FKApiClient.postInputs(projectID, expeditionID, inputName)
      .then(res => {
        console.log('server response:', res)
        if (!res) {
          // error
          console.log('error with input registration')
        } else {
          console.log('success')
          browserHistory.push('/admin/' + projectID + '/new-expedition/confirmation')
          dispatch({
            type: SET_ERROR,
            error: null
          })
        }
      })
  }
}

export function saveExpedition () {
  return function (dispatch, getState) {
    dispatch ({
      type: SAVE_EXPEDITION
    })
  }
}


export function updateExpedition (expedition) {
  return function (dispatch, getState) {
    window.setTimeout(() => {
      dispatch(expeditionUpdated(expedition))
    }, 1000)
    dispatch({
      type: UPDATE_EXPEDITION,
      expedition
    })
  }
}



export function expeditionUpdated (expedition) {
  return {
    type: EXPEDITION_UPDATED,
    expedition
  }
}


/*

AUTH ACTIONS

*/

export function requestSignIn (username, password) {
  return function (dispatch, getState) {
    dispatch({
      type: LOGIN_REQUEST
    })
    FKApiClient.login(username, password)
      .then(() => {
        FKApiClient.onLogin()
        dispatch({
          type: LOGIN_SUCCESS
        })
        browserHistory.push('/admin')
      })
      .catch(error => {
        console.log('signin error:', error)
        if(error.response) {
          switch(error.response.status){
            case 429:
              dispatch({
                type: SET_ERROR,
                error: I.fromJS('Try again later.')
              })
              break
            case 401:
              dispatch({
                type: SET_ERROR,
                error: I.fromJS('Username or password incorrect.')
              })
              break
          }
        } else {
          dispatch({
            type: SET_ERROR,
            error: I.fromJS('A server error occured')
          })
        }
      })
  }
}

export function requestSignUp (email, username, password, invite) {
  return function (dispatch, getState) {
    dispatch({
      type: SIGNUP_REQUEST
    })
    const params = {
      'email': email,
      'username': username,
      'password': password,
      'invite': invite,
    }
    FKApiClient.register(params)
      .then(() => {
        dispatch({
          type: SIGNUP_SUCCESS
        })
        browserHistory.push('/signin')
      })
      .catch(error => {
        console.log('signup error:', error)
        if (error.response && error.response.status === 400) {
          error.response.json()
            .then(message => {
              dispatch({
                type: SIGNUP_ERROR,
                message
              })
            })
        } else {
          dispatch({
            type: SIGNUP_ERROR,
            message: 'A server error occurred.'
          })
        }
      })
  }
}

export function requestSignOut () {
  return function (dispatch, getState) {
    dispatch({
      type: SIGNOUT_REQUEST
    })
    FKApiClient.logout()
      .then(() => {
        dispatch({
          type: SIGNOUT_SUCCESS
        })
        browserHistory.push('/')
      })
      .catch(error => {
        console.log('signout error:', error)
        if (error.response && error.response.status === 400) {
          error.response.json()
            .then(message => {
              dispatch({
                type: SIGNOUT_ERROR,
                message
              })
            })
        } else {
          dispatch({
            type: SIGNOUT_ERROR,
            message: 'A server error occurred.'
          })
        }
      })
  }
}