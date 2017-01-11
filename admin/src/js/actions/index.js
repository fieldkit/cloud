
import fetch from 'whatwg-fetch'
// import { fetch } from "../vendor_modules/redux-auth"
import * as d3 from 'd3'
import { browserHistory } from 'react-router'
import {FKApiClient} from '../api/api.js'

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

AUTH ACTIONS

*/

export function requestSignIn (email, password) {
  return function (dispatch, getState) {
    dispatch(loginRequest())
    FKApiClient.get().login(email, password)
      .then(() => {
        dispatch(loginSuccess())
        browserHistory.push('/admin/okavango_16')
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

export function requestSignUp (email, userName, firstName, lastName, password) {
  return function (dispatch, getState) {
    dispatch(signupRequest())

    const params = {
      'email': email,
      'username': userName,
      'first_name': firstName,
      'last_name': lastName,
      'password': password
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


// function timestampToString (t) {
//   var d = new Date(t)
//   var year = d.getUTCFullYear()
//   var month = (d.getUTCMonth() + 1) + ''
//   if (month.length === 1) month = '0' + month
//   var date = (d.getUTCDate()) + ''
//   if (date.length === 1) date = '0' + date
//   return year + '-' + month + '-' + date
// }

// export const CLOSE_LIGHTBOX = 'CLOSE_LIGHTBOX'

// export function closeLightBox () {
//   return {
//     type: CLOSE_LIGHTBOX
//   }
// }

// export const SHOW_360_PICTURE = 'SHOW_360_PICTURE'

// export function show360Picture (post) {
//   return {
//     type: SHOW_360_PICTURE,
//     post
//   }
// }

// export const ENABLE_CONTENT = 'ENABLE_CONTENT'

// export function enableContent () {
//   return {
//     type: ENABLE_CONTENT
//   }
// }

// export const SET_PAGE = 'SET_PATH'

// export function setPage () {
//   return (dispatch, getState) => {
//     dispatch({
//       type: SET_PAGE,
//       location: location.pathname
//     })
//     if (location.pathname.indexOf('/journal') > -1) dispatch(checkFeedContent())
//   }
// }

// export function checkFeedContent () {
//   return (dispatch, getState) => {
//     const state = getState()
//     const expeditionID = state.selectedExpedition
//     const expedition = state.expeditions[expeditionID]
//     const dayCount = expedition.dayCount
//     const posts = d3.values(expedition.features)
//     const postsByDay = expedition.postsByDay
//     const contentHeight = d3.select('#content').node().offsetHeight
//     const scrollTop = d3.select('#content').node().scrollTop
//     const feedElement = d3.select('#feed').node()
//     const viewRange = [scrollTop, scrollTop + contentHeight]

//     if (feedElement) {
//       const postElements = d3.select(feedElement).selectAll('div.post')._groups[0]
//       var visibleDays = []
//       var visibleElements = []
//       if (postElements) {
//         postElements.forEach((p) => {
//           var postRange = [p.offsetTop - 100, p.offsetTop + p.offsetHeight - 100]
//           if ((viewRange[0] > postRange[0] && viewRange[0] <= postRange[1]) || (viewRange[0] <= postRange[0] && viewRange[1] > postRange[0]) || (viewRange[1] > postRange[0] && viewRange[1] <= postRange[1])) {
//             visibleElements.push(p.className.split(' ')[1])
//           }
//         })
//       }
//       visibleElements.forEach(p => {
//         var feature = expedition.features[p]
//         var day = Math.floor((new Date(feature.properties.DateTime).getTime() - expedition.start.getTime()) / (1000 * 3600 * 24))
//         if (visibleDays.indexOf(day) === -1) visibleDays.push(day)
//       })
//       for (var i = 0; i < visibleDays.length - 1; i++) {
//         if (Math.abs(visibleDays[i] - visibleDays[i + 1])) {
//           dispatch(fetchPostsByDay(expeditionID, null, i))
//           break
//         }
//       }

//       const feedHeight = feedElement.offsetHeight
//       if ((posts.length === 0) || feedHeight < contentHeight || (scrollTop <= 100 && !postsByDay[dayCount]) || (scrollTop >= feedHeight - contentHeight - 100 && !postsByDay[0])) {
//         dispatch(fetchPostsByDay(expeditionID, expedition.currentDate))
//       }
//     } else {
//       if ((posts.length === 0) || (scrollTop <= 100 && !postsByDay[dayCount])) {
//         dispatch(fetchPostsByDay(expeditionID, expedition.currentDate))
//       }
//     }
//   }
// }

// export const FETCH_POSTS_BY_DAY = 'FETCH_POSTS_BY_DAY'

// export function fetchPostsByDay (_expeditionID, date, expeditionDay) {
//   return function (dispatch, getState) {
//     var i
//     var state = getState()
//     // if (state.isFetchingPosts > 0) return
//     var expeditionID = _expeditionID || state.selectedExpedition
//     var expedition = state.expeditions[expeditionID]
//     if (!expeditionDay) {
//       if (!date) date = expedition.currentDate
//       expeditionDay = Math.floor((date.getTime() - expedition.start.getTime()) / (1000 * 3600 * 24))
//     }

//     var daysToFetch = []

//     if (!expedition.postsByDay[expeditionDay]) daysToFetch.push(expeditionDay)
//     else if (expedition.postsByDay[expeditionDay] === 'loading') return
//     else {
//       for (i = expeditionDay - 1; i >= 0; i--) {
//         if (expedition.postsByDay[i] === 'loading') break
//         if (!expedition.postsByDay[i]) {
//           daysToFetch.push(i)
//           daysToFetch[0] = i
//           break
//         }
//       }
//       for (i = expeditionDay + 1; i < expedition.dayCount; i++) {
//         if (expedition.postsByDay[i] === 'loading') break
//         if (!expedition.postsByDay[i]) {
//           daysToFetch.push(i)
//           break
//         }
//       }
//     }

//     if (daysToFetch.length === 0) return
//     const datesToFetch = []
//     daysToFetch.forEach(function (d, i) {
//       var t = expedition.start.getTime() + d * (1000 * 3600 * 24)
//       datesToFetch[i] = t
//     })
//     var range = [
//       timestampToString(d3.min(datesToFetch)),
//       timestampToString(d3.max(datesToFetch) + (1000 * 3600 * 24))
//     ]

//     dispatch({
//       type: FETCH_POSTS_BY_DAY,
//       expeditionID: expeditionID,
//       daysToFetch
//     })

//     var queryString = 'https://intotheokavango.org/api/features?limit=0&FeatureType=blog,audio,image,tweet&limit=0&Expedition=' + state.selectedExpedition + '&startDate=' + range[0] + '&endDate=' + range[1]
//     // console.log('querying posts:', queryString)
//     fetch(queryString)
//       .then(response => response.json())
//       .then(json => {
//         var results = json.results.features
//         // console.log('done with post query! Received:' + results.length + ' features.')
//         return dispatch(receivePosts(expeditionID, results, range))
//       })
//       .then(() => {
//         dispatch(checkFeedContent())
//       })
//   }
// }

// export const RECEIVE_POSTS = 'RECEIVE_POSTS'

// export function receivePosts (expeditionID, data, timeRange) {
//   return {
//     type: RECEIVE_POSTS,
//     expeditionID,
//     data,
//     timeRange
//   }
// }


// export const COMPLETE_DAYS = 'COMPLETE_DAYS'

// export function completeDays (expeditionID) {
//   return {
//     type: COMPLETE_DAYS,
//     expeditionID
//   }
// }

// export const SHOW_LOADING_WHEEL = 'SHOW_LOADING_WHEEL'

// export function showLoadingWheel () {
//   return {
//     type: SHOW_LOADING_WHEEL
//   }
// }

// export const HIDE_LOADING_WHEEL = 'HIDE_LOADING_WHEEL'

// export function hideLoadingWheel () {
//   return {
//     type: HIDE_LOADING_WHEEL
//   }
// }

// export const JUMP_TO = 'JUMP_TO'

// export function jumpTo (date, expeditionID) {
//   return function (dispatch, getState) {
//     var state = getState()
//     var expedition = state.expeditions[expeditionID]
//     var expeditionDay = Math.floor((date.getTime() - expedition.start.getTime()) / (1000 * 3600 * 24))
//     if (expedition.days[expeditionDay]) {
//       dispatch(updateTime(date, true, expeditionID))
//       return dispatch(fetchDay(date))
//     } else {
//       dispatch(showLoadingWheel())
//       return dispatch(fetchDay(date))
//     }
//   }
// }

// export const START = 'START'

// export function startAnimation () {
//   return {
//     type: START
//   }
// }

// export const REQUEST_EXPEDITIONS = 'REQUEST_EXPEDITIONS'

// export function requestExpeditions () {
//   return {
//     type: REQUEST_EXPEDITIONS
//   }
// }

// export const UPDATE_TIME = 'UPDATE_TIME'

// export function updateTime (currentDate, updateMapState, expeditionID) {
//   return {
//     type: UPDATE_TIME,
//     currentDate,
//     updateMapState,
//     expeditionID
//   }
// }

// export const UPDATE_MAP = 'UPDATE_MAP'

// export function updateMap (currentDate, coordinates, viewGeoBounds, zoom, expeditionID) {
//   return function (dispatch, getState) {
//     var state = getState()
//     var expedition = state.expeditions[expeditionID]
//     var tiles = expedition.featuresByTile
//     var tileResolution = Math.floor((expedition.geoBounds[2] - expedition.geoBounds[0]) * 111 / 10)

//     const coordinatesToTile = (coordinates, geoBounds) => {
//       var x = Math.floor((coordinates[0] - geoBounds[0]) * 111 / 10)
//       var y = Math.floor((coordinates[1] - geoBounds[3]) * 111 / 10)
//       return {x, y}
//     }

//     const tileToCoordinates = (tile, geoBounds) => {
//       var lng = (tile.x * 10 / 111) + geoBounds[0]
//       var lat = (tile.y * 10 / 111) + geoBounds[3]
//       return [lng, lat]
//     }

//     var west = viewGeoBounds[0]
//     var north = viewGeoBounds[1]
//     var east = viewGeoBounds[2]
//     var south = viewGeoBounds[3]

//     // TODO TEMPORARY: limiting max range
//     var centroid = [(west + east) / 2, (south + north) / 2]
//     west = centroid[0] + Math.max(west - centroid[0], -0.1)
//     east = centroid[0] + Math.min(east - centroid[0], 0.1)
//     north = centroid[1] + Math.min(north - centroid[1], 0.1)
//     south = centroid[1] + Math.max(south - centroid[1], -0.1)
//     // TEMPORARY

//     var northWestTile = coordinatesToTile([west, north], expedition.geoBounds)
//     var southEastTile = Object.assign({}, northWestTile)
//     while (tileToCoordinates(southEastTile, expedition.geoBounds)[0] <= east) {
//       southEastTile.x++
//     }
//     while (tileToCoordinates(southEastTile, expedition.geoBounds)[1] >= south) {
//       southEastTile.y--
//     }

//     var tileRange = []
//     var tilesInView = []
//     for (var x = northWestTile.x; x <= southEastTile.x; x++) {
//       for (var y = northWestTile.y; y >= southEastTile.y; y--) {
//         var tile = x + y * tileResolution
//         if (!tiles[tile]) tileRange.push({x, y})
//         tilesInView.push(x + y * tileResolution)
//       }
//     }

//     var queryNorthWest = [180, -90]
//     var querySouthEast = [-180, 90]
//     tileRange.forEach((t) => {
//       var northWest = tileToCoordinates(t, expedition.geoBounds)
//       var southEast = tileToCoordinates({x: t.x + 1, y: t.y - 1}, expedition.geoBounds)
//       if (queryNorthWest[0] > northWest[0]) queryNorthWest[0] = northWest[0]
//       if (queryNorthWest[1] < northWest[1]) queryNorthWest[1] = northWest[1]
//       if (querySouthEast[0] < southEast[0]) querySouthEast[0] = southEast[0]
//       if (querySouthEast[1] > southEast[1]) querySouthEast[1] = southEast[1]
//     })
//     var queryGeoBounds = [queryNorthWest[0], queryNorthWest[1], querySouthEast[0], querySouthEast[1]]

//     tileRange.forEach((t, i, a) => {
//       a[i] = t.x + t.y * tileResolution
//     })

//     if (tileRange.length > 0) {
//       var queryString = 'https://intotheokavango.org/api/features?limit=0&FeatureType=blog,audio,image,tweet,sighting&Expedition=' + state.selectedExpedition + '&geoBounds=' + queryGeoBounds.toString()
//       // console.log('querying features by tile:', queryString)
//       fetch(queryString)
//         .then(response => response.json())
//         .then(json => {
//           var results = json.results.features
//           // console.log('done with feature query! Received ' + results.length + ' features.')
//           dispatch(receiveFeatures(state.selectedExpedition, results, tileRange))
//         })
//     }

//     return dispatch({
//       type: UPDATE_MAP,
//       expeditionID,
//       currentDate,
//       coordinates,
//       viewGeoBounds,
//       tilesInView,
//       zoom,
//       tileRange
//     })
//   }
// }

// export const RECEIVE_EXPEDITIONS = 'RECEIVE_EXPEDITIONS'

// export function receiveExpeditions (data) {
//   for (var k in data.results) {
//     if (data.results[k].Days < 1) {
//       delete data.results[k]
//     }
//   }

//   return {
//     type: RECEIVE_EXPEDITIONS,
//     data
//   }
// }

// export function fetchExpeditions (parameters) {
//   return function (dispatch, getState) {

//     const startDate = parameters.date ? new Date(parseInt(parameters.date)) : new Date('2016-08-20 09:30:00+00:00') 

//     dispatch(requestExpeditions())
//     return fetch('https://intotheokavango.org/api/expeditions')
//       .then(response => response.json())
//       .then(json => {
//         return {
//           ...json,
//           results: {
//             ...json.results,
//             okavango_16: {
//               ...json.results.okavango_16,
//               Days: 18
//             }
//           }
//         }
//       })
//       .then(json => dispatch(receiveExpeditions(json)))
//       .then(() => dispatch(fetchDay(startDate, null, null, true)))
//       .then(() => {
//         var state = getState()
//         dispatch(fetchTotalSightings(state.selectedExpedition))
//         if (location.pathname.indexOf('/journal') > -1) {
//           dispatch(checkFeedContent())
//         }
//       })
//   }
// }

// export function fetchTotalSightings (id) {
//   return function (dispatch, getState) {
//     return fetch('https://intotheokavango.org/api/sightings?FeatureType=sighting&limit=0&Expedition=' + id)
//       .then(response => response.json())
//       .then(json => {
//         if (id !== 'okavango_16') {
//           return json
//         } else {
//           return {
//             ...json,
//             results: json.results.slice(0, 20)
//           }
//         }
//       })
//       .then(json => dispatch(receiveTotalSightings(id, json)))
//   }
// }

// export const RECEIVE_TOTAL_SIGHTINGS = 'RECEIVE_TOTAL_SIGHTINGS'

// export function receiveTotalSightings (id, data) {
//   return {
//     type: RECEIVE_TOTAL_SIGHTINGS,
//     id,
//     data
//   }
// }

// export function fetchDay (date, initialDate, _expeditionID, initialize) {
//   if (!initialDate) initialDate = date
//   return function (dispatch, getState) {
//     var state = getState()
//     var expeditionID = _expeditionID || state.selectedExpedition
//     var expedition = state.expeditions[expeditionID]
//     if (!date) date = expedition.currentDate
//     var expeditionDay = Math.floor((date.getTime() - expedition.start.getTime()) / (1000 * 3600 * 24))
//     var daysToFetch = []
//     if (!expedition.days[expeditionDay - 1] && expeditionDay - 1 >= 0) daysToFetch.push(expeditionDay - 1)
//     if (!expedition.days[expeditionDay]) daysToFetch.push(expeditionDay)
//     if (!expedition.days[expeditionDay + 1] && expeditionDay + 1 < expedition.dayCount) daysToFetch.push(expeditionDay + 1)

//     if (daysToFetch.length === 0) return

//     daysToFetch.forEach(function (d, i, a) {
//       var t = expedition.start.getTime() + d * (1000 * 3600 * 24)
//       a[i] = t
//     })
//     var range = [
//       timestampToString(d3.min(daysToFetch)),
//       timestampToString(d3.max(daysToFetch) + (1000 * 3600 * 24))
//     ]

//     const goFetch = (featureTypes, results, expeditionID) => {
//       var type = featureTypes.shift()
//       var queryString = 'https://intotheokavango.org/api/features?limit=0&FeatureType=' + type + '&Expedition=' + expeditionID + '&startDate=' + range[0] + '&endDate=' + range[1]
//       if (type === 'ambit_geo') queryString += '&resolution=2'
//       // console.log('querying:', queryString)
//       fetch(queryString)
//         .then(response => response.json())
//         .then(json => {
//           results = results.concat(json.results.features)
//           if (featureTypes.length > 0) {
//             // console.log('received ' + json.results.features.length + ' ' + type)
//             goFetch(featureTypes, results, expeditionID)
//           } else {
//             // console.log('done with query! Received ' + json.results.features.length + ' ' + type, initialize)
//             dispatch(receiveDay(expeditionID, results, range))
//             dispatch(completeDays(expeditionID))
//             var state = getState()
//             var expedition = state.expeditions[state.selectedExpedition]
//             var days = expedition.days
//             var incompleteDays = []
//             d3.keys(expedition.days).forEach((k) => {
//               if (expedition.days[k].incomplete) incompleteDays.push(k)
//             })
//             if (incompleteDays.length === 0) {
//               if (!state.animate && initialize) dispatch(startAnimation())
//               dispatch(updateTime(initialDate, false, expeditionID))
//               dispatch(hideLoadingWheel())
//             } else {
//               // console.log('incomplete days', incompleteDays)
//               var nextTarget = -1
//               for (var i = 0; i < incompleteDays.length; i++) {
//                 var id = incompleteDays[i]
//                 for (var j = Math.max(0, id - 1); j < expedition.dayCount; j++) {
//                   if (!days[j]) {
//                     nextTarget = j
//                     break
//                   }
//                 }
//                 if (nextTarget > -1) break
//               }
//               if (nextTarget > -1) {
//                 nextTarget = new Date(expedition.start.getTime() + nextTarget * (1000 * 3600 * 24))
//                 dispatch(fetchDay(nextTarget, initialDate, expeditionID, initialize))
//               }
//             }
//           }
//         })
//     }
//     goFetch(['ambit_geo', 'beacon'], [], expeditionID)
//   }
// }

// export const SET_EXPEDITION = 'SET_EXPEDITION'

// export function setExpedition (id) {
//   return function (dispatch, getState) {
//     var state = getState()
//     var expedition = state.expeditions[id]
//     if (expedition.totalSightings.length === 0) {
//       dispatch(fetchTotalSightings(id))
//     }
//     dispatch({
//       type: SET_EXPEDITION,
//       id
//     })
//   }
// }

// export const SET_CONTROL = 'SET_CONTROL'

// export function setControl (target, mode) {
//   return {
//     type: SET_CONTROL,
//     target,
//     mode
//   }
// }

// export const REQUEST_DAY = 'REQUEST_DAY'

// export function requestDay (expeditionID, dayID) {
//   return {
//     type: REQUEST_DAY,
//     expeditionID,
//     dayID
//   }
// }

// export const RECEIVE_DAY = 'RECEIVE_DAY'

// export function receiveDay (expeditionID, data, dateRange) {
//   return {
//     type: RECEIVE_DAY,
//     expeditionID,
//     data,
//     dateRange
//   }
// }

// export const RECEIVE_FEATURES = 'RECEIVE_FEATURES'

// export function receiveFeatures (expeditionID, data, tileRange) {
//   return {
//     type: RECEIVE_FEATURES,
//     expeditionID,
//     data,
//     tileRange
//   }
// }

// export const SELECT_FEATURE = 'SELECT_FEATURE'

// export function selectFeature () {
//   return {
//     type: SELECT_FEATURE
//   }
// }

// export const UNSELECT_FEATURE = 'UNSELECT_FEATURE'

// export function unselectFeature () {
//   return {
//     type: UNSELECT_FEATURE
//   }
// }
