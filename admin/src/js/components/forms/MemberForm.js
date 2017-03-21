// @flow weak

import React, { Component } from 'react'

import { FormContainer } from '../containers/FormContainer';
import { errorsFor, slugify } from '../../common/util';

import type { APIErrors, APINewMember } from '../../api/types';

type Props = {
  teamId: number,

  cancelText?: string;
  saveText?: ?string;
  onCancel?: () => void;
  onSave: (teamId: number, e: APINewMember) => Promise<?APIErrors>; 
}

export class MemberForm extends Component {
  props: Props;
  state: {
    userId: number,
    role: string,
    saveDisabled: boolean,
    errors: ?APIErrors
  };

  constructor(props: Props) {
    super(props)
    this.state = {
      userId: 0,
      role: '',
      saveDisabled: false,
      errors: null
    }

    this.loadData();
  }

  componentWillReceiveProps(nextProps: Props) {
    this.setState({
      userId: 0,
      role: '',
      saveDisabled: false,
      errors: null
    });
  }

  async save() {
    const errors = await this.props.onSave(this.props.teamId, {
      user_id: parseInt(this.state.userId),
      role: this.state.role
    });
    if (errors) {
      this.setState({ errors });
    }
  }

  handleInputChange(event) {
    const target = event.target;
    const value = target.type === 'checkbox' ? target.checked : target.value;
    const name = target.name;

    this.setState({ [name]: value });
  }  

  async loadData() {
    // TO DO: load real list of users from server
  }

  render () {
    const users = [
        {userId: 1234, name: 'adjany', username: 'adjany', avatar_url: 'img/test.png'},
        {userId: 1235, name: 'steve', username: 'steve', avatar_url: 'img/test.png'},
        {userId: 1236, name: 'jer', username: 'jer', avatar_url: 'img/test.png'},
        {userId: 1237, name: 'chris', username: 'chris', avatar_url: 'img/test.png'},
        {userId: 1238, name: 'john', username: 'john', avatar_url: 'img/test.png'}
      ];

    return (
      <FormContainer
        onSave={this.save.bind(this)}
        onCancel={this.props.onCancel}
        saveText={this.props.saveText}
        cancelText={this.props.cancelText}>

        <div className="form-group">
          <label htmlFor="userId">Member</label>
          <select name="userId"  className='lg' value={this.state.userId} onChange={this.handleInputChange.bind(this)}>
            <option value={null}>Select a user</option>
            { users.map((user, i) => 
              <option key={i} value={user.userId}>{user.username}</option>) }
          </select>
          { errorsFor(this.state.errors, 'userId') }          
        </div>

        <div className="form-group">
          <label htmlFor="role">Role</label>
          <input type="text" name="role" className="lg" value={this.state.role} onChange={this.handleInputChange.bind(this)} />
          { errorsFor(this.state.errors, 'role') }
        </div>
      </FormContainer>
    )
  }
}  