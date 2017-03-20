// @flow weak

import React, { Component } from 'react'
import { Route, Link } from 'react-router-dom'
import ReactModal from 'react-modal';

import { MembersTable } from '../shared/MembersTable';
import { TeamForm } from '../forms/TeamForm';
import { MemberForm } from '../forms/MemberForm';
import { FKApiClient } from '../../api/api';

import type { APIProject, APIExpedition, APITeam, APINewTeam, APINewMember } from '../../api/types';

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
    teams: APITeam[]
  }

  constructor (props: Props) {
    super(props);
    this.state = {
      teams: []
    }

    this.loadData();
  }

  async loadData() {
    const teamsRes = await FKApiClient.get().getTeamsBySlugs(this.props.project.slug, this.props.expedition.slug);
    if (teamsRes.type === 'ok' && teamsRes.payload) {
      this.setState({ teams: teamsRes.payload.teams || [] })
    }
  }
  
  createTeam(expeditionId: number, values: APINewTeam): Promise<FKAPIResponse<APITeam>> {
    return this.postWithErrors(`/expeditions/${expeditionId}/team`, values)
  }

  async onTeamCreate(e: APINewTeam) {
    
    const { expedition, match } = this.props;

    const teamRes = await FKApiClient.get().createTeam(expedition.id, e);
    if (teamRes.type === 'ok') {
      await this.loadData();
      this.props.history.push(`${match.url}`);
    } else {
      return teamRes.errors;
    }
  }

  async onMemberAdding(team_id: number, e: APINewMember) {

    const { match } = this.props;

    const memberRes = await FKApiClient.get().addMember(team_id, e);
    if (memberRes.type === 'ok') {
      await this.loadData();
      this.props.history.push(`${match.url}`);
    } else {
      return memberRes.errors;
    }
  }

  render() {
    const { match } = this.props;
    
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
              team_id={props.match.params.teamId}
              onCancel={() => this.props.history.push(`${match.url}`)}
              onSave={this.onMemberAdding.bind(this)} />
          </ReactModal> } />

        <h1>Teams</h1>

        <div id="teams">
        { this.state.teams.map((team, i) =>
          <table className="teams-table">
            <tbody>
              <tr>
                <td key={i} name={team.name}>
                  {team.name}, {team.id}
                  <button>Edit</button>
                  <Link className="button" to={`${match.url}/${team.id}/add-member`}>Add Member</Link>
                </td>
              </tr>
            </tbody>
          </table> ) }
        { this.state.teams.length === 0 &&
          <span className="empty">No teams</span> }
        </div>
        <Link className="button" to={`${match.url}/new-team`}>Create New Team</Link>
      </div>
    )
  }
}
