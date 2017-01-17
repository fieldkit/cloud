
/*

  TODO:
  Request expeditions
  Reroute callback


  styling
  window resize

  drag and zoom interactions
  webGLOverlay rendering optimization
  favicon
  path linewidth
  change Map shouldComponentUpdate to allow update when reaching end of timeline
*/

import '../scss/app.scss'

import 'babel-polyfill'

import React from 'react'
import ReactDOM from 'react-dom'
import { createStore, applyMiddleware, combineReducers } from 'redux'
import { Provider } from 'react-redux'
import thunkMiddleware from 'redux-thunk'
import multiMiddleware from 'redux-multi'
import { batchedSubscribe } from 'redux-batched-subscribe'

import * as actions from './actions'
import expeditionReducer from './reducers/expeditions'
import authReducer from './reducers/auth'
import { Router, Route, IndexRoute, Redirect, browserHistory } from 'react-router'

import RootContainer from './containers/RootContainer'

import MapPageContainer from './containers/MapPageContainer'

import {FKApiClient} from './api/api.js';

document.getElementById('root').remove()

const createStoreWithMiddleware = applyMiddleware(
  thunkMiddleware,
  multiMiddleware,
)(createStore)
// const createStoreWithBatching = batchedSubscribe(
//   fn => fn()
// )(createStoreWithMiddleware)
const reducer = combineReducers({
  auth: authReducer,
  expeditions: expeditionReducer,
  // routing: routerReducer
})
const store = createStoreWithMiddleware(reducer)


const routes = (
  <Route path="/" component={RootContainer}>
    <IndexRoute onEnter={(nextState, replace) => {
      if (nextState.location.pathname === '/') {
        // console.log('aga', store.getState().expeditions.get('expeditions').toJS())
        // const expeditionID = store.getState().expeditions.get('expeditions').toList().get(0).get('id')
        const expeditionID = 'okavango'
        replace({
          pathname: '/' + expeditionID,
          state: { nextPathname: nextState.location.pathname }
        })
      }
    }} />
    <Route path=":expeditionID" onEnter={(state) => {
      store.dispatch(actions.requestExpedition(state.params.expeditionID))
    }}>
      <IndexRoute component={MapPageContainer} />
      <Route path="map" component={MapPageContainer}/>
    </Route>
  </Route>
)

var render = function () {
  ReactDOM.render(
    (
      <Provider store={store}>
        <Router history={browserHistory} routes={routes}/>
      </Provider>
    ),
    document.getElementById('fieldkit')
  )
}

// function onLogout () {
//   // todo
// }
// FKApiClient.setup('http://localhost:3000' || '', onLogout);
render()