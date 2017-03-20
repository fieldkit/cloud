// @flow weak

import React, { Component } from 'react'
import { Route, Link } from 'react-router-dom'
import ReactModal from 'react-modal'
import _ from 'lodash'

import { ProjectExpeditionForm } from '../forms/ProjectExpeditionForm';
import { InputForm } from '../forms/InputForm';
import { FKApiClient } from '../../api/api';
import { RouteOrLoading } from '../shared/RequiredRoute';

import type { APIProject, APIExpedition, APINewExpedition, APIInput, APINewInput } from '../../api/types';

type Props = {
  project: APIProject;
  expedition: APIExpedition;
  onExpeditionUpdate: (newSlug: ?string) => void;

  match: Object;
  location: Object;
  history: Object;
}

function coerceInt(id: string): ?number {
  let intId = id ? parseInt(id, 10) : null;
  return (intId === null || _.isNaN(intId)) ? null : intId;
}

export class DataSources extends Component {
  props: Props;
  state: {
    inputs: APIInput[],
    editingInput: ?APIInput
  }

  constructor(props: Props) {
    super(props);
    this.state = {
      inputs: [],
      editingInput: null
    }

    const { inputId } = props.match.params;

    Promise.all([
      this.loadInputs(),
      this.loadSelectedInput(inputId),
    ]);
  }

  componentWillReceiveProps(nextProps: Props) {
    const promises = [];
    const stateChange = {};

    const inputId = this.editingInputId();
    const newInputId = coerceInt(nextProps.match.params.inputId);

    if (inputId != newInputId) {
      promises.push(this.loadSelectedInput(newInputId));
      stateChange.editingInput = null;
    } else if (!newInputId) {
      stateChange.editingInput = null;
    }

    if (promises.length > 0 || Object.keys(stateChange).length > 0) {
      this.setState(stateChange);
      Promise.all(promises);
    }
  }

  editingInputId(id: string = this.props.match.params.inputId): ?number {
    return coerceInt(id);
  }

  async loadInputs() {
    const inputsRes = await FKApiClient.get().getExpeditionInputs(this.props.expedition.id);
    if (inputsRes.type === 'ok' && inputsRes.payload) {
      this.setState({ inputs: inputsRes.payload.inputs || [] })
    }
  }

  async loadSelectedInput(inputId: ?number = this.editingInputId()) {
    if (inputId) {
      const inputRes = await FKApiClient.get().getInput(inputId);
      if (inputRes.type == 'ok' && inputRes.payload) {
        this.setState({ editingInput: inputRes.payload });
      }
    } else {
      this.setState({ editingInput: null });
    }
  }

  async onInputCreate(i: APINewInput) {
    const { project, expedition } = this.props;

    const inputRes = await FKApiClient.get().createInput(expedition.id, i);
    if (inputRes.type === 'ok') {
      await this.loadInputs();
      this.props.history.push(`/projects/${project.slug}/expeditions/${expedition.slug}/datasources`);
    } else {
      return inputRes.errors;
    }
  }

  async onInputSave(inputId: ?number, i: APINewInput) {
    // TODO: better error
    if (!inputId) return;

    const { project, expedition } = this.props;

    // TODO: this isn't on the server yet
    // const inputRes = await FKApiClient.get().updateInput(inputId, i);
    // if (inputRes.type === 'ok') {
    //   await this.loadInputs();
    //   this.props.history.push(`/projects/${project.slug}/expeditions/${expedition.slug}/datasources`);
    // } else {
    //   return inputRes.errors;
    // }

    this.props.history.push(`/projects/${project.slug}/expeditions/${expedition.slug}/datasources`);
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
        <RouteOrLoading
          path="/projects/:projectSlug/expeditions/:expeditionSlug/datasources/:inputId/edit"
          required={[this.state.editingInput]}
          component={() =>
            <ReactModal isOpen={true} contentLabel="Edit datasource form">
              <h1>Edit data source</h1>
              <InputForm
                input={this.state.editingInput}
                projectSlug={projectSlug}
                onCancel={() => this.props.history.push(`/projects/${projectSlug}/expeditions/${expeditionSlug}/datasources`)}
                onSave={(i) => this.onInputSave(this.editingInputId(), i)} />
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
                <th></th>
              </tr>
            { this.state.inputs.map((input, i) =>
              <tr key={`input-${i}`} className="input-item">
                <td>{input.name}</td>
                <td>Sensor</td>
                <td>Jer Thorp</td>
                <td>Water Quality</td>
                <td><Link to={`/projects/${projectSlug}/expeditions/${expeditionSlug}/datasources/${input.id}/edit`}>edit</Link></td>
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
