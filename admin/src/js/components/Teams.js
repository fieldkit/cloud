import React, { Component } from 'react'
import { MembersTable } from './shared/MembersTable';

export class TeamRow extends Component {	
	render() {
		return (
			<tr>
				<td>
					{this.props.name}
					<MembersTable members={this.props.members}>
					</MembersTable>
					<button>Add Members</button>
				</td>
			</tr>
		)
	}
}

export class TeamsTable extends Component {
	constructor() {
		super();
		this.state = {
			teams: [
				{
					name: 'river',
					members: [
						{name: 'adjany', username: 'some_username', role: 'explorer'},
						{name: 'steve', username: 'some_username', role: 'explorer'},
						{name: 'jer', username: 'some_username', role: 'explorer'},
						{name: 'chris', username: 'some_username', role: 'explorer'},
						{name: 'john', username: 'some_username', role: 'explorer'}
					]
				},
				{
					name: 'ground',
					members: [
						{name: 'kate', username: 'some_username', role: 'creative researcher'},
						{name: 'eric', username: 'some_username', role: 'creative researcher'},
						{name: 'chris', username: 'some_username', role: 'creative researcher'},
						{name: 'noa', username: 'some_username', role: 'creative researcher'}
					]					
				}
			]
		}
	}
	render() {
		const rows = [];
		this.state.teams.forEach(function(team) {
			rows.push(<TeamRow name={team.name} members={team.members} />);
		});
		return (
			<table className="teams-table">
				<tbody>{rows}</tbody>
			</table>
		)
	}
}

export class Teams extends Component {
	render() {
		return (
			<div className="teams">
				<h1>Teams</h1>
				<TeamsTable>
				</TeamsTable>
				<button>Add Team</button>
			</div>
		)
	}
}