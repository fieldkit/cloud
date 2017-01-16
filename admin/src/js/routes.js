
/*

  CHANGE API URL BEFORE BUILD

  indicate navigation state
  project creation flow
  load projects
  fix id issue
  new project call
  new sighting call
  mailing list form
  broken image links
  display profile information
  signup form error feedback
  change the map endpoints on ItO to call the API with https

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
import ForgotPasswordPage from './components/ForgotPasswordPage'

import UploaderSection from './components/UploaderSection'
import SourcesSection from './components/SourcesSection'
import EditorSection from './components/EditorSection'
import IdentitySection from './components/IdentitySection'
import ProfileSection from './components/ProfileSection'

import NewProjectContainer from './containers/NewProjectContainer'

import LandingPageContainer from './containers/LandingPageContainer'
import AdminPageContainer from './containers/AdminPageContainer'
import SignUpPageContainer from './containers/SignUpPageContainer'
import SignInPageContainer from './containers/SignInPageContainer'
import TeamsSectionContainer from './containers/TeamsSectionContainer'

import DashboardSectionContainer from './containers/DashboardSectionContainer'
import NewGeneralSettingsContainer from './containers/NewGeneralSettingsContainer'
import NewInputsContainer from './containers/NewInputsContainer'
import NewConfirmationContainer from './containers/NewConfirmationContainer'

// import NewTeamsContainer from './containers/NewTeamsContainer'
// import NewOutputsContainer from './containers/NewOutputsContainer'

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


function requestProjects(nextState, replace) {  
  if (!FKApiClient.get().loggedIn()) {
    replace({
      pathname: '/signin',
      state: { nextPathname: nextState.location.pathname }
    })
  }

  store.dispatch(actions.requestProjects())

  // let currentProjectID = store.getState().expeditions.get('currentProjectID')
  // if (!currentProjectID) currentProjectID = store.getState().expeditions.get('projects').toList().get(0).get('id')
  // let currentExpeditionID = store.getState().expeditions.get('currentExpeditionID')
  // if (!currentExpeditionID) currentExpeditionID = store.getState().expeditions.getIn(['projects', currentProjectID, 'expeditions']).get(0)
  // if (!currentExpeditionID) currentExpeditionID = 'new-expedition'
  // if (nextState.location.pathname === '/admin' || nextState.location.pathname === '/admin/') {
  //   replace({
  //     pathname: '/admin/' + currentProjectID + '/' + currentExpeditionID
  //   })
  // }
}



const routes = (
  <Route path="/" component={Root}>
    <IndexRoute component={LandingPageContainer}/>
    <Route component={LandingPageContainer}>
      <Route path="signup"/>
      <Route path="signin"/>
    </Route>
    <Route path="admin" 
      component={AdminPageContainer} 
      onEnter={requestProjects}
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
      <Route 
        path="profile"
        component={ProfileSection}
      />

      <Route path=":projectID" onEnter={(state) => {
        store.dispatch(actions.requestExpeditions())
      }}>
        <Route 
          path="new-expedition" 
          onEnter={() => store.dispatch(actions.initNewExpeditionSection())}
        >
          <IndexRoute component={NewGeneralSettingsContainer}/>
          <Route path="general-settings" component={NewGeneralSettingsContainer}/>
          <Route path="inputs" component={NewInputsContainer}/>
          <Route path="confirmation" component={NewConfirmationContainer}/>
          {/*
          <Route
            path="teams"
            component={NewTeamsContainer}
            onEnter={() => store.dispatch(actions.initNewTeamsSection())}
          />
          <Route path="outputs" component={NewOutputsContainer}/>
          */}
        </Route>

        <Route path=":expeditionID" onEnter={(state) => {
          store.dispatch(actions.setCurrentExpedition(state.params.expeditionID))
        }}>
          <IndexRoute component={DashboardSectionContainer}/>
          <Route path="dashboard" component={DashboardSectionContainer}/>
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

FKApiClient.setup('https://fieldkit.org')
render()