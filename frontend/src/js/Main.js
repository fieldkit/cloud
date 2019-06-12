// @flow weak

import React, { Component } from 'react';
import { Route, Switch } from 'react-router-dom';

import log from 'loglevel';

import Login from './pages/Login';
import SingleUserMap from './pages/SingleUserMap';
import Files from './pages/Files';
import Map from './pages/Map';

import '../css/main.css';

type Props = {
};

export class Main extends Component {
    props: Props
    state = {
    }

    constructor(props: Props) {
        super(props);

        log.setLevel('trace');
    }

    render() {
        return (
            <div className="main">
                <Switch>
                    <Route exact path={'/files'} component={Files} />
                    <Route exact path={'/files/:deviceId'} render={props => {
                        const { match } = props;
                        return ( <Files deviceId={match.params.deviceId} /> );
                    }} />
                    <Route exact path={'/map'} component={SingleUserMap} />
                    <Route exact path={'/old-map'} component={Map} />
                    <Route path={'/'} component={Login} />
                </Switch>
            </div>
        );
    }
}

export default Main;
