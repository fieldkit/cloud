import React, { Component } from 'react'

export class MemberRow extends Component {
	render() {
		return (
			<tr>
				<td>
					{this.props.name}<br />
					{this.props.username}
				</td>
				<td>
					{this.props.role}
				</td>
				<td>
					<button className="bt-icon edit"></button>
				</td>				
				<td>
					<button className="bt-icon remove"></button>
				</td>
			</tr>
		)
	}
}

export class MembersTable extends Component {
	constructor() {
		super();
		this.state = {
			members: [
				{name: 'adjany', username: 'some_username', role: 'explorer', avatar_url: 'img/test.png'},
				{name: 'steve', username: 'some_username', role: 'explorer', avatar_url: 'img/test.png'},
				{name: 'jer', username: 'some_username', role: 'explorer', avatar_url: 'img/test.png'},
				{name: 'chris', username: 'some_username', role: 'explorer', avatar_url: 'img/test.png'},
				{name: 'john', username: 'some_username', role: 'explorer', avatar_url: 'img/test.png'}
			]
		}
	}	
	render() {
		const rows = [];
		this.state.members.forEach(function(member) {
			rows.push(<MemberRow name={member.name} username={member.username} role={member.role}/>);
		});
		return (
			<div>
				<table>
					<thead>
						<tr>
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