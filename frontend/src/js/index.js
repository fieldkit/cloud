
import '../scss/app.scss'

import 'babel-polyfill'
import React from 'react'
import ReactDOM from 'react-dom'
import { createStore, applyMiddleware, combineReducers } from 'redux'
import { Provider } from 'react-redux'
import thunkMiddleware from 'redux-thunk'
import multiMiddleware from 'redux-multi'
import { Router, Route, IndexRoute, Redirect, browserHistory } from 'react-router'

import { rootSaga } from './actions/sagas'
import createSagaMiddleware from 'redux-saga'

import * as actions from './actions'
import expeditionReducer from './reducers/expeditions'

import RootContainer from './containers/Root/Root'
import MapPageContainer from './containers/MapPage/MapPage'
import JournalPageContainer from './containers/JournalPage/JournalPage'
import DataPageContainer from './containers/DataPage/DataPage'

import FKApiClient from './api/api.js'

document.getElementById('root').remove()

const sagaMiddleware = createSagaMiddleware();

const createStoreWithMiddleware = applyMiddleware(
    thunkMiddleware,
    multiMiddleware,
    sagaMiddleware,
)(createStore)

const reducer = combineReducers({
  expeditions: expeditionReducer
})

const store = createStoreWithMiddleware(reducer)

sagaMiddleware.run(rootSaga);

const routes = (
  <Route 
    path="/"
    component={RootContainer}
  >
    <IndexRoute onEnter={(nextState, replace) => {
      if (nextState.location.pathname === '/') {
        const projectID = location.hostname.split('.')[0]
        FKApiClient.getExpeditions(projectID)
          .then(resExpeditions => {
            console.log('server response received:', resExpeditions)
            if (!resExpeditions || resExpeditions.length == 0) {
              console.log('expedition data empty');
            } else {
              const expedition = resExpeditions[0];
              console.log('replacing')
              browserHistory.push(`/${expedition.slug}`)
            }
          });
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
          store.dispatch(actions.selectFocusType('expedition'))
          store.dispatch(actions.setCurrentPage('journal'))
        }}
      />
      <Route
        path="data"
        component={ DataPageContainer }
        onEnter={(state, replace) => {
          store.dispatch(actions.selectPlaybackMode('pause'))
          store.dispatch(actions.setCurrentPage('data'))
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
