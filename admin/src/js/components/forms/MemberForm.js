// @flow weak

import React, { Component } from 'react'

import { FormContainer } from '../containers/FormContainer';
import { errorsFor, slugify } from '../../common/util';

import type { APIErrors, APINewMember } from '../../api/types';

type Props = {
  teamId?: number,

  cancelText?: string;
  saveText?: ?string;
  onCancel?: () => void;
  onSave: (e: APINewMember) => Promise<?APIErrors>; 
}

export class MemberForm extends Component {
  props: Props;
  state: {
    member: APINewMember,
    saveDisabled: boolean,
    errors: ?APIErrors
  };

  constructor(props: Props) {
    super(props)
    this.state = {
      member: {
        user_id: 0,
        role: ''
      },
      saveDisabled: false,
      errors: null
    }

    this.loadData();
  }

  componentWillReceiveProps(nextProps: Props) {
    this.setState({
      member: {
        user_id: 0,
        role: ''
      },
      saveDisabled: false,
      errors: null
    });
  }

  async save() {
    const team_id = this.props.team_id;
    const errors = await this.props.onSave(team_id, {
      user_id: this.state.member.user_id,
      role: this.state.member.role
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
        {user_id: 1234, name: 'adjany', username: 'adjany', avatar_url: 'img/test.png'},
        {user_id: 1235, name: 'steve', username: 'steve', avatar_url: 'img/test.png'},
        {user_id: 1236, name: 'jer', username: 'jer', avatar_url: 'img/test.png'},
        {user_id: 1237, name: 'chris', username: 'chris', avatar_url: 'img/test.png'},
        {user_id: 1238, name: 'john', username: 'john', avatar_url: 'img/test.png'}
      ];

    return (
      <FormContainer
        onSave={this.save.bind(this)}
        onCancel={this.props.onCancel}
        saveText={this.props.saveText}
        cancelText={this.props.cancelText}>

        <div className="form-group">
          <label htmlFor="member">Member</label>
          <select name="member" value={this.state.member.user_id} onChange={this.handleInputChange.bind(this)}>
            { errorsFor(this.state.errors, 'member') }
            { users.map((user, i) => 
              <option value={user.user_id}>{user.username}</option>) }
          </select>
        </div>
      </FormContainer>      
    )
  }
}  