// @flow weak

import React, { Component } from 'react'
import type { APIMember, APIUser, APIErrors } from '../../api/types';

import { RemoveIcon, EditIcon } from '../icons/Icons'

type Props = {
  teamId: number,
  members: APIMember[],
  users: APIUser[],
  onDelete: (teamId: number, userId: number) => Promise<?APIErrors>
}

export class MembersTable extends Component {
  props: Props;
  state: {
    users: {[id: number]: APIUser},
    errors: ?APIErrors    
  }

  constructor(props: Props) {
    super(props);
    this.state = {
      users: {},
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

  async delete(userId: number) {
    const { teamId } = this.props;

    const errors = await this.props.onDelete(teamId, userId);
    if (errors) {
      this.setState({ errors });
    }    
  }

  render() {
    const { teamId, members } = this.props;
    let { users } = this.state

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
                  <span>{member.role}</span>
                  <div className="bt-icon medium">
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
