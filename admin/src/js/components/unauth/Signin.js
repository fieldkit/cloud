// @flow weak

import React, { Component } from 'react'
import { Redirect } from 'react-router'

import log from 'loglevel';

import { FKApiClient } from '../../api/api';
import { errorClass } from '../../common/util';

import { Unauth } from '../containers/Unauth';
import type { APIErrors } from '../../api/types';

type Props = {
  requestSignIn: (u: string, p: string) => Promise<?APIErrors>;
  location: Object;
};

export class Signin extends Component {
  props: Props;
  state: {
    errors: ?APIErrors,
    redirectToReferrer: boolean
  }
  onSubmit: Function;

  constructor(props: Props) {
    super(props)
    this.state = {
      errors: null,
      redirectToReferrer: false
    }
    this.onSubmit = this.onSubmit.bind(this)
  }

  async onSubmit(event) {
    event.preventDefault()
    const response = await FKApiClient.get().signIn(this.refs.username.value, this.refs.password.value);
    log.info(response);
    if (response.type === 'err') {
      if (response.errors) {
        this.setState({ errors: response.errors });
      } else {
        // TODO: someday it'll be fixed
        const fakeError = { code: '100', detail: '', id: '', meta: {}, status: 200 };
        this.setState({ errors: fakeError })
      }
    } else {
      this.setState({ redirectToReferrer: true });
    }
  }

  render() {
    if (this.state.redirectToReferrer) {
      let from = this.props.location.state || '/';
      if (from === '/signin') { from = '/' };

      return <Redirect to={from} />;
    }

    const { errors } = this.state;

    return (
      <Unauth>
        <div className="signin">
          <header>
            <h1>Sign in</h1>
          </header>

          <form onSubmit={this.onSubmit}>
            { this.state.errors &&
              <p className="errors">
                Username or password invalid. Check your information and try again.
              </p> }
            <div className="content">
              <div className={`group ${errorClass(errors, 'username')}`}>
                <label htmlFor="username">Username</label>
                <input ref="username" id="username" name="username" type="text" placeholder="" />
              </div>
              <div className={`group ${errorClass(errors, 'password')}`}>
                <label htmlFor="password">Password</label>
                <input ref="password" id="password" name="password" type="password" placeholder="" />
              </div>
            </div>
            <input type="submit" value="Submit"/>
          </form>
        </div>
      </Unauth>
    )
  }
}
