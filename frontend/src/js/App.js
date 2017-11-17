// @flow weak

import React, { Component } from 'react';
import { BrowserRouter as Router, Route } from 'react-router-dom';
import Main from './components/Main';

import '../css/App.css';

class App extends Component {
    render() {
        return (
            <div className="App">
                <Router>
                    <Route path="/" component={ Main } />
                </Router>
            </div>
        );
    }
}

export default App;
