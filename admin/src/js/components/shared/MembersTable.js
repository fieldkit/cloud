// @flow weak

import React, { Component } from 'react'
import type { APIMember, APIUser, APIErrors } from '../../api/types';

import { RemoveIcon, EditIcon } from '../icons/Icons'

type Props = {
  teamId: number,
  members: APIMember[],
  users: APIUser[],
  onDelete: (teamId: number, userId: number) => Promise<?APIErrors>,
  onMemberUpdate: (teamId: number, memberId: number, role: string) => void
}

export class MembersTable extends Component {
  props: Props;
  state: {
    users: {[id: number]: APIUser},
    editingMemberId: number,
    errors: ?APIErrors
  }

  constructor(props: Props) {
    super(props);
    this.state = {
      users: {},
      editingMemberId: -1,
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
    const { teamId, onMemberUpdate } = this.props;
    // console.log(teamId, memberId, value);
    onMemberUpdate(teamId, memberId, value);
  }

  handleKeyUp(event) {
    console.log(event.key);
    if (event.key === 'Enter') {
      console.log('save');
    }
  }

  startMemberEdit(memberId: number) {
    console.log(memberId);
    this.setState({editingMemberId: memberId});
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
    let { users, editingMemberId } = this.state

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
                    disabled={member.user_id !== editingMemberId}
                    onChange={this.handleInputChange.bind(this, member.user_id)}/>
                  <div id={member.user_id} className={'bt-icon medium ' + (member.user_id === editingMemberId ? 'disabled' : '') } onClick={this.startMemberEdit.bind(this, member.user_id)}>
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
