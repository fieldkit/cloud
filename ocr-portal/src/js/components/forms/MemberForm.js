// @flow weak

import React, { Component } from 'react'

import { FormContainer } from '../containers/FormContainer';
import { errorsFor, slugify } from '../../common/util';

import type { APIErrors, APIUser, APINewMember, APIMember } from '../../api/types';
import { FKApiClient } from '../../api/api';

import { FormSelectItem } from './FormSelectItem'

type Props = {
teamId: number,
member?: APIMember,
members: APIMember[],

cancelText?: string;
saveText?: ?string;
onCancel?: () => void;
onSave: (teamId: number, e: APINewMember) => Promise<?APIErrors>;
}

export class MemberForm extends Component {
    props: Props;
    state: {
    users: APIUser[],
    userId: number,
    role: string,
    saveDisabled: boolean,
    errors: ?APIErrors
    };

    constructor(props: Props) {
        super(props)
        this.state = {
            users: [],
            userId: -1,
            role: '',
            saveDisabled: false,
            errors: null
        }

        this.loadUsers();
    }

    componentWillReceiveProps(nextProps: Props) {
        this.setState({
            users: [],
            userId: -1,
            role: '',
            saveDisabled: false,
            errors: null
        });
    }

    async loadUsers() {
        const usersRes = await FKApiClient.get().getUsers();
        if (usersRes.type === 'ok' && usersRes.payload) {
            const membersIds = this.props.members.map(member => member.user_id);
            const availableUsers = usersRes.payload.users.filter(user => membersIds.indexOf(user.id) < 0);
            this.setState({
                    users: availableUsers
                } || []);
        }
    }

    async save() {
        const errors = await this.props.onSave(this.props.teamId, {
            user_id: parseInt(this.state.userId),
            role: this.state.role
        });
        if (errors) {
            this.setState({
                errors
            });
        }
    }

    handleInputChange(event) {
        const target = event.target;
        const value = target.type === 'checkbox' ? target.checked : target.value;
        const name = target.name;

        this.setState({
            [name]: value
        });
    }

    render() {
        const {users, errors} = this.state;
        const options = users.map((user, i) => {
            return {
                value: user.id,
                text: user.username
            }
        });

        return (
            <FormContainer onSave={ this.save.bind(this) } onCancel={ this.props.onCancel } saveText={ this.props.saveText } cancelText={ this.props.cancelText }>
                <FormSelectItem labelText={ 'Member' } name={ 'userId' } className={ 'lg' } value={ this.state.userId } firstOptionText={ 'Select a user' }
                    options={ options } errors={ errors } onChange={ this.handleInputChange.bind(this) } />
                <div className="form-group">
                    <label htmlFor="role">Role</label>
                    <input type="text" name="role" className="lg" value={ this.state.role } onChange={ this.handleInputChange.bind(this) } />
                    { errorsFor(this.state.errors, 'role') }
                </div>
            </FormContainer>
        )
    }
}