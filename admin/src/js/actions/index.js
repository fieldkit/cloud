
import fetch from 'whatwg-fetch'
// import { fetch } from "../vendor_modules/redux-auth"
import * as d3 from 'd3'
import { browserHistory } from 'react-router'
import {FKApiClient} from '../api/api.js'
import I from 'immutable'

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

export const ADD_EXPEDITION = 'ADD_EXPEDITION'
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
            .getIn(['expeditions', getState().expeditions.get('currentExpeditionID'), 'documentTypes'])
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
    dispatch({
      type: ADD_DOCUMENT_TYPE,
      id,
      collectionType
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

export function initNewExpeditionSection () {
  return function (dispatch, getState) {
    dispatch(
      {
        type: ADD_EXPEDITION
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


export const REQUEST_PROJECTS = 'REQUEST_PROJECTS'
export const RECEIVE_PROJECTS = 'RECEIVE_PROJECTS'
export const REQUEST_EXPEDITIONS = 'REQUEST_EXPEDITIONS'
export const RECEIVE_EXPEDITIONS = 'RECEIVE_EXPEDITIONS'
export const SUBMIT_GENERAL_SETTINGS = 'SUBMIT_GENERAL_SETTINGS'
export const RECEIVE_TOKEN = 'RECEIVE_TOKEN'

export function requestProjects () {
  return function (dispatch, getState) {
    FKApiClient.get().getProjects()
      .then(res => {
        console.log('projects received:', res)
        if (!res) {
          dispatch(createProject ('new project'))
        } else {
          const projectMap = {}
          res.forEach(p => {
            projectMap[p.slug] = p
          })
          const projects = I.fromJS(projectMap)
            .map(p => {
              return p.merge(I.fromJS({expeditions: []}))
            })
          dispatch(receiveProjects(projects))

        }
      })
  }
}

export function createProject (name) {
  return function (dispatch, getState) {
    FKApiClient.get().createProjects(name)
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
        dispatch(receiveProjects(projects))
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
    browserHistory.push('/admin/' + projects.toList().getIn([0, 'slug']))
  }
}

export function requestExpeditions () {
  return function (dispatch, getState) {
    const projectID = getState().expeditions.get('currentProjectID')
    FKApiClient.get().getExpeditions(projectID)
      .then(res => {
        console.log('expeditions received:', res)
        if (!res) {
          browserHistory.push('/admin/' + projectID + '/new-expedition')
        } else {
          const expeditionMap = {}
          res.forEach(e => {
            expeditionMap[e.slug] = e
          })
          const expeditions = I.fromJS(expeditionMap)
            .map(e => {
              return e.merge(I.fromJS({
                id: e.get('slug'),
                token: 'd0sid0239ud29h2ijbe109eudsoijdo2109u2wdlkn',
                selectedDocumentType: {},
                documentTypes: {},              
              }))
            })
          dispatch(receiveExpeditions(projectID, expeditions, false))
        }
      })
  }
}

export function submitGeneralSettings () {
  return function (dispatch, getState) {
    const projectID = getState().expeditions.get('currentProjectID')
    const expeditionID = getState().expeditions.get('currentExpeditionID')
    const expeditionName = getState().expeditions.getIn(['expeditions', expeditionID, 'name'])
    FKApiClient.get().postGeneralSettings(projectID, expeditionName)
      .then(res => {
        console.log('expeditions successfully saved:', res)
        if (!res) {
          // error
          console.log('error with expedition creation')
        } else {
          const expeditions = I.fromJS({
            [res.slug]: {
              id: res.slug,
              name: res.name,
              token: '',
              selectedDocumentType: {},
              documentTypes: {},
            }
          })
          dispatch(receiveExpeditions(projectID, expeditions, false))
          FKApiClient.get().addExpeditionToken(projectID, expeditionID)
            .then(res => {
              console.log('server response:', res)
              if(!res) {
                console.log('no response')
              } else {
                console.log('storing token')
                // {"ID":"B6XFKIV772JO2SZOHNK57AUQ6LL3PBV4","ExpeditionID":"LRRUHSZIYTM4XZVMYKVGSIQ5UVMYLNPF"}
                dispatch({
                  type: RECEIVE_TOKEN,
                  expeditionID,
                  token: res.ExpeditionID
                })
                browserHistory.push('/admin/' + projectID + '/new-expedition/inputs')
              }
            })
        }
      })
  }
}

export function submitInputs () {
  return function (dispatch, getState) {

    const projectID = getState().expeditions.get('currentProjectID')
    const expeditionID = getState().expeditions.get('currentExpeditionID')
    const inputName = getState().expeditions.getIn(['expeditions', expeditionID, 'documentTypes']).toList().get(0).get('id')
    console.log('sending expedition', projectID, expeditionID, inputName)
    FKApiClient.get().postInputs(projectID, expeditionID, inputName)
      .then(res => {
        console.log('server response:', res)
        if (!res) {
          // error
          console.log('error with input registration')
        } else {
          console.log('success')
          browserHistory.push('/admin/' + projectID + '/new-expedition/confirmation')
        }
      })
  }
}

export function receiveExpeditions (projectID, expeditions, updateCurrentExpedition) {
  return function (dispatch, getState) {
    dispatch(
      {
        type: RECEIVE_EXPEDITIONS,
        projectID,
        expeditions
      }
    )
    if (updateCurrentExpedition) {
      const expeditionID = expeditions.toList().getIn([0, 'id'])
      dispatch(setCurrentExpedition(expeditionID))
    }
  }
}


/*

AUTH ACTIONS

*/


export function requestSignIn (username, password) {
  return function (dispatch, getState) {
    dispatch(loginRequest())
    console.log('requesting sign in ', username, password)
    FKApiClient.get().login(username, password)
      .then(() => {
        FKApiClient.get().onLogin()
        dispatch(loginSuccess())
        browserHistory.push('/admin')
      })
      .catch(error => {
        console.log('signin error:', error)
        if(error.response) {
          switch(error.response.status){
            case 429:
              dispatch(loginError('Try again later.'))
              break
            case 401:
              dispatch(loginError('Username or password incorrect.'))
              break
          }
        } else {
          dispatch(loginError('A server error occured.'))
        }
      })
  }
}

export function loginRequest() {
  return {
    type: LOGIN_REQUEST
  }
}

export function loginSuccess () {
  console.log('login success!')
  return {
    type: LOGIN_SUCCESS
  }
}

export function loginError (message) {
  return {
    type: LOGIN_ERROR,
    message
  }
}

export function requestSignUp (email, username, password, invite, project) {
  return function (dispatch, getState) {
    dispatch(signupRequest())

    const params = {
      'email': email,
      'username': username,
      'password': password,
      'invite': invite,
      'project': project
    }

    FKApiClient.get().register(params)
      .then(() => {
        dispatch(signupSuccess())
        browserHistory.push('/signin')
      })
      .catch(error => {
        console.log('signup error:', error)
        if (error.response && error.response.status === 400) {
          error.response.json()
            .then(err => {
              dispatch(signupError(err))
            })
        } else {
          dispatch(signupError('A server error occurred.'))
        }
      })
  }
}

export function signupRequest () {
  return {
    type: SIGNUP_REQUEST
  }
}

export function signupSuccess () {
  console.log('signup success!')
  return {
    type: SIGNUP_SUCCESS
  }
}

export function signupError (message) {
  return {
    type: SIGNUP_ERROR,
    message
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

export function jumpTo (date, expeditionID) {
  return function (dispatch, getState) {
    var state = getState()
    var expedition = state.expeditions[expeditionID]
    var expeditionDay = Math.floor((date.getTime() - expedition.start.getTime()) / (1000 * 3600 * 24))
    if (expedition.days[expeditionDay]) {
      dispatch(updateTime(date, true, expeditionID))
      return dispatch(fetchDay(date))
    } else {
      dispatch(showLoadingWheel())
      return dispatch(fetchDay(date))
    }
  }
}



export function connect () {
  browserHistory.push('/admin/okavango_16')
  return {
    type: CONNECT
  }
}



export function disconnect () {
  browserHistory.push('/')
  return {
    type: DISCONNECT
  }
}