
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
      onEnter={(state, replace, callback) => {
        checkAuthentication(state, replace)
        store.dispatch(actions.requestUser((user) => {
          store.dispatch(actions.requestProjects((projects) => {
            callback()
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
        }))
      }}
    >

      <IndexRoute
        component={ProfileSection}
        onEnter={(state) => {
          store.dispatch(actions.setBreadcrumbs(0, null))
        }}
      />
      <Route 
        path="profile"
        component={ProfileSection}
      />

      <Route 
        path="new-project" 
        component={NewProjectContainer}
        onEnter={() => {
          store.dispatch(actions.newProject())
          store.dispatch(actions.setBreadcrumbs(0, 'New Project', '/new-project'))
        }}
      />

      <Route
        path=":projectID"
        onEnter={(state, replace, callback) => {
          const projectID = state.params.projectID
          const projectName = store.getState().expeditions.getIn(['projects', projectID, 'name'])
          let expeditionID = state.params.expeditionID
          store.dispatch(actions.setCurrentProject(projectID))
          store.dispatch(actions.setBreadcrumbs(0, 'Project: ' + projectName, '/admin/' + projectID))
          store.dispatch(actions.requestExpeditions(projectID, (expeditions) => {
            callback()
            if (expeditions.size === 0) {
              browserHistory.push('/admin/' + projectID + '/new-expedition')
              store.dispatch(actions.setBreadcrumbs(1, 'New Expedition', '/admin/' + projectID + '/new-expedition'))
            } else {
              let expeditionName = ''
              if (!expeditionID) {
                expeditionID = expeditions.first().get('id')
                expeditionName = store.getState().expeditions.getIn(['expeditions', expeditionID, 'name'])
                browserHistory.push('/admin/' + projectID + '/' + expeditionID)
              } else {
                if (expeditions.some(e => e.get('id') === expeditionID)) {
                  expeditionName = store.getState().expeditions.getIn(['expeditions', expeditionID, 'name'])
                  store.dispatch(actions.setCurrentExpedition(expeditionID))
                } else {
                  expeditionID = expeditions.first().get('id')
                  expeditionName = store.getState().expeditions.getIn(['expeditions', expeditionID, 'name'])
                  browserHistory.push('/admin/' + projectID + '/' + expeditions.first().get('id'))
                }
              }
              store.dispatch(actions.setBreadcrumbs(1, 'Expedition: ' + expeditionName, '/admin/' + projectID + '/' + expeditionID))
            }
          }))
        }}
        onChange={(state) => {
          const projectID = state.params.projectID
          const projectName = store.getState().expeditions.getIn(['projects', projectID, 'name'])
          store.dispatch(actions.setBreadcrumbs(0, 'Project: ' + projectName, '/admin/' + projectID))
        }}
      >

        <Route 
          path="new-expedition" 
          onEnter={(state) => {
            store.dispatch(actions.newExpedition())
            store.dispatch(actions.setBreadcrumbs(1, 'New Expedition', '/admin/' + state.params.projectID + '/new-expedition'))
          }}
          onChange={(state) => {
            store.dispatch(actions.setBreadcrumbs(1, 'New Expedition', '/admin/' + state.params.projectID + '/new-expedition'))
          }}
        >
          <IndexRoute component={NewGeneralSettingsContainer}/>
          <Route 
            path="general-settings"
            component={NewGeneralSettingsContainer}
            onEnter={(state) => {
              store.dispatch(actions.setBreadcrumbs(2, 'General Settings', '/admin/' + state.params.projectID + '/new-expedition/general-settings'))
            }}
          />
          <Route 
            path="inputs"
            component={NewInputsContainer}
            onEnter={(state) => {
              store.dispatch(actions.setBreadcrumbs(2, 'Input Setup', '/admin/' + state.params.projectID + '/new-expedition/inputs'))
            }}
          />
          <Route 
            path="confirmation"
            component={NewConfirmationContainer}
            onEnter={(state) => {
              store.dispatch(actions.setBreadcrumbs(2, 'Confirmation', '/admin/' + state.params.projectID + '/new-expedition/confirmation'))
            }}
          />
        </Route>

        <Route path=":expeditionID" 
          onLeave={() => {
            store.dispatch(actions.setCurrentExpedition(null))
          }}
          onEnter={(state) => {
            const expeditionID = state.params.expeditionID
            const expeditionName = store.getState().expeditions.getIn(['expeditions', expeditionID, 'name'])
            store.dispatch(actions.setCurrentExpedition(expeditionID))
            store.dispatch(actions.setBreadcrumbs(1, 'Expedition: ' + expeditionName, '/admin/' + state.params.projectID + '/' + expeditionID))
          }}
          onChange={(state) => {
            const expeditionID = state.params.expeditionID
            const expeditionName = store.getState().expeditions.getIn(['expeditions', expeditionID, 'name'])
            store.dispatch(actions.setBreadcrumbs(1, 'Expedition: ' + expeditionName, '/admin/' + state.params.projectID + '/' + expeditionID))
          }}
        >
          <IndexRoute component={ExpeditionPageContainer}/>
          <Route
            path="dashboard"
            component={ExpeditionPageContainer}
            onEnter={(state) => {
              store.dispatch(actions.setBreadcrumbs(1, 'Expedition: ' + expeditionName, '/admin/' + state.params.projectID + '/' + expeditionID))
            }}
          />
          <Route
            path="general-settings"
            component={GeneralSettingsContainer}
            onEnter={(state) => {
              store.dispatch(actions.setBreadcrumbs(2, 'General Settings', '/admin/' + state.params.projectID + '/' + state.params.expeditionID + '/general-settings'))
            }}
          />
          <Route
            path="inputs"
            component={InputsContainer}
            onEnter={(state, replace, callback) => {
              store.dispatch(actions.initInputPage(() => {
                store.dispatch(actions.setBreadcrumbs(2, 'Input Setup', '/admin/' + state.params.projectID + '/' + state.params.expeditionID + '/inputs'))
                callback()
              }))
            }}
          />
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