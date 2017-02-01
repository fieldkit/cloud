
import fetch from 'whatwg-fetch'
// import { fetch } from "../vendor_modules/redux-auth"
import * as d3 from 'd3'
import { browserHistory } from 'react-router'
import FKApiClient from '../api/api.js'
import I from 'immutable'
import slug from 'slug'

export const SET_ERROR = 'SET_ERROR'

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
          dispatch([
            {
              type: RECEIVE_PROJECTS,
              projects
            },
            {
              type: SET_CURRENT_PROJECT,
              projectID: projects.toList().getIn([0, 'slug'])
            }
          ])
        }
      })
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

/*

EXPEDITION ACTIONS

*/

export const NEW_EXPEDITION = 'NEW_EXPEDITION'
export const SET_CURRENT_EXPEDITION = 'SET_CURRENT_EXPEDITION'
export const SET_EXPEDITION_PROPERTY = 'SET_EXPEDITION_PROPERTY'
export const SAVE_EXPEDITION = 'SAVE_EXPEDITION'
export const REQUEST_EXPEDITIONS = 'REQUEST_EXPEDITIONS'
export const RECEIVE_EXPEDITIONS = 'RECEIVE_EXPEDITIONS'
export const ADD_DOCUMENT_TYPE = 'ADD_DOCUMENT_TYPE'
export const REMOVE_DOCUMENT_TYPE = 'REMOVE_DOCUMENT_TYPE'
export const RECEIVE_TOKEN = 'RECEIVE_TOKEN'

export function newExpedition () {
  return function (dispatch, getState) {
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

export function requestExpeditions (callback) {
  return function (dispatch, getState) {
    const projectID = getState().expeditions.getIn(['newProject', 'id'])
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
    const projectID = getState().expeditions.getIn(['newProject', 'id'])
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

    const projectID = getState().expeditions.getIn(['newProject', 'id'])
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
    const projectID = getState().expeditions.getIn(['newProject', 'id'])
    const expedition = getState().expeditions.get('newExpedition')
    const expeditionID = expedition.get('id')
    FKApiClient.addInput(projectID, expeditionID, id)
      .then(res => {
        console.log('server response:', res)
        if (!res) {
          console.log('adding input, error')
        } else {
          console.log('input successfully added')
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
export const SIGNIN_ERROR = 'SIGNIN_ERROR'
export const SIGNUP_REQUEST = 'SIGNUP_REQUEST'
export const SIGNUP_SUCCESS = 'SIGNUP_SUCCESS'
export const SIGNUP_ERROR = 'SIGNUP_ERROR'
export const SIGNOUT_REQUEST = 'SIGNOUT_REQUEST'
export const SIGNOUT_SUCCESS = 'SIGNOUT_SUCCESS'
export const SIGNOUT_ERROR = 'SIGNOUT_ERROR'

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
    FKApiClient.signUp(params)
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
    FKApiClient.signOut()
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