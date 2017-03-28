
import '../scss/app.scss'

import 'babel-polyfill'
import React from 'react'
import ReactDOM from 'react-dom'
import { createStore, applyMiddleware, combineReducers } from 'redux'
import { Provider } from 'react-redux'
import thunkMiddleware from 'redux-thunk'
import multiMiddleware from 'redux-multi'
import { Router, Route, IndexRoute, Redirect, browserHistory } from 'react-router'

import * as actions from './actions'
import expeditionReducer from './reducers/expeditions'

import RootContainer from './containers/Root/Root'
import MapPageContainer from './containers/MapPage/MapPage'
import JournalPageContainer from './containers/JournalPage/JournalPage'

import FKApiClient from './api/api.js'

document.getElementById('root').remove()

const createStoreWithMiddleware = applyMiddleware(
  thunkMiddleware,
  multiMiddleware,
)(createStore)
const reducer = combineReducers({
  expeditions: expeditionReducer
})
const store = createStoreWithMiddleware(reducer)


const routes = (
  <Route 
    path="/"
    component={RootContainer}
  >
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
      <IndexRoute
        component={ MapPageContainer }
        onEnter={(state, replace) => {
          store.dispatch(actions.setCurrentPage('map'))
        }}
      />
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