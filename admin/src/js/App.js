// @flow weak

import React, { Component } from 'react';
import { BrowserRouter as Router, Route, Redirect, Switch } from 'react-router-dom'

import { FKApiClient } from './api/api';

import { Landing } from './components/Landing';
import { Signin } from './components/unauth/Signin';
import { Signup } from './components/unauth/Signup';

import { Main } from './components/Main';

import '../css/App.css';

const PrivateRoute = ({ component, ...rest }) => (
  <Route {...rest} render={props => (
    FKApiClient.get().signedIn() ? (
      React.createElement(component, props)
    ) : (
      <Redirect to={{
        pathname: '/landing',
        state: { from: props.location }
      }}/>
    )
  )}/>
)

export class App extends Component {
  constructor(props) {
    super(props);

    let API_HOST = 'https://fieldkit.org';
    if (process.env.NODE_ENV === 'development') {
      API_HOST = 'http://localhost:8080';
    }

    FKApiClient.setup(API_HOST, this.onUnauthorizedAccess.bind(this));
  }

  signOut() {
    FKApiClient.get().signOut();
    return <Redirect to="/landing" />;
  }

  onUnauthorizedAccess() {
    this.forceUpdate();
  }

  render() {
    return (
      <Router>
        <Switch>
          <Route exact path="/landing" component={Landing} />
          <Route exact path="/signin" component={Signin} />
          <Route exact path="/signup" component={Signup} />
          <Route exact path="/signout" render={() => this.signOut()} />

          <PrivateRoute path="/projects/:projectSlug/expeditions/:expeditionSlug/teams" component={Main} />
          <PrivateRoute path="/projects/:projectSlug/expeditions/:expeditionSlug/datasources" component={Main} />
          <PrivateRoute path="/projects/:projectSlug/expeditions/:expeditionSlug" component={Main} />
          <PrivateRoute path="/projects/:projectSlug" component={Main} />
          <PrivateRoute path="/" component={Main} />
        </Switch>
      </Router>
    );
  }
}
