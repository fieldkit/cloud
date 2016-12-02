
/*

*/

import 'babel-polyfill'

import React from 'react'
import ReactDOM from 'react-dom'
import { createStore, applyMiddleware, combineReducers } from 'redux'
import { Provider } from 'react-redux'
import thunkMiddleware from 'redux-thunk'
// import { syncHistory, syncParams, routeParamsReducer } from 'react-router-redux-params'
// import { routerReducer, syncHistoryWithStore } from 'react-router-redux'

import { fetchExpeditions } from './actions'
import expeditionReducer from './reducers/expeditions'
import authReducer from './reducers/auth'
import { Router, Route, IndexRoute, Redirect, browserHistory } from 'react-router'

import Root from './components/Root'
import LandingPage from './components/LandingPage'
import ForgotPasswordPage from './components/ForgotPasswordPage'
import AdminPage from './components/AdminPage'

import DashboardSection from './components/DashboardSection'
import UploaderSection from './components/UploaderSection'
import SourcesSection from './components/SourcesSection'
import EditorSection from './components/EditorSection'
import IdentitySection from './components/IdentitySection'
import ProfileSection from './components/ProfileSection'

import SignUpPageContainer from './containers/SignUpPageContainer'
import SignInPageContainer from './containers/SignInPageContainer'
import TeamsSectionContainer from './containers/TeamsSectionContainer'

import {FKApiClient} from './api/api.js';

document.getElementById('root').remove()

let store = createStore(
  combineReducers({
    auth: authReducer,
    expeditions: expeditionReducer
  }),
  applyMiddleware(
    thunkMiddleware
  )
)

function requireAuth(nextState, replace): void {
  // temporary as we're setting up auth
  return
  if (!FKApiClient.get().loggedIn()) {
    replace({
      pathname: '/signup',
      state: { nextPathname: nextState.location.pathname }
    })
  }
}

function onLogout () {
  // todo
}

const routes = (
  <Route path="/" component={Root}>
    <IndexRoute component={LandingPage}/>
    <Route path="signup" component={SignUpPageContainer}/>
    <Route path="signin" component={SignInPageContainer}/>
    <Route path="forgot" component={ForgotPasswordPage}/>
    <Route path="admin" component={AdminPage} onEnter={requireAuth}>
      <IndexRoute component={ProfileSection}/>
      <Route path="profile" component={ProfileSection}/>
      <Route path=":expeditionID">
        <IndexRoute component={DashboardSection}/>
        <Route path="dashboard" component={DashboardSection}/>
        <Route path="uploader" component={UploaderSection}/>
        <Route path="sources" component={SourcesSection}/>
        <Route path="teams" component={TeamsSectionContainer}/>
        <Route path="editor" component={EditorSection}/>
        <Route path="identity" component={IdentitySection}/>
      </Route>
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

store.subscribe(render)
FKApiClient.setup('http://localhost:3000' || '', onLogout);
render()