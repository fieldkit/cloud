// @flow weak

import React, { Component } from 'react'
import { Route, Link } from 'react-router-dom'

import { FKApiClient } from '../../api/api';
import { FormItem } from '../forms/FormItem';

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
    const {
      name,
      bio,
      errors,
      username,
      email,
      oldPassword,
      newPassword,
      newPasswordConfirmation
    } = this.state;

    const onChange = this.handleInputChange.bind(this);

    return (
      <div className="profile-page">
        <h1>Profile</h1>
        
        <div className="row">
          <div className="profile-form two-columns">
            <FormItem labelText="Name" name="name" value={name} errors={errors} onChange={onChange} className="lg" />
            <FormItem labelText="Bio" name="bio" value={bio} errors={errors} onChange={onChange} className="lg" />
          </div>
        </div>

        <div className="row">
          <div className="account-settings two-columns">
            <div className="form-group">
              <label htmlFor="username">Username</label>
              <span className="disabled-form">{username}</span>
            </div>
            <div className="form-group">
              <label htmlFor="email">Email</label>
              <span className="disabled-form">{email}</span>
            </div>
          </div>
        </div>

        <div className="row">
          <div className="password-form two-columns">
            <h3>Change Password</h3>
            <FormItem labelText="Old Password" name="oldPassword" value={oldPassword} errors={errors} onChange={onChange} className="lg" />
            <FormItem labelText="New Password" name="newPassword" value={newPassword} errors={errors} onChange={onChange} className="lg" />
            <FormItem labelText="Confirm New Password" name="newPasswordConfirmation" value={newPasswordConfirmation} errors={errors} onChange={onChange} className="lg" />

            <input type="submit" onClick={this.onPasswordChange.bind(this)} value="Change Password" />
          </div>
        </div>

      </div>
    )
  }
}
