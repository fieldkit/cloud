// @flow weak

import React, { Component } from 'react'

import { FormContainer } from '../containers/FormContainer';
import { errorsFor, slugify } from '../../common/util';

import type { APIErrors, APIProject, APIUser, APINewAdministrator } from '../../api/types';
import { FKApiClient } from '../../api/api';

type Props = {
  project: APIProject,
//   member?: APIMember,

//   cancelText?: string;
//   saveText?: ?string;
//   onCancel?: () => void;
//   onSave: (teamId: number, e: APINewMember) => Promise<?APIErrors>; 
}

export class AdministratorForm extends Component {
  props: Props;
  state: {
    users: APIUser[],    
    userId: number,
    saveDisabled: boolean,
    errors: ?APIErrors    
  }

  constructor(props: Props) {
    super(props);
    this.state = {
      users: [],
      userId: -1,
      saveDisabled: true,
      errors: null      
    }
    this.loadUsers();
  }

  async loadUsers() {
    /* TO DO:
      remove already added members from this list.
      Could probably pass current members in the props from the
      Teams table and use componentWillReceiveProps to set the state */
    const usersRes = await FKApiClient.get().getUsers();
    if (usersRes.type === 'ok' && usersRes.payload) {
      this.setState({users: usersRes.payload.users} || []);
    }
  }

  handleInputChange(event) {
    const target = event.target;
    const value = target.type === 'checkbox' ? target.checked : target.value;
    const name = target.name;

    this.setState({ [name]: value });
  }

  async save() {
  //   const errors = await this.props.onSave(this.props.teamId, {
  //     user_id: parseInt(this.state.userId),
  //     role: this.state.role
  //   });
  //   if (errors) {
  //     this.setState({ errors });
  //   }
  }

  render () {
    
    const { users, userId } = this.state;

    return (
      <FormContainer
        onSave={this.save.bind(this)}>

        <div className="form-group">
          <label htmlFor="userId">Member</label>
          <select name="userId"  className='lg' value={this.state.userId} onChange={this.handleInputChange.bind(this)}>
            <option value={null}>Select a user</option>
            { users.map((user, i) => 
              <option key={i} value={user.id}>{user.username}</option>) }
          </select>
          { errorsFor(this.state.errors, 'userId') }          
        </div>

      </FormContainer>
    )
  }
}