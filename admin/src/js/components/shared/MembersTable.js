// @flow weak

import React, { Component } from 'react'
import type { APIMember, APIErrors } from '../../api/types';

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
          <div>
            <table className="members-table">
              <thead>
                <tr>
                  <th></th>
                  <th>Members ({ members.length })</th>
                  <th>Role</th>
                </tr>
              </thead>
              <tbody>
                { members.map((member, i) =>
                    <tr>
                      <td>
                        <div className="user-avatar">
                        </div>
                      </td>
                      <td>
                        {/*member.username*/}
                        {member.user_id}
                      </td>
                      <td>
                        {member.role} <button className="bt-icon edit"></button>
                        <button className="bt-icon remove" onClick={this.delete.bind(this, member.user_id)}></button>
                      </td>
                    </tr> )}                
              </tbody>
            </table>
          </div>
        )
    }
    return (<span className="empty">No members</span>)
  }
}
