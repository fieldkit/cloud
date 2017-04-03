// @flow weak

import React, { Component } from 'react'
import { Route, Link, Redirect } from 'react-router-dom'
import ReactModal from 'react-modal'
import _ from 'lodash'
import log from 'loglevel';

import { ProjectExpeditionForm } from '../forms/ProjectExpeditionForm';
import { InputForm } from '../forms/InputForm';
import { FKApiClient } from '../../api/api';
import { RouteOrLoading } from '../shared/RequiredRoute';

import type { APIProject, APIExpedition, APINewExpedition, APIInputs, APITwitterInputCreateResponse, APINewTwitterInput, APINewFieldkitInput, APITeam, APIMember, APIUser } from '../../api/types';

type Props = {
  project: APIProject;
  expedition: APIExpedition;

  match: Object;
  location: Object;
  history: Object;
}

export class DataSources extends Component {
  props: Props;
  state: {
    inputs: APIInputs,
    teams: APITeam[],
    members: { [teamId: number]: APIMember[] },
    users: APIUser[]
  }

  constructor(props: Props) {
    super(props);
    this.state = {
      inputs: {},
      teams: [],
      members: {},
      users: []
    }

    this.loadInputs();
    this.loadTeams();
  }

  async loadInputs() {
    const inputsRes = await FKApiClient.get().getExpeditionInputs(this.props.expedition.id);
    if (inputsRes.type === 'ok') {
      this.setState({ inputs: inputsRes.payload });
      console.log(inputsRes.payload);
    }
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
      for (const member of members[teamId]) {
        await this.loadMemberName(teamId, member.user_id);
      }      
    }
  }

  async loadMemberName(teamId: number, userId: number){
    const userRes = await FKApiClient.get().getUserById(userId);
    if (userRes.type === 'ok' && userRes.payload) {
      const { users } = this.state;
      if(userRes.payload && !users.find(user => userRes.payload)){
        users.push(userRes.payload);
      }
      this.setState({users: users});
    }
  }  

  async onTwitterCreate(event, t: APINewTwitterInput) {
    const res = await FKApiClient.get().createTwitterInput(this.props.expedition.id, t);
    if (res.type === 'ok') {
      const redirect = res.payload.location;
      window.location = redirect;
    } else {
      log.error('Bad request to create twitter input!');
    }
  }

  async onFieldkitCreate(f: APINewFieldkitInput) {
    const { expedition, match } = this.props;
    const fieldkitRes = await FKApiClient.get().createFieldkitInput(expedition.id, f);
    console.log(fieldkitRes);
    if (fieldkitRes.type === 'ok') {
      console.log(fieldkitRes.payload);
      await this.loadInputs();
      this.props.history.push(`${match.url}`);
    } else {
      return fieldkitRes.errors
    }
  }

  render() {
    const { expedition, match } = this.props;
    const { twitter_accounts, fieldkit_inputs } = this.state.inputs;
    const { users } = this.state;

    return (
      <div className="data-sources-page">


      <Route path={`${match.url}/add-fieldkit-input`} render={props =>
        <ReactModal isOpen={true} contentLabel="Add Fieldkit Input" className="modal" overlayClassName="modal-overlay">
          <h2>Add Fieldkit Input</h2>
          <InputForm
            expeditionId={props.match.params.expeditionId}
            users={users}
            onCancel={() => this.props.history.push(`${match.url}`)}
            onSave={this.onFieldkitCreate.bind(this)} 
            saveText="Add" />
        </ReactModal> } />    

        <h1>Data Sources</h1>

        <div className="input-section">
          <h3>Twitter</h3>
          { twitter_accounts && twitter_accounts.length > 0 &&
            <table>
              <tbody>
                <tr>
                  <th>Username</th>
                  <th>Binding</th>
                  <th></th>
                </tr>
              { twitter_accounts.map((t, i) =>
                <tr key={i} className="input-item">
                  <td>{t.screen_name}</td>
                  <td>None</td>
                  {/* TODO: hook up */}
                  <td><a href="#">Delete</a></td>
                </tr> )}
              </tbody>
            </table> }
          <button className="add-button" onClick={this.onTwitterCreate.bind(this)}>Add Twitter Account</button>

          <h3>Fieldkit Sensors</h3>
          { fieldkit_inputs && fieldkit_inputs.length > 0 &&
            <table>
              <tbody>
                <tr>
                  <th>Name</th>
                  <th>Binding</th>
                  <th></th>
                </tr>
              { fieldkit_inputs.map((t, i) =>
                <tr key={i} className="input-item">
                  <td>{t.name}</td>
                  <td>None</td>
                  {/* TODO: hook up */}
                  <td><a href="#">Delete</a></td>
                </tr> )}
              </tbody>
            </table> }

          <Link className="button" to={`${match.url}/add-fieldkit-input`}>Add Fieldkit Input</Link>
        </div>
      </div>
    )
  }
}
