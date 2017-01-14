
/*

  TODO:
  playback
  timeline
  notifications
  identity

*/

import '../css/index.scss'

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

import Root from './components/Root'

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
  <Route path="/" component={Root}>
    <Route path=":expeditionID" onEnter={(state) => {
      store.dispatch(actions.initializeExpedition(state.params.expeditionID))
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