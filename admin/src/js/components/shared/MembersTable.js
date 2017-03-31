// @flow weak

import React, { Component } from 'react'
import type { APIMember, APIUser, APIErrors, APIBaseMember } from '../../api/types';

import { RemoveIcon, EditIcon } from '../icons/Icons'

type Props = {
  teamId: number,
  members: APIMember[],
  users: APIUser[],
  onDelete: (teamId: number, userId: number) => Promise<?APIErrors>,
  onStartUpdate: (teamId: number, memberId: number, values: APIBaseMember) => void,
  onUpdate: (teamId: number, memberId: number) => void
}

export class MembersTable extends Component {
  props: Props;
  state: {
    users: {[id: number]: APIUser},
    user_id: number,
    role: string,
    errors: ?APIErrors
  }

  constructor(props: Props) {
    super(props);
    this.state = {
      users: {},
      user_id: -1,
      role: '',
      errors: null
    }
  }

  componentWillReceiveProps(nextProps: Props) {
    let { users } = nextProps;
    if(users){
      const mappedUsers = {};
      for(var u of users){
        mappedUsers[u.id] = u;
        this.setState({users: mappedUsers});
      }
    }
  }

  handleInputChange(memberId: number, event) {
    const target = event.target;
    const value = target.value;
    const { teamId, onStartUpdate } = this.props;

    onStartUpdate(teamId, memberId, {role: value});
  }

  handleKeyUp(memberId: number, event) {
    const { teamId, onUpdate } = this.props;
    const { user_id } = this.state;
    if (event.key === 'Enter') {
      onUpdate(teamId, memberId);
      this.setState({user_id: -1});
    }
  }

  startMemberEdit(memberId: number) {
    this.setState({user_id: memberId});
  }

  handleBlur(){
    this.setState({user_id: -1});
  }

  async delete(userId: number) {
    const { teamId } = this.props;

    const errors = await this.props.onDelete(teamId, userId);
    if (errors) {
      this.setState({ errors });
    }    
  }

  render() {
    const { teamId, members } = this.props;
    let { users, user_id } = this.state

    return (
      <table className="members-table">
        <thead>
          <tr>
            <th></th>
            <th>Members ({ members.length })</th>
            <th>Role</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          { members.map((member, i) =>
              <tr key={i}>
                <td>
                  <div className="user-avatar medium">
                  </div>
                </td>
                <td>
                  {users[member.user_id] &&
                    <div>
                      <p>{users[member.user_id].name}</p>
                      <p className="type-small">{users[member.user_id].username}</p>
                      </div> }
                </td>
                <td>
                  <input type="text"
                    value={member.role}
                    name="role"
                    disabled={member.user_id !== user_id}
                    onChange={this.handleInputChange.bind(this, member.user_id)}
                    onKeyUp={this.handleKeyUp.bind(this, member.user_id)}
                    onBlur={this.handleBlur.bind(this)}/>
                  <div id={member.user_id} className={'bt-icon medium ' + (member.user_id === user_id ? 'disabled' : '') } onClick={this.startMemberEdit.bind(this, member.user_id)}>
                    <EditIcon />
                  </div>                        
                </td>
                <td>
                  <div className="bt-icon medium" onClick={this.delete.bind(this, member.user_id)}>
                    <RemoveIcon />
                  </div>
                </td>
              </tr> )}                
        </tbody>
      </table>
    )
  }
}
