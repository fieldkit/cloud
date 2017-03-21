// @flow weak

import React from 'react';
import ReactDOM from 'react-dom';
import {App} from './js/App';
import './css/index.css';

import log from 'loglevel';
log.setLevel(log.levels.TRACE);

ReactDOM.render(
  <App />,
  document.getElementById('root')
);
