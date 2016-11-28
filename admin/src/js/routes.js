
/*

*/

import 'babel-polyfill'

import React from 'react'
import ReactDOM from 'react-dom'
import { createStore, applyMiddleware, combineReducers } from 'redux'
import { Provider } from 'react-redux'
import thunkMiddleware from 'redux-thunk'
import { syncHistory, syncParams, routeParamsReducer } from 'react-router-redux-params'

import { fetchExpeditions } from './actions'
import okavangoReducer from './reducers'
import { Router, Route, IndexRoute, Redirect, browserHistory } from 'react-router'

import RootContainer from './containers/RootContainer'
import LandingPage from './components/LandingPage'
import SignUpPage from './components/SignUpPage'
import SignInPage from './components/SignInPage'
import ForgotPasswordPage from './components/ForgotPasswordPage'
import AdminPage from './components/AdminPage'

import DashboardSection from './components/DashboardSection'
import UploaderSection from './components/UploaderSection'
import SourcesSection from './components/SourcesSection'
import TeamsSection from './components/TeamsSection'
import EditorSection from './components/EditorSection'
import IdentitySection from './components/IdentitySection'
import ProfileSection from './components/ProfileSection'

document.getElementById('root').remove()

let store = createStore(
  okavangoReducer,
  applyMiddleware(
    thunkMiddleware
  )
)

const routes = (
  <Route path="/" component={RootContainer}>
    <IndexRoute component={LandingPage}/>
    <Route path="signup" component={SignUpPage}/>
    <Route path="signin" component={SignInPage}/>
    <Route path="forgot" component={ForgotPasswordPage}/>
    <Route path="admin" component={AdminPage}>
      <Route path="profile" component={ProfileSection}/>
      <Route path=":expeditionID">
        <IndexRoute component={DashboardSection}/>
        <Route path="dashboard" component={DashboardSection}/>
        <Route path="uploader" component={UploaderSection}/>
        <Route path="sources" component={SourcesSection}/>
        <Route path="teams" component={TeamsSection}/>
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
render()
// // store.dispatch(fetchExpeditions(getURLParameters()))

// window.onclick = function (event) {
//   if (!event.target.matches('.dropbtn')) {
//     var dropdowns = document.getElementsByClassName('dropdown-content')
//     var i
//     for (i = 0; i < dropdowns.length; i++) {
//       var openDropdown = dropdowns[i]
//       if (openDropdown.classList.contains('show')) {
//         openDropdown.classList.remove('show')
//       }
//     }
//   }
// }

