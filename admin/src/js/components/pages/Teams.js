// @flow weak

import React, { Component } from 'react'
import { Route, Link } from 'react-router-dom'
import ReactModal from 'react-modal';

import { MembersTable } from '../shared/MembersTable';
import { TeamForm } from '../forms/TeamForm';
import { MemberForm } from '../forms/MemberForm';
import { FKApiClient } from '../../api/api';

import type { APIProject, APIExpedition, APITeam, APINewTeam, APINewMember, APIMember } from '../../api/types';

type Props = {
  project: APIProject;
  expedition: APIExpedition;

  match: Object;
  location: Object;
  history: Object;
}

export class Teams extends Component {

  props: Props;
  state: {
    teams: APITeam[],
    members: { [teamId: number]: APIMember[] }
  }

  constructor (props: Props) {
    super(props);
    this.state = {
      teams: [],
      members: {}
    }

    this.loadTeams();
  }

  async loadTeams() {
    const teamsRes = await FKApiClient.get().getTeamsBySlugs(this.props.project.slug, this.props.expedition.slug);
    if (teamsRes.type === 'ok' && teamsRes.payload) {
      const { teams } = teamsRes.payload;
      this.setState({ teams: teams });
      for (const team of teams) {
        await this.loadMembers(team.id);
      }
    }
  }

  async loadMembers(teamId: number) {
    const membersRes = await FKApiClient.get().getMembers(teamId);  
    if (membersRes.type === 'ok' && membersRes.payload) {
      const { members } = this.state;
      members[teamId] = membersRes.payload.members;
      this.setState({members: members});
    }
  }

  async onTeamCreate(e: APINewTeam) {
    
    const { expedition, match } = this.props;

    const teamRes = await FKApiClient.get().createTeam(expedition.id, e);
    if (teamRes.type === 'ok') {
      await this.loadTeams();
      this.props.history.push(`${match.url}`);
    } else {
      return teamRes.errors;
    }
  }

  async onMemberAdd(teamId: number, e: APINewMember) {
    const { match } = this.props;
    const memberRes = await FKApiClient.get().addMember(teamId, e);
    if (memberRes.type === 'ok') {
      await this.loadTeams();
      this.props.history.push(`${match.url}`);
    } else {
      return memberRes.errors;
    }
  }

  async onMemberDelete(teamId: number, userId: number) {
    const { match } = this.props;
    const memberRes = await FKApiClient.get().deleteMember(teamId, userId);
    if (memberRes.type === 'ok') {
      await this.loadTeams();
      this.props.history.push(`${match.url}`);
    } else {
      return memberRes.errors;
    }    
  }

  render() {
    const { match } = this.props;
    const { members, teams } = this.state;
    
    return (
      <div className="teams">
        <Route path={`${match.url}/new-team`} render={() =>
          <ReactModal isOpen={true} contentLabel="Create New Team">
            <h1>Create a new team</h1>
            <TeamForm
              onCancel={() => this.props.history.push(`${match.url}`)}
              onSave={this.onTeamCreate.bind(this)} />
          </ReactModal> } />

        <Route path={`${match.url}/:teamId/add-member`} render={props =>
          <ReactModal isOpen={true} contentLabel="Add Members">
            <h1>Add Members</h1>
            <MemberForm
              teamId={props.match.params.teamId}
              onCancel={() => this.props.history.push(`${match.url}`)}
              onSave={this.onMemberAdd.bind(this)} />
          </ReactModal> } />

        <h1>Teams</h1>

        <div id="teams">
        { teams.map((team, i) =>
          <table key={i} className="teams-table">
            <tbody>
              <tr>
                <td name={team.name}>
                  {team.name}<br/>
                  {team.description}                  
                  <button className="secondary">Edit</button>

                  { <MembersTable
                      teamId={team.id}
                      members={members[team.id]}
                      onDelete={this.onMemberDelete.bind(this)}/> }
                  { !members[team.id] &&
                    <span className="empty">No members</span> }
                  <Link className="button secondary" to={`${match.url}/${team.id}/add-member`}>Add Member</Link>                

                </td>
              </tr>
            </tbody>
          </table> ) }
        { teams.length === 0 &&
          <span className="empty">No teams</span> }
        </div>
        <Link className="button" to={`${match.url}/new-team`}>Create New Team</Link>
      </div>
    )
  }
}
