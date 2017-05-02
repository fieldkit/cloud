// @flow weak

import React, { Component } from 'react';
import { BrowserRouter as Router, Route, Redirect, Switch } from 'react-router-dom'
import { FKApiClient } from './api/api';
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
        pathname: '/signin',
        state: { from: props.location }
      }}/>
    )
  )}/>
)

export class App extends Component {
  constructor(props) {
    super(props);

    let API_HOST = 'https://api.fieldkit.org';
    if (process.env.NODE_ENV === 'development') {
      API_HOST = 'http://localhost:8080';
    } else if (process.env.NODE_ENV === 'staging' || window.location.hostname.endsWith('fieldkit.team')) {
      API_HOST = 'https://api.fieldkit.team';
    }

    FKApiClient.setup(API_HOST, this.onUnauthorizedAccess.bind(this));
  }

  signOut() {
    FKApiClient.get().signOut();
    return <Redirect to="/signin" />;
  }

  onUnauthorizedAccess() {
    this.forceUpdate();
  }

  render() {
    return (
      <Router basename="/admin">
        <Switch>
          <Route exact path="/signin" component={Signin} />
          <Route exact path="/signup" component={Signup} />
          <Route exact path="/signout" render={() => this.signOut()} />
          <PrivateRoute path="/projects/:projectSlug/expeditions/:expeditionSlug" component={Main} />
          <PrivateRoute path="/projects/:projectSlug" component={Main} />
          <PrivateRoute path="/" component={Main} />
        </Switch>
      </Router>
    );
  }
}
