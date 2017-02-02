
import '../css/index.scss'
import 'react-select/dist/react-select.css';

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
import ProfileSection from './components/AdminPage/ProfileSection'
import LandingPageContainer from './containers/LandingPage/LandingPage'
import AdminPageContainer from './containers/AdminPage/AdminPage'
import NewProjectContainer from './containers/AdminPage/NewProjectPage/NewProjectPage'
import NewGeneralSettingsContainer from './containers/AdminPage/NewExpeditionPage/NewExpeditionPage'
import NewInputsContainer from './containers/AdminPage/NewExpeditionPage/InputsSection'
import NewConfirmationContainer from './containers/AdminPage/NewExpeditionPage/ConfirmationSection'
import ExpeditionPageContainer from './containers/AdminPage/ExpeditionPage/ExpeditionPage'
import GeneralSettingsContainer from './containers/AdminPage/ExpeditionPage/GeneralSettingsSection'
import InputsContainer from './containers/AdminPage/ExpeditionPage/InputsSection'

import FKApiClient from './api/api.js'

document.getElementById('root').remove()

const createStoreWithMiddleware = applyMiddleware(
  thunkMiddleware,
  multiMiddleware,
)(createStore)
const reducer = combineReducers({
  auth: authReducer,
  expeditions: expeditionReducer,
})
const store = createStoreWithMiddleware(reducer)


function checkAuthentication(state, replace) {  
  if (!FKApiClient.signedIn()) {
    replace({
      pathname: '/signin',
      state: { nextPathname: state.location.pathname }
    })
  } 
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
      onEnter={(state, replace) => {
        checkAuthentication(state, replace)
        store.dispatch(actions.requestProjects((projects) => {
          // browserHistory.push('/admin/new-project')
          const projectID = state.params.projectID
          if (projects.size === 0) {
            browserHistory.push('/admin/new-project')
          } else {
            if (!projectID) {
              browserHistory.push('/admin/' + projects.first().get('id'))
            } else {
              if (projects.some(p => p.get('id') === projectID)) {
                store.dispatch(actions.setCurrentProject(projectID))
              } else {
                browserHistory.push('/admin/' + projects.first().get('id'))
              }
            }
          }
        }))
      }}
    >

      <IndexRoute component={ProfileSection}/>
      <Route 
        path="profile"
        component={ProfileSection}
      />

      <Route 
        path="new-project" 
        component={NewProjectContainer}
        onEnter={() => store.dispatch(actions.newProject())}
      />

      <Route path=":projectID" onEnter={(state) => {
        const projectID = state.params.projectID
        const expeditionID = state.params.expeditionID
        store.dispatch(actions.setCurrentProject(projectID))
        store.dispatch(actions.requestExpeditions(projectID, (expeditions) => {
          if (expeditions.size === 0) {
            browserHistory.push('/admin/' + projectID + '/new-expedition')
          } else {
            if (!expeditionID) {
              browserHistory.push('/admin/' + projectID + '/' + expeditions.first().get('id'))
            } else {
              if (expeditions.some(e => e.get('id') === expeditionID)) {
                store.dispatch(actions.setCurrentExpedition(expeditionID))
              } else {
                browserHistory.push('/admin/' + projectID + '/' + expeditions.first().get('id'))
              }
            }
          }
        }))
      }}>

        <Route 
          path="new-expedition" 
          onEnter={() => {
            store.dispatch(actions.newExpedition())
          }}
        >
          <IndexRoute component={NewGeneralSettingsContainer}/>
          <Route path="general-settings" component={NewGeneralSettingsContainer}/>
          <Route path="inputs" component={NewInputsContainer}/>
          <Route path="confirmation" component={NewConfirmationContainer}/>
        </Route>

        <Route path=":expeditionID" 
          onLeave={() => {
            store.dispatch(actions.setCurrentExpedition(null))
          }}
          onEnter={(state) => {
            store.dispatch(actions.setCurrentExpedition(state.params.expeditionID))
          }}
        >
          <IndexRoute component={ExpeditionPageContainer}/>
          <Route path="dashboard" component={ExpeditionPageContainer}/>
          <Route path="general-settings" component={GeneralSettingsContainer}/>
          <Route path="inputs" component={InputsContainer}/>
        </Route>
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

render()