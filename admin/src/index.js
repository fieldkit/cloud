// @flow weak

import React from 'react';
import ReactDOM from 'react-dom';

import './css/index.css';
import {App} from './js/App';

import log from 'loglevel';
log.setLevel(log.levels.TRACE);

ReactDOM.render(
  <App />,
  document.getElementById('root')
);
