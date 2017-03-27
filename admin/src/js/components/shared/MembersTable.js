// @flow weak

import React, { Component } from 'react'
import type { APIMember, APIErrors } from '../../api/types';
import removeImg from '../../../img/icons/icon-remove-small.png'
import editImg from '../../../img/icons/icon-edit-small.png'

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
                  <th></th>
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
                        <span>{member.role}</span>
                        <div className="bt-icon">
                          <img src={editImg} alt="external link" />
                        </div>                        
                      </td>
                      <td>
                        <div className="bt-icon" onClick={this.delete.bind(this, member.user_id)}>
                          <img src={removeImg} alt="external link" />
                        </div>
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
