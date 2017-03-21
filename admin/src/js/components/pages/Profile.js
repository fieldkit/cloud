// @flow weak

import React, { Component } from 'react'
import { Route, Link } from 'react-router-dom'

import { FKApiClient } from '../../api/api';
import { errorsFor } from '../../common/util';

import type { APIErrors, APIUserProfile } from '../../api/types';

type Props = {
  user: APIUserProfile;
  onUserUpdate: () => void;

  match: Object;
  location: Object;
  history: Object;
}

/*flow-include
type State = {
  oldPassword: string;
  newPassword: string;
  newPasswordConfirmation: string;
  passwordMessage: ?string;
  passwordErrors: ?APIErrors;

  ...$Exact<APIUserProfile>;
  errors: ?APIErrors;
};
*/

export class Profile extends Component {
  props: Props;
  state: State;

  constructor(props: Props) {
    super(props);

    this.state = {
      oldPassword: '',
      newPassword: '',
      newPasswordConfirmation: '',
      passwordMessage: null,
      passwordErrors: null,
      ...props.user,
      errors: null
    }
  }

  async onUserSave() {
    this.props.onUserUpdate();
  }

  handleInputChange(event) {
    const target = event.target;
    const value = target.type === 'checkbox' ? target.checked : target.value;
    const name = target.name;

    this.setState({ [name]: value });
  }

  onPasswordChange(event) {
    event.preventDefault();

    const { newPassword, newPasswordConfirmation } = this.state;

    if (!newPassword || !newPasswordConfirmation || newPassword !== newPasswordConfirmation) {
      return this.setState({ passwordMessage: 'Invalid fields', passwordErrors: null });
    }

    const passRes = {}; // await FKApiClient.get().updateUserPassword(this.props.user.id, this.state.oldPassword, this.state.newPassword);
    if (passRes.type !== 'ok') {
      this.setState({ passwordMessage: null, passwordErrors: passRes.errors });
    } else {
      this.setState({ passwordMessage: 'Success!', passwordErrors: null })
    }
  }

  render() {
    return (
      <div className="profile-page">
        <h1>Profile</h1>
        <div className="profile-form">
          <div className="form-group">
            <label htmlFor="name">Name</label>
            <input ref="name" id="name" name="name" type="text" value={this.state.name} />
            { errorsFor(this.state.errors, 'name') }
          </div>
          <div className="form-group">
            <label htmlFor="bio">Bio</label>
            <input ref="bio" id="bio" name="bio" type="text" value={this.state.bio} />
            { errorsFor(this.state.errors, 'bio') }
          </div>
        </div>

        <div className="account-settings">
          <div className="form-group">
            <label htmlFor="username">Username</label>
            <span name="username">{this.state.username}</span>
          </div>
          <div className="form-group">
            <label htmlFor="email">Email</label>
            <span name="username">{this.state.email}</span>
          </div>
        </div>

        <div className="password-form">
          <h3>Change Password</h3>
          <div className="form-group">
            <label htmlFor="oldPassword">Old Pasword</label>
            <input type="text" name="oldPassword" className="lg" value={this.state.oldPassword} onChange={this.handleInputChange.bind(this)} />
            { errorsFor(this.state.passwordErrors, 'oldPassword') }
          </div>
          <div className="form-group">
            <label htmlFor="newPassword">Old Pasword</label>
            <input type="text" name="newPassword" className="lg" value={this.state.newPassword} onChange={this.handleInputChange.bind(this)} />
            { errorsFor(this.state.passwordErrors, 'newPassword') }
          </div>
          <div className="form-group">
            <label htmlFor="newPasswordConfirmation">Old Pasword</label>
            <input type="text" name="newPasswordConfirmation" className="lg" value={this.state.newPasswordConfirmation} onChange={this.handleInputChange.bind(this)} />
            { errorsFor(this.state.passwordErrors, 'newPasswordConfirmation') }
          </div>
          <input type="submit" onClick={this.onPasswordChange.bind(this)} value="Change Password" />
        </div>
      </div>
    )
  }
}
