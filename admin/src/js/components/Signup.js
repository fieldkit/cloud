// @flow weak

import React, { Component } from 'react'

import { errorsFor } from '../common/util';
import type { ErrorMap } from '../common/util';

type Props = {
  requestSignUp: (e: string, u: string, p: string, i: string) => Promise<?ErrorMap>;
};

export class Signup extends Component {
  props: Props;
  state: {
    errors: ?ErrorMap
  }
  onSubmit: Function;

  constructor(props: Props) {
    super(props)
    this.state = {
      errors: null
    }
    this.onSubmit = this.onSubmit.bind(this)
  }

  async onSubmit(event) {
    event.preventDefault()
    const errors = await this.props.requestSignUp(
      this.refs.email.value,
      this.refs.username.value,
      this.refs.password.value,
      this.refs.invite.value
    );
    this.setState({ errors });
  }

  render() {
    return (
      <div className="signup">
        <header>
          <h1>Join us</h1>
        </header>

        <form onSubmit={this.onSubmit}>
          <div className="content">
            <div className="group">
              <label htmlFor="email">Email</label>
              <input ref="email" id="email" name="email" type="email" placeholder="explorer@email.com" />
              { errorsFor(this.state.errors, 'email') }
            </div>

            <div className="group">
              <label htmlFor="username">Username</label>
              <input ref="username" id="username" name="username" type="text" placeholder="explorer123" />
              { errorsFor(this.state.errors, 'username') }
            </div>

            <div className="group">
              <label htmlFor="password">Password</label>
              <input ref="password" id="password" name="password" type="password" placeholder="correct horse battery staple" />
              { errorsFor(this.state.errors, 'password') }
            </div>

            <div className="group">
              <label htmlFor="invite">Invitation token</label>
              <input ref="invite" id="invite" name="invite" type="text"  placeholder="E3NANDTJ3YXM5LNMGNZTF2373LAFTOCC"/>
              { errorsFor(this.state.errors, 'invite') }
            </div>
          </div>

          <footer>
            { this.state.errors &&
              <p className="errors">
                We found one or multiple errors. Please check your information above or try again later.
              </p> }
            <input type="submit" value="Submit"/>
          </footer>
        </form>
      </div>
    )
  }
}
