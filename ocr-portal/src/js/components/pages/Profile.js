// @flow weak

import React, { Component } from 'react'
import { Route, Link } from 'react-router-dom'

import { FKApiClient } from '../../api/api';
import { ProfileForm } from '../forms/ProfileForm';
import { FormItem } from '../forms/FormItem';

import type { APIErrors, APIUser, APIBaseUser } from '../../api/types';

type Props = {
user: APIUser;
onUpdate: (user: APIUser) => void;

match: Object;
location: Object;
history: Object;
}

export class Profile extends Component {
    props: Props;
    state: {
    id: number,
    username: string,
    name: string,
    bio: string,
    email: string,
    oldPassword: string,
    newPassword: string,
    newPasswordConfirmation: string,
    passwordMessage: ?string,
    passwordErrors: ?APIErrors,
    errors: null
    }

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

    componentWillReceiveProps(nextProps: Props) {
        this.setState({
            name: nextProps.user.name || '',
            bio: nextProps.user.bio || ''
        });
    }

    handleInputChange(event) {
        const target = event.target;
        const value = target.type === 'checkbox' ? target.checked : target.value;
        const name = target.name;

        this.setState({
            [name]: value
        });
    }

    async onUserSave(name: string, bio: string) {
        const {id, username, email} = this.props.user;

        const userRes = await FKApiClient.get().updateUserById(id, {
            username: username,
            name: name,
            bio: bio,
            email: email
        });
        if (userRes.type === 'ok' && userRes.payload) {
            this.props.onUpdate(userRes.payload);
        } else if (userRes.errors) {
            return userRes.errors;
        }
    }

    onPasswordChange(event) {
        event.preventDefault();

        const {newPassword, newPasswordConfirmation} = this.state;

        if (!newPassword || !newPasswordConfirmation || newPassword !== newPasswordConfirmation) {
            return this.setState({
                passwordMessage: 'Invalid fields',
                passwordErrors: null
            });
        }

        const passRes = {}; // await FKApiClient.get().updateUserPassword(this.props.user.id, this.state.oldPassword, this.state.newPassword);
        if (passRes.type !== 'ok') {
            this.setState({
                passwordMessage: null,
                passwordErrors: passRes.errors
            });
        } else {
            this.setState({
                passwordMessage: 'Success!',
                passwordErrors: null
            })
        }
    }

    render() {
        const {name, bio, errors, username, email, oldPassword, newPassword, newPasswordConfirmation} = this.state;

        const onChange = this.handleInputChange.bind(this);

        return (
            <div className="profile-page">
                <h1>Profile</h1>
                <div className="profile">
                    <div className="profile-form two-columns">
                        <h3>Main</h3>
                        <ProfileForm name={ name } bio={ bio } onSave={ this.onUserSave.bind(this) }>
                        </ProfileForm>
                    </div>
                    <div className="two-columns">
                        <div className="account-settings">
                            <h3>Account Settings</h3>
                            <div className="form-group">
                                <label htmlFor="username">Username</label>
                                <span className="disabled-form">{ username }</span>
                            </div>
                            <div className="form-group">
                                <label htmlFor="email">Email</label>
                                <span className="disabled-form">{ email }</span>
                            </div>
                        </div>
                        <div className="password-form">
                            <h3>Change Password</h3>
                            <FormItem type="password" labelText="Old Password" name="oldPassword" value={ oldPassword } errors={ errors } onChange={ onChange } className="lg"
                            />
                            <FormItem type="password" labelText="New Password" name="newPassword" value={ newPassword } errors={ errors } onChange={ onChange } className="lg"
                            />
                            <FormItem type="password" labelText="Confirm New Password" name="newPasswordConfirmation" value={ newPasswordConfirmation } errors={ errors } onChange={ onChange }
                                className="lg" />
                            <input type="submit" onClick={ this.onPasswordChange.bind(this) } value="Change Password" />
                        </div>
                    </div>
                </div>
            </div>
        )
    }
}
