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
import { RemoveIcon, EditIcon } from '../icons/Icons'

import type { APIProject, APIExpedition, APINewExpedition, APIInputs, APITwitterInputCreateResponse, APINewTwitterInput, APINewFieldkitInput, APITeam, APIMember, APIUser, APIMutableInput } from '../../api/types';

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
    members: APIMember[],
    users: {[id: number]: APIUser},
  }

  constructor(props: Props) {
    super(props);
    this.state = {
      inputs: {},
      teams: [],
      members: [],
      users: {}
    }

    this.loadInputs();
    this.loadTeams();
  }

  async loadInputs() {
    const inputsRes = await FKApiClient.get().getExpeditionInputs(this.props.expedition.id);
    if (inputsRes.type === 'ok') {
      this.setState({ inputs: inputsRes.payload });
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
      const members = this.state.members.slice();

      membersRes.payload.members.forEach(member =>
        !members.find(m => m.user_id === member.user_id) &&
          members.push(member)
      );
      this.setState({members: members});
      for (const member of members) {
        await this.loadMemberName(teamId, member.user_id);
      }
    }
  }

  async loadMemberName(teamId: number, userId: number){
    const userRes = await FKApiClient.get().getUserById(userId);
    if (userRes.type === 'ok' && userRes.payload) {
      const { users } = this.state;
      const userId = userRes.payload.id;
      users[userId] = userRes.payload;
      this.setState({users: users});
    }
  }  

  async onTwitterCreate(t: APINewTwitterInput) {
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
    if (fieldkitRes.type === 'ok') {
      await this.loadInputs();
      this.props.history.push(`${match.url}`);
    } else {
      return fieldkitRes.errors
    }
  }

  async onInputUpdate(inputId: number, inputData: APIMutableInput) {
    // updateInput(): Promise<FKAPIResponse<APIBaseInput>> {
    // const teamRes = await FKApiClient.get().updateTeam(teamId, team);
    // if(teamRes.type === 'ok' && teamRes.payload) {
    //   await this.loadTeams();
    // } else {
    //   return teamRes.errors;
    // }
  }

  getInputById (id: number): ?APIMutableInput {
    const { inputs } = this.state;
    for(let inputType in inputs){
      if(inputs[inputType]){
        const twitterInput = inputs[inputType].find(input => input.id === id);
        if(twitterInput){
          const { name, team_id, user_id } = twitterInput;  
          return {
            name: name,
            team_id: team_id,
            user_id: user_id
          };        
        }
      }
    }
  }

  render() {
    const { expedition, match } = this.props;
    const { twitter_account_inputs, fieldkit_inputs } = this.state.inputs;
    const { users, teams } = this.state;

    return (
      <div className="data-sources-page">

      <Route path={`${match.url}/add-fieldkit-input`} render={props =>
        <ReactModal isOpen={true} contentLabel="Add Fieldkit Input" className="modal" overlayClassName="modal-overlay">
          <h2>Add Fieldkit Input</h2>
          <InputForm
            expeditionId={props.match.params.expeditionId}
            users={users}
            teams={teams}
            onCancel={() => this.props.history.push(`${match.url}`)}
            onSave={this.onFieldkitCreate.bind(this)}
            saveText="Add" />
        </ReactModal> } />

      <Route path={`${match.url}/add-twitter-input`} render={props =>
        <ReactModal isOpen={true} contentLabel="Add Twitter Input" className="modal" overlayClassName="modal-overlay">
          <h2>Add Twitter Input</h2>
          <InputForm
            expeditionId={props.match.params.expeditionId}
            users={users}
            teams={teams}
            onCancel={() => this.props.history.push(`${match.url}`)}
            onSave={this.onTwitterCreate.bind(this)}
            saveText="Add" />
        </ReactModal> } />

      <Route path={`${match.url}/:inputId/edit`} render={props =>
        <ReactModal isOpen={true} contentLabel="Edit Input" className="modal" overlayClassName="modal-overlay">
          <h2>Edit Input</h2>
          <InputForm
            expeditionId={props.match.params.expeditionId}
            users={users}
            teams={teams}
            input={this.getInputById(parseInt(props.match.params.inputId))}
            onCancel={() => this.props.history.push(`${match.url}`)}
            onSave={ (props.match.params.inputType === 'twitter') ? this.onTwitterCreate.bind(this) : this.onFieldkitCreate.bind(this)} 
            saveText="Add" />
        </ReactModal> } />        

        <h1>Data Sources</h1>

        <div className="input-section">
          <h3>Twitter</h3>
          { twitter_account_inputs && twitter_account_inputs.length > 0 &&
            <table className="twitter-account-inputs-table">
              <thead>
                <tr>
                  <th>Username</th>
                  <th>Binding</th>
                  <th></th>
                  <th></th>                  
                </tr>
              </thead>                
              <tbody>
              { twitter_account_inputs.map((t, i) =>
                <tr key={i} className="input-item">
                  <td>{t.name}</td>
                  <td>None</td>
                  <td>
                    <Link className="bt-icon medium" to={`${match.url}/${t.id}/edit`}>
                      <EditIcon />
                    </Link>
                  </td>
                  <td>
                    <div className="bt-icon medium">
                      <RemoveIcon />
                    </div>
                  </td>
                </tr> )}
              </tbody>
            </table> }
          <Link className="button" to={`${match.url}/add-twitter-input`}>Add Twitter Account</Link>

          <br/>
          <br/>          

          <h3>Fieldkit Sensors</h3>
          { fieldkit_inputs && fieldkit_inputs.length > 0 &&
            <table className="fieldkit-inputs-table">
              <thead>
                <tr>
                  <th>Name</th>
                  <th>Binding</th>
                  <th></th>
                  <th></th>                  
                </tr>
              </thead>                
              <tbody>                
              { fieldkit_inputs.map((f, i) =>
                <tr key={i} className="input-item">
                  <td>{f.name}</td>
                  <td>None</td>
                  <td>
                    <Link className="bt-icon medium" to={`${match.url}/${f.id}/edit`}>
                      <EditIcon />
                    </Link>
                  </td>
                  <td>
                    <div className="bt-icon medium">
                      <RemoveIcon />
                    </div>
                  </td>
                </tr> )}
              </tbody>
            </table> }

          <Link className="button" to={`${match.url}/add-fieldkit-input`}>Add Fieldkit Input</Link>
        </div>
      </div>
    )
  }
}
