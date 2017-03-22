// @flow weak

import React, { Component } from 'react'
import type { APIMember } from '../../api/types';

type Props = {
  members: APIMember[];
}

export class MemberRow extends Component {
  render() {
    return (
      <tr>
        <td>
          <div className="user-avatar">
          </div>
        </td>
        <td>
          {/*this.props.username*/}
          {this.props.userId}
        </td>
        <td>
          {this.props.role} <button className="bt-icon edit"></button>
          <button className="bt-icon remove"></button>
        </td>
      </tr>
    )
  }
}

export class MembersTable extends Component {
  props: Props;
  constructor(props: Props) {
    super(props);
  }

  render() {
    const { members } = this.props;

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
                  <MemberRow key={i} userId={member.user_id} role={member.role} /> )}
              </tbody>
            </table>
          </div>
        )
    }
    return (<span className="empty">No members</span>)
  }
}
