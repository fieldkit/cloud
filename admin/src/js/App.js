// @flow

import React, { Component } from 'react';
import {FKApiClient} from './api/api';

import logo from '../images/logo.svg';
import '../css/App.css';

export class App extends Component {
  constructor() {
    super(arguments);

    let API_HOST = 'https://fieldkit.org';
    if (process.env.NODE_ENV === 'development') {
      API_HOST = 'http://localhost:3000';
    }

    FKApiClient.setup(API_HOST, this.onLogout.bind(this));
  }

  onLogout() {
    console.log('Logged out!');
  }

  render() {
    return (
      <div className="App">
        <div className="App-header">
          <img src={logo} className="App-logo" alt="logo" />
          <h2>Welcome to React</h2>
        </div>
        <p className="App-intro">
          To get started, edit <code>src/App.js</code> and save to reload.
        </p>
      </div>
    );
  }
}
