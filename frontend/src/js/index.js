

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
import { Router, Route, IndexRoute, Redirect, browserHistory } from 'react-router'

import RootContainer from './containers/Root'

import MapPageContainer from './containers/MapPage'
import JournalPageContainer from './containers/JournalPage'

import FKApiClient from './api/api.js'

document.getElementById('root').remove()

const createStoreWithMiddleware = applyMiddleware(
  thunkMiddleware,
  multiMiddleware,
)(createStore)
// const createStoreWithBatching = batchedSubscribe(
//   fn => fn()
// )(createStoreWithMiddleware)
const reducer = combineReducers({
  expeditions: expeditionReducer,
  // routing: routerReducer
})
const store = createStoreWithMiddleware(reducer)


const routes = (
  <Route path="/" component={RootContainer}>
    <IndexRoute onEnter={(nextState, replace) => {
      if (nextState.location.pathname === '/') {
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
      <IndexRoute component={ MapPageContainer } />
      <Route 
        path="map"
        component={ MapPageContainer }
        onEnter={(state, replace) => {
          store.dispatch(actions.setCurrentPage('map'))
        }}
      />
      <Route
        path="journal"
        component={ JournalPageContainer }
        onEnter={(state, replace) => {
          store.dispatch(actions.selectPlaybackMode('pause'))
          store.dispatch(actions.setCurrentPage('journal'))
        }}
      />
      <Route
        path="about"
        component={ JournalPageContainer }
        onEnter={(state, replace) => {
          store.dispatch(actions.setCurrentPage('about'))
        }}
      />
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

render()