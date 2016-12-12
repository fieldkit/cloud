
/*

*/

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
import LandingPage from './components/LandingPage'
import ForgotPasswordPage from './components/ForgotPasswordPage'

import DashboardSection from './components/DashboardSection'
import UploaderSection from './components/UploaderSection'
import SourcesSection from './components/SourcesSection'
import EditorSection from './components/EditorSection'
import IdentitySection from './components/IdentitySection'
import ProfileSection from './components/ProfileSection'

import AdminPageContainer from './containers/AdminPageContainer'
import SignUpPageContainer from './containers/SignUpPageContainer'
import SignInPageContainer from './containers/SignInPageContainer'
import TeamsSectionContainer from './containers/TeamsSectionContainer'

import NewGeneralSettingsContainer from './containers/NewGeneralSettingsContainer'
import NewInputsContainer from './containers/NewInputsContainer'
import NewTeamsContainer from './containers/NewTeamsContainer'
import NewOutputsContainer from './containers/NewOutputsContainer'

import {FKApiClient} from './api/api.js';

import 'react-select/dist/react-select.css';

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
    <Route path="admin" 
      component={AdminPageContainer} 
      onEnter={requireAuth}
      onChange={(prevState, nextState, replace) => {
        const previousSection = prevState.location.pathname.split('/')[3]
        const nextSection = nextState.location.pathname.split('/')[3]
        if (previousSection === 'teams' && nextSection !== 'teams' && !!store.getState().expeditions.get('editedTeam')) {
          store.dispatch(actions.promptModalConfirmChanges(nextState.location.pathname))
          replace(prevState.location.pathname)
        } 
      }}
    >

      <IndexRoute component={ProfileSection}/>
      <Route path="profile" component={ProfileSection}/>

      <Route path="new-expedition">
        <IndexRoute component={NewGeneralSettingsContainer}/>
        <Route path="general-settings" component={NewGeneralSettingsContainer}/>
        <Route path="inputs" component={NewInputsContainer}/>
        <Route path="teams" component={NewTeamsContainer}/>
        <Route path="outputs" component={NewOutputsContainer}/>
      </Route>

      <Route path=":expeditionID">
        <IndexRoute component={DashboardSection}/>
        <Route path="dashboard" component={DashboardSection}/>
        <Route path="uploader" component={UploaderSection}/>
        <Route path="sources" component={SourcesSection}/>
        <Route path="teams" 
          component={TeamsSectionContainer} 
          onEnter={() => store.dispatch(actions.initTeamSection())}
        />
        <Route path="editor" component={EditorSection}/>
        <Route path="identity" component={IdentitySection}/>
      </Route>
    </Route>
  </Route>
)


// const history = syncHistoryWithStore(browserHistory, store)

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

FKApiClient.setup('http://localhost:3000' || '', onLogout);
render()