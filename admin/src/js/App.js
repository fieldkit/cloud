// @flow weak

import React, { Component } from 'react';
import { BrowserRouter as Router, Route, Redirect, Switch } from 'react-router-dom'

import { FKApiClient } from './api/api';
import type { ErrorMap } from './common/util';

import { Landing } from './components/Landing';
import { Signin } from './components/unauth/Signin';
import { Signup } from './components/unauth/Signup';
import { Home } from './components/Home';

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
  state: {
    redirectTo: ?string;
  }

  constructor(props) {
    super(props);
    this.state = { redirectTo: null };

    let API_HOST = 'https://fieldkit.org';
    if (process.env.NODE_ENV === 'development') {
      API_HOST = 'http://localhost:8080';
    }

    FKApiClient.setup(API_HOST, this.onLogout.bind(this));
  }

  async signOut() {
    await FKApiClient.get().signOut();
    // TODO: how to redirect?
  }

  onLogout() {
    // TODO: how to handle this?
  }

  render() {
    return (
      <Router>
        <Switch>
          <Route exact path="/" component={Landing} />

          <Route exact path="/signin" component={Signin} />
          <Route exact path="/signup" component={Signup} />
          <Route exact path="/signout" render={() => this.signOut()} />

          <PrivateRoute path="/app" component={Home} />
        </Switch>
      </Router>
    );
  }
}
