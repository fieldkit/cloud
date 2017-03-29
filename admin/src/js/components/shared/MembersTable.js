// @flow weak

import React, { Component } from 'react'
import type { APIMember, APIErrors } from '../../api/types';

import { RemoveIcon, EditIcon } from '../icons/Icons'

type Props = {
  teamId: number,
  members: APIMember[],

  onDelete: (teamId: number, userId: number) => Promise<?APIErrors>
}

export class MembersTable extends Component {
  props: Props;
  state: {
    errors: ?APIErrors      
  }

  constructor(props: Props) {
    super(props);
    state: {
      errors: null
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
    if(members){
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
                <tr>
                  <td>
                    <div className="user-avatar medium">
                    </div>
                  </td>
                  <td>
                    {/*member.username*/}
                    {member.user_id}
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
    return (<span className="empty">No members</span>)
  }
}
