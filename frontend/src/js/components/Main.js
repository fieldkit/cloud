// @flow weak

import React, { Component } from 'react';
import { connect } from 'react-redux'
import { Route, Switch, NavLink } from 'react-router-dom';

import log from 'loglevel';
import Map from './pages/Map';
import Sensors from './pages/Sensors';
import About from './pages/About';

import type { ActiveExpedition  } from '../types';

import '../../css/main.css';

type Props = {
    activeExpedition: ActiveExpedition
};

export class Main extends Component {
    props: Props

    constructor(props: Props) {
        super(props);

        log.setLevel('trace');
    }

    render() {
        const { activeExpedition } = this.props;

        return (
            <div className="main">
                <div className="header">
                    <div className="project-name">
                        { activeExpedition && activeExpedition.project && activeExpedition.project.name }
                    </div>
                    <div className="nav-bar">
                        <div className="navigation-tabs">
                            <NavLink exact to={ '/' }><span>Map</span></NavLink>
                            <NavLink exact to={ '/sensors' }><span>Sensors</span></NavLink>
                            <NavLink exact to={ '/journal' }><span>Journal</span></NavLink>
                            <NavLink exact to={ '/about' }><span>About</span></NavLink>
                        </div>
                    </div>
                </div>
                <Switch>
                    <Route path={ '/sensors' } component={ Sensors } />
                    <Route path={ '/about' } component={ About } />
                    <Route path={ '/' } component={ Map } />
                </Switch>
            </div>
        );
    }
}

const mapStateToProps = state => ({
    activeExpedition: state.activeExpedition
});

export default connect(mapStateToProps, {
})(Main);
