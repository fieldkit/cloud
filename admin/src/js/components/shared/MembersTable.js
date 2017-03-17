import React, { Component } from 'react'

export class MemberRow extends Component {
  render() {
    return (
      <tr>
        <td>
          <div className="user-avatar">
          </div>
        </td>
        <td>
          {this.props.name}<br />
          {this.props.username}
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
  constructor() {
    super();
  }
  render() {
    const rows = this.props.members.map((member, i) => <MemberRow key={i} name={member.name} username={member.username} role={member.role}/>);
    return (
      <div>
        <table className="members-table">
          <thead>
            <tr>
              <th></th>
              <th>Members ({rows.length})</th>
              <th>Role</th>
            </tr>
          </thead>
          <tbody>{rows}</tbody>
        </table>
      </div>
    )
  }
}
