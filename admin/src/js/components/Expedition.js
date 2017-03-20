// @flow weak

import React, { Component } from 'react'
import { Route, Link } from 'react-router-dom'
import ReactModal from 'react-modal'

import { ProjectExpeditionForm } from './forms/ProjectExpeditionForm';
import { InputForm } from './forms/InputForm';
import { FKApiClient } from '../api/api';

import type { APIProject, APIExpedition, APINewExpedition, APIInput, APINewInput } from '../api/types';

import '../../css/home.css'

type Props = {
  project: APIProject;
  expedition: APIExpedition;
  onExpeditionUpdate: (newSlug: ?string) => void;

  match: Object;
  location: Object;
  history: Object;
}

export class Expedition extends Component {
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

  async onExpeditionSave(expedition: APINewExpedition) {
    const expeditionRes = await FKApiClient.get().updateExpedition(this.props.expedition.id, expedition);
    if (expeditionRes.type !== 'ok') {
      return expeditionRes.errors;
    }

    if (expeditionRes.slug != this.props.expedition.slug && expeditionRes.payload) {
      this.props.onExpeditionUpdate(expeditionRes.payload.slug);
    } else {
      this.props.onExpeditionUpdate();
    }
  }

  async onInputCreate(i: APINewInput) {
    const { project, expedition } = this.props;

    const inputRes = await FKApiClient.get().createInput(expedition.id, i);
    if (inputRes.type === 'ok') {
      await this.loadData();
      this.props.history.push(`/projects/${project.slug}/expeditions/${expedition.slug}`);
    } else {
      return inputRes.errors;
    }
  }

  render() {
    const { project, expedition } = this.props;
    const projectSlug = project.slug;
    const expeditionSlug = expedition.slug;

    return (
      <div className="expedition">
        <Route path="/projects/:projectSlug/expeditions/:expeditionSlug/new-input" render={() =>
          <ReactModal isOpen={true} contentLabel="New input form">
            <h1>Create a new input</h1>
            <InputForm
              projectSlug={projectSlug}
              onCancel={() => this.props.history.push(`/projects/${projectSlug}/expeditions/${expeditionSlug}`)}
              onSave={this.onInputCreate.bind(this)} />
          </ReactModal> } />

        <div id="inputs">
          <h4>Inputs</h4>
          { this.state.inputs.map((input, i) =>
            <div key={`input-${i}`} className="input-item">
              { JSON.stringify(input) }
            </div> )}
          { this.state.inputs.length === 0 &&
            <span className="empty">No inputs!</span> }
        </div>
        <Link to={`/projects/${projectSlug}/expeditions/${expeditionSlug}/new-input`}>Show new input modal</Link>

        <h1>Expedition Settings</h1>
        <ProjectExpeditionForm
          projectSlug={projectSlug}
          name={expedition.name}
          slug={expedition.slug}
          description={expedition.description}
          onSave={this.onExpeditionSave.bind(this)} />
      </div>
    )
  }
}
