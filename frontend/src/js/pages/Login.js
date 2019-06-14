// @flow weak

import React, { Component } from 'react';
import { Redirect } from 'react-router';

import UserSession from '../api/session';

import '../../css/login.css';

class Login extends Component {
    props: Props
    state = {
    }

    isAuthenticated() {
        return new UserSession().authenticated();
    }

    async onSubmit(ev) {
        ev.preventDefault();

        try {
            const user = await new UserSession().login(this.refs.email.value, this.refs.password.value);

            if (user) {
                this.setState({
                    errors: false
                });
            }
            else {
                this.setState({
                    errors: true
                });
            }
        }
        catch (error) {
            console.log("error", error);

            this.setState({
                errors: true
            });
        }
    }

    render() {
        if (this.isAuthenticated()) {
            return <Redirect to={ "/map" } />;
        }

        return (
            <div className="map page unauth">
                <div className="logo"></div>

                <div className="contents">
                    <div className="signin">
                        <header>
                            <h1>Log In to Your Account</h1>
                        </header>

                        <form onSubmit={this.onSubmit.bind(this)}>
                            { this.state.errors && <p className="errors">
                                Email/username or password invalid. Check your information and try again.
                            </p> }
                            <div className="form-section">
                                <div className="control-group">
                                    <input ref="email" id="email" name="email" type="text" placeholder="Email" />
                                </div>
                                <div className="control-group">
                                    <input ref="password" id="password" name="password" type="password" placeholder="Password" />
                                </div>
                            </div>
                            <input type="submit" value="Log In" />
                        </form>
                    </div>
                </div>
            </div>
        );
    }
};

export default Login;
