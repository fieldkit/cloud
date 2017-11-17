// @flow weak

import React, { Component, PropTypes } from 'react';
import { BrowserRouter as Router, Route } from 'react-router-dom';
import Main from './components/Main';

import { Provider } from 'react-redux'

import '../css/App.css';

class App extends Component {
    props: {
        state: PropTypes.object.isRequired
    }

    render() {
        const { store } = this.props

        return (
            <div className="App">
                <Provider store={store}>
                    <Router>
                        <Route path="/" component={ Main } />
                    </Router>
                </Provider>
            </div>
        );
    }
}

export default App;
