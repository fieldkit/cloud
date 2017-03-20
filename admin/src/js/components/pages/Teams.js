import React, { Component } from 'react'
import { Route, Link } from 'react-router-dom'
import ReactModal from 'react-modal';

import { MembersTable } from '../shared/MembersTable';
import { TeamForm } from '../forms/TeamForm';
import { FKApiClient } from '../../api/api';

import type { APIProject, APIExpedition, APITeam, APINewTeam } from '../../api/types';

type Props = {
  project: APIProject;
  expedition: APIExpedition;

  match: Object;
  location: Object;
  history: Object;
}

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

// export class TeamsTable extends Component {
//   constructor() {
//     super();
//     this.state = {
//       teams: [
//         {
//           name: 'river',
//           members: [
//             {name: 'adjany', username: 'some_username', role: 'explorer'},
//             {name: 'steve', username: 'some_username', role: 'explorer'},
//             {name: 'jer', username: 'some_username', role: 'explorer'},
//             {name: 'chris', username: 'some_username', role: 'explorer'},
//             {name: 'john', username: 'some_username', role: 'explorer'}
//           ]
//         },
//         {
//           name: 'ground',
//           members: [
//             {name: 'kate', username: 'some_username', role: 'creative researcher'},
//             {name: 'eric', username: 'some_username', role: 'creative researcher'},
//             {name: 'chris', username: 'some_username', role: 'creative researcher'},
//             {name: 'noa', username: 'some_username', role: 'creative researcher'}
//           ]
//         }
//       ]
//     }
//   }
//   render() {

//     return (

//     )
//   }
// }

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
    const { project, expedition } = this.props;
    const projectSlug = project.slug;
    const expeditionSlug = expedition.slug;

    const teamRes = await FKApiClient.get().createTeam(expedition.id, e);
    if (teamRes.type === 'ok') {
      await this.loadData();
      this.props.history.push(`/projects/${projectSlug}/expeditions/${expeditionSlug}/teams`);
    } else {
      return teamRes.errors;
    }
  }  

  render() {
    const { project } = this.props;
    const projectSlug = project.slug;
    const { expedition } = this.props;
    const expeditionSlug = expedition.slug;
    
    return (
      <div className="teams">
        <Route path="/projects/:projectSlug/expeditions/:expeditionSlug/teams/new-team" render={() =>
          <ReactModal isOpen={true} contentLabel="Create New Team">
            <h1>Create a new team</h1>
            <TeamForm
              projectSlug={projectSlug}
              expeditionSlug={expeditionSlug}
              onCancel={() => this.props.history.push(`/projects/${projectSlug}/expeditions/${expeditionSlug}/teams`)}
              onSave={this.onTeamCreate.bind(this)} />
          </ReactModal> } />

        <h1>Teams</h1>

        <div id="teams">
        { this.state.teams.map((team, i) =>
          <table className="teams-table">
            <tbody>
              <TeamRow key={i} name={team.name} members={team.members} />
            </tbody>
          </table> ) }
        { this.state.teams.length === 0 &&
          <span className="empty">No teams</span> }
        </div>

        {/* <button>Add Team</button> */}
        <Link to={`/projects/${projectSlug}/expeditions/${expeditionSlug}/teams/new-team`}>Show new team modal</Link>
      </div>
    )
  }
}
