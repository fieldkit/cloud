// @flow weak

import React, { Component } from 'react';
import { connect } from 'react-redux';
import { Redirect } from 'react-router';

import log from 'loglevel';

import UserSession from '../api/session';

import '../../css/login.css';

class Login extends Component {
    props: Props

    state = {
        authenticated: null
    }

    async componentDidMount() {
        const authenticated = await new UserSession().authenticated();

        this.setState({
            authenticated
        });
    }

    async onSubmit(ev) {
        ev.preventDefault();

        const user = await new UserSession().login(this.refs.email.value, this.refs.password.value);

        log.info(user);
    }

    render() {
        const { authenticated } = this.state;

        if (authenticated) {
            return <Redirect to={ "/map" } />;
        }

        return (
            <div className="map page unauth">
                <div className="contents">
                    <div className="signin">
                        <header>
                            <h1>Sign in</h1>
                        </header>
                        <form onSubmit={this.onSubmit.bind(this)}>
                            { this.state.errors && <p className="errors">
                                Email/username or password invalid. Check your information and try again.
                            </p> }
                            <div className="content">
                                <div>
                                    <label htmlFor="email">E-Mail</label>
                                    <input ref="email" id="email" name="email" type="text" placeholder="" />
                                </div>
                                <div>
                                    <label htmlFor="password">Password</label>
                                    <input ref="password" id="password" name="password" type="password" placeholder="" />
                                </div>
                            </div>
                            <input type="submit" value="Submit" />
                        </form>
                    </div>
                </div>
            </div>
        );
    }
};

const mapStateToProps = state => ({
});

export default connect(mapStateToProps, {
})(Login);
