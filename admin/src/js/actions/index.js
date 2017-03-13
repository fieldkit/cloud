
import fetch from 'whatwg-fetch'
// import { fetch } from "../vendor_modules/redux-auth"
import * as d3 from 'd3'
import { browserHistory } from 'react-router'
import FKApiClient from '../api/api.js'
import I from 'immutable'
import slug from 'slug'

/*

COMMON ACTIONS

*/

export const SET_ERROR = 'SET_ERROR'
export const SET_BREADCRUMBS = 'SET_BREADCRUMBS'

export function setBreadcrumbs(level, name, url) {
  return function(dispatch, getState) {
    dispatch({
      type: SET_BREADCRUMBS,
      level,
      value: !!name ? I.fromJS({
        name,
        url
      }) : null
    })
  }
}


/*

PROJECT ACTIONS

*/

export const NEW_PROJECT = 'NEW_PROJECT'
export const SET_CURRENT_PROJECT = 'SET_CURRENT_PROJECT'
export const SET_PROJECT_PROPERTY = 'SET_PROJECT_PROPERTY'
export const SAVE_PROJECT = 'SAVE_PROJECT'
export const REQUEST_PROJECTS = 'REQUEST_PROJECTS'
export const RECEIVE_PROJECTS = 'RECEIVE_PROJECTS'

export function newProject () {
  return function (dispatch, getState) {
    dispatch(
      {
        type: NEW_PROJECT
      }
    )
  }
}

export function requestProjects (callback) {
  return function (dispatch, getState) {
    FKApiClient.getProjects()
      .then(res => {
        console.log('projects received:', res)
        if (!res) {
          const projects = I.Map()
          dispatch([
            {
              type: RECEIVE_PROJECTS,
              projects
            },
            {
              type: SET_ERROR,
              errors: null
            }
          ])
          if (callback) callback(projects)
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
          dispatch([
            {
              type: RECEIVE_PROJECTS,
              projects
            },
            {
              type: SET_CURRENT_PROJECT,
              projectID: projects.toList().getIn([0, 'slug'])
            },
            {
              type: SET_ERROR,
              errors: null
            }
          ])
          if (callback) callback(projects)
        }
      })
      .catch(err => {})
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

export function setProjectProperty (keyPath, value) {
  return function (dispatch, getState) {
    dispatch({
      type: SET_PROJECT_PROPERTY,
      keyPath,
      value
    })
  }
}

export function saveProject (id, name) {
  return function (dispatch, getState) {
    FKApiClient.createProjects(name)
      .then(res => {
        console.log('project created', res)
        dispatch({
          type: SAVE_PROJECT,
          id
        })
        browserHistory.push('/admin/' + id)
      })
      .catch(error => {
        dispatch({
          type: SET_ERROR,
          errors: I.fromJS(error.message)
        })
      })
  }
}

/*

EXPEDITION ACTIONS

*/

export const NEW_EXPEDITION = 'NEW_EXPEDITION'
export const SET_CURRENT_EXPEDITION = 'SET_CURRENT_EXPEDITION'
export const SET_EXPEDITION_PROPERTY = 'SET_EXPEDITION_PROPERTY'
export const SAVE_EXPEDITION = 'SAVE_EXPEDITION'
export const REQUEST_EXPEDITIONS = 'REQUEST_EXPEDITIONS'
export const RECEIVE_EXPEDITIONS = 'RECEIVE_EXPEDITIONS'
export const ADD_INPUT = 'ADD_INPUT'
export const REMOVE_INPUT = 'REMOVE_INPUT'
export const RECEIVE_TOKEN = 'RECEIVE_TOKEN'

export function newExpedition () {
  return function (dispatch, getState) {
    console.log('wat? newexp')
    dispatch(
      {
        type: NEW_EXPEDITION
      }
    )
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

export function setExpeditionProperty (keyPath, value) {
  return function (dispatch, getState) {
    dispatch({
      type: SET_EXPEDITION_PROPERTY,
      keyPath,
      value
    })
  }
}

export function requestExpeditions (projectID, callback) {
  return function (dispatch, getState) {
    FKApiClient.getExpeditions(projectID)
      .then(res => {
        console.log('expeditions received:', res)
        if (!res) {
          const expeditions = I.Map()
          dispatch([
            {
              type: RECEIVE_EXPEDITIONS,
              projectID,
              expeditions
            },
            {
              type: SET_ERROR,
              errors: null
            }
          ])
          if (callback) callback(expeditions)
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
                inputs: [
                ]
              }))
            })
          dispatch([
            {
              type: RECEIVE_EXPEDITIONS,
              projectID,
              expeditions
            },
            {
              type: SET_ERROR,
              errors: null
            }
          ])
          if (callback) callback(expeditions)
        }
      })
      .catch(err => {})
  }
}

export function saveGeneralSettings (callback) {
  return function (dispatch, getState) {
    const projectID = getState().expeditions.getIn(['currentProject', 'id'])
    const expedition = getState().expeditions.get('currentExpedition')
    const expeditionName = expedition.get('name')
    const expeditionID = expedition.get('id')
    FKApiClient.postGeneralSettings(projectID, expeditionName)
      .then(res => {
        console.log('expeditions successfully saved:', res)
        if (!res) {
          // error
          console.log('error with expedition creation')
        } else {
          dispatch({
            type: SET_EXPEDITION_PROPERTY,
            keyPath: ['id'],
            value: expeditionID
          })
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
                    errors: null
                  }
                ])
                if (!!callback) callback()
              }
            })
            .catch(err => {})
        }
      })
      .catch(error => {
        dispatch({
          type: SET_ERROR,
          errors: I.fromJS(error.message)
        })
      })
  }
}

export function submitInputs () {
  return function (dispatch, getState) {

    const projectID = getState().expeditions.getIn(['currentProject', 'id'])
    const expedition = getState().expeditions.get('currentExpedition')
    const expeditionID = expedition.get('id')
    // TODO: submit multiple inputs
    const inputName = expedition.get('inputs').toList().get(0).get('id')
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
            errors: null
          })
        }
      })
      .catch(err => {})
  }
}

export function fetchSuggestedInputs (input, type, callback) {
  return function (dispatch, getState) {
    window.setTimeout(() => {
      const inputs = getState().expeditions
        .get('inputs')
        .filter((d) => {
          const nameCheck = d.get('name').toLowerCase().indexOf(input.toLowerCase()) > -1
          const typeCheck = d.get('type') === type
          const membershipCheck = getState().expeditions
            .getIn(['currentExpedition', 'inputs'])
            .has(d.get('id'))
          return (nameCheck) && !membershipCheck && typeCheck
        })
        .map((m) => {
          return { value: m.get('id'), label: m.get('name')}
        })
        .toArray()

      callback(null, {
        options: inputs,
        complete: true
      })
    }, 500)
  }
}

export function addInput (id, collectionType) {
  return function (dispatch, getState) {
    const projectID = getState().expeditions.getIn(['currentProject', 'id'])
    const expedition = getState().expeditions.get('currentExpedition')
    const expeditionID = expedition.get('id')
    FKApiClient.addInput(projectID, expeditionID, id)
      .then(res => {
        console.log('server response:', res)
        if (!res) {
          console.log('adding input, error')
        } else {
          console.log('input successfully added')
          dispatch({
            type: ADD_INPUT,
            id,
            collectionType,
            token: res.id
          })
        }
      })
      .catch(err => {})
  }
}

export function removeInput (id) {
  return function (dispatch, getState) {
    dispatch({
      type: REMOVE_INPUT,
      id
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

/*

AUTH ACTIONS

*/

export const SIGNIN_REQUEST = 'SIGNIN_REQUEST'
export const SIGNIN_SUCCESS = 'SIGNIN_SUCCESS'
export const SIGNUP_REQUEST = 'SIGNUP_REQUEST'
export const SIGNUP_SUCCESS = 'SIGNUP_SUCCESS'
export const SIGNOUT_REQUEST = 'SIGNOUT_REQUEST'
export const SIGNOUT_SUCCESS = 'SIGNOUT_SUCCESS'
export const RECEIVE_USER = 'RECEIVE_USER'

export function requestUser (callback) {
  return function (dispatch, getState) {
    FKApiClient.getUser()
      .then(res => {
        const user = I.fromJS({
          username: res.username,
          email: res.email
        })
        dispatch({
          type: RECEIVE_USER,
          user
        })
        callback(null, user)
      })
      .catch(err => {
        console.log('Error getting user:', err)
        callback(err, null)
      })
  }
}

export function requestSignIn (username, password) {
  return function (dispatch, getState) {
    dispatch({
      type: SIGNIN_REQUEST
    })
    FKApiClient.signIn(username, password)
      .then(() => {
        FKApiClient.onSignIn()
        dispatch({
          type: SIGNIN_SUCCESS
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
                errors: I.fromJS('Try again later.')
              })
              break
            case 401:
              dispatch({
                type: SET_ERROR,
                errors: I.fromJS('Username or password incorrect.')
              })
              break
          }
        } else {
          dispatch({
            type: SET_ERROR,
            errors: I.fromJS('A server error occured')
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
    FKApiClient.signUp(params)
      .then(() => {
        dispatch({
          type: SIGNUP_SUCCESS
        })
        browserHistory.push('/signin')
      })
      .catch(error => {
        dispatch({
          type: SET_ERROR,
          errors: I.fromJS(error.message)
        })
      })
  }
}

export function requestSignOut () {
  return function (dispatch, getState) {
    dispatch({
      type: SIGNOUT_REQUEST
    })
    FKApiClient.signOut()
      .then(() => {
        dispatch({
          type: SIGNOUT_SUCCESS
        })
        browserHistory.push('/')
      })
      .catch(error => {
        dispatch({
          type: SET_ERROR,
          errors: I.fromJS(error.message)
        })
      })
  }
}


/*

INPUT PAGE ACTIONS

*/

export const REQUEST_INPUTS = 'REQUEST_INPUTS'
export const RECEIVE_INPUTS = 'RECEIVE_INPUTS'

export function initInputPage (callback) {
  return function (dispatch, getState) {
    dispatch({
      type: REQUEST_INPUTS
    })
    const projectID = getState().expeditions.getIn(['currentProject', 'id'])
    const expeditionID = getState().expeditions.getIn(['currentExpedition', 'id'])
    FKApiClient.getInputs(projectID, expeditionID)
      .then((res) => {
        const inputs = I.fromJS(res).map((i) => {
          return i.get('slug')
        })
        dispatch({
          type: RECEIVE_INPUTS,
          expeditionID,
          inputs
        })
        callback()
      })
  }
}

