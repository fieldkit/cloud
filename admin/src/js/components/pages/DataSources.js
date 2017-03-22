// @flow weak

import React, { Component } from 'react'
import { Route, Link, Redirect } from 'react-router-dom'
import ReactModal from 'react-modal'
import _ from 'lodash'
import log from 'loglevel';

import { ProjectExpeditionForm } from '../forms/ProjectExpeditionForm';
// import { InputForm } from '../forms/InputForm';
import { FKApiClient } from '../../api/api';
import { RouteOrLoading } from '../shared/RequiredRoute';

import type { APIProject, APIExpedition, APINewExpedition, APIInputs, APITwitterInputCreateResponse } from '../../api/types';

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
    inputs: APIInputs
  }

  constructor(props: Props) {
    super(props);
    this.state = {
      inputs: {}
    }

    this.loadInputs();
  }

  async loadInputs() {
    const inputsRes = await FKApiClient.get().getExpeditionInputs(this.props.expedition.id);
    if (inputsRes.type === 'ok') {
      this.setState({ inputs: inputsRes.payload })
    }
  }

  async onTwitterCreate(event) {
    const res = await FKApiClient.get().createTwitterInput(this.props.expedition.id);
    if (res.type === 'ok') {
      const redirect = res.payload.location;
      window.location = redirect;
    } else {
      log.error('Bad request to create twitter input!');
    }
  }

  render() {
    const { twitter_accounts } = this.state.inputs;

    return (
      <div className="data-sources-page">
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
        </div>
      </div>
    )
  }
}
