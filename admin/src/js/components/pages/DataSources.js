// @flow weak

import React, { Component } from 'react'
import { Route, Link } from 'react-router-dom'
import ReactModal from 'react-modal'

import { ProjectExpeditionForm } from '../forms/ProjectExpeditionForm';
import { InputForm } from '../forms/InputForm';
import { FKApiClient } from '../../api/api';

import type { APIProject, APIExpedition, APINewExpedition, APIInput, APINewInput } from '../../api/types';

type Props = {
  project: APIProject;
  expedition: APIExpedition;
  onExpeditionUpdate: (newSlug: ?string) => void;

  match: Object;
  location: Object;
  history: Object;
}

export class DataSources extends Component {
  props: Props;
  state: {
    inputs: APIInput[]
  }

  constructor(props: Props) {
    super(props);
    this.state = {
      inputs: []
    }
    this.loadData();
  }

  async loadData() {
    const inputsRes = await FKApiClient.get().getExpeditionInputs(this.props.expedition.id);
    if (inputsRes.type === 'ok' && inputsRes.payload) {
      this.setState({ inputs: inputsRes.payload.inputs || [] })
    }
  }

  async onInputCreate(i: APINewInput) {
    const { project, expedition } = this.props;

    const inputRes = await FKApiClient.get().createInput(expedition.id, i);
    if (inputRes.type === 'ok') {
      await this.loadData();
      this.props.history.push(`/projects/${project.slug}/expeditions/${expedition.slug}/datasources`);
    } else {
      return inputRes.errors;
    }
  }

  render() {
    const { project, expedition } = this.props;
    const projectSlug = project.slug;
    const expeditionSlug = expedition.slug;

    return (
      <div className="inputs">
        <Route path="/projects/:projectSlug/expeditions/:expeditionSlug/datasources/new-datasource" render={() =>
          <ReactModal isOpen={true} contentLabel="New datasource form">
            <h1>Create a new data source</h1>
            <InputForm
              projectSlug={projectSlug}
              onCancel={() => this.props.history.push(`/projects/${projectSlug}/expeditions/${expeditionSlug}/datasources`)}
              onSave={this.onInputCreate.bind(this)} />
          </ReactModal> } />

        <h1>Data Sources</h1>
        { this.state.inputs &&
          <table>
            <tbody>
              <tr>
                <th>Name</th>
                <th>Type</th>
                <th>Binding</th>
                <th>Inputs</th>
              </tr>
            { this.state.inputs.map((input, i) =>
              <tr key={`input-${i}`} className="input-item">
                <td>{input.name}</td>
                <td>Sensor</td>
                <td>Jer Thorp</td>
                <td>Water Quality</td>
              </tr> )}
            </tbody>
          </table> }
        { this.state.inputs.length === 0 &&
          <span className="empty">No inputs!</span> }

      <Link to={`/projects/${projectSlug}/expeditions/${expeditionSlug}/datasources/new-datasource`}>Show new datasource modal</Link>
      </div>
    )
  }
}
