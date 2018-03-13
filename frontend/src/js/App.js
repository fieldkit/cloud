// @flow weak

import React, { Component } from 'react';
import { BrowserRouter as Router, Route } from 'react-router-dom';
import Main from './Main';

import { Provider } from 'react-redux';

import '../css/App.css';

type Props = {
    store: mixed,
}

class App extends Component {
    props: Props

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
