// @flow weak

import React, { Component } from 'react'
import { Route, Link } from 'react-router-dom'
import ReactModal from 'react-modal';

import { ProjectForm } from './forms/ProjectForm';
import { ProjectExpeditionForm } from './forms/ProjectExpeditionForm';
import { FKApiClient } from '../api/api';

import type { APIProject } from '../api/types';

import '../../css/home.css'

type Props = {
  project: APIProject;
  onProjectUpdate: (newSlug: ?string) => void;

  match: Object;
  location: Object;
  history: Object;
}

export class Project extends Component {
  props: Props;
  state: {
    expeditions: Object[]
  }

  constructor(props: Props) {
    super(props);

    this.state = {
      project: null,
      expeditions: []
    };

    this.loadData();
  }

  async loadData() {
    const expeditionsRes = await FKApiClient.get().getExpeditionsByProjectSlug(this.props.project.slug);
    if (expeditionsRes.type === 'ok' && expeditionsRes.payload) {
      this.setState({ expeditions: expeditionsRes.payload.expeditions || [] })
    }
  }

  async onExpeditionCreate(name: string, slug: string, description: string) {
    const { project } = this.props;

    const expeditionRes = await FKApiClient.get().createExpedition(project.id, { name, slug, description });
    if (expeditionRes.type === 'ok') {
      await this.loadData();
      this.props.history.push(`/projects/${project.slug}`);
    } else {
      return expeditionRes.errors;
    }
  }

  async onProjectSave(name: string, slug: string, description: string) {
    // TODO: this isn't implemented on the backend yet!

    const project = await FKApiClient.get().updateProject(this.props.project.id, { name, slug, description });
    if (project.type === 'ok') {
      await this.loadData();
      this.props.history.push("/");
    } else {
      return project.errors;
    }

    if (slug != this.props.project.slug) {
      this.props.onProjectUpdate(slug);
    } else {
      this.props.onProjectUpdate();
    }
  }

  render () {
    const { project } = this.props;
    const projectSlug = project.slug;

    return (
      <div className="project">
        <Route path="/projects/:projectSlug/new-expedition" render={() =>
          <ReactModal isOpen={true} contentLabel="New expedition form">
            <h1>Create a new expedition</h1>
            <ProjectExpeditionForm
              projectSlug={projectSlug}
              onCancel={() => this.props.history.push(`/projects/${projectSlug}`)}
              onSave={this.onExpeditionCreate.bind(this)} />
          </ReactModal> } />

        <div id="expeditions">
          { this.state.expeditions.map((e, i) =>
            <div key={`expedition-${i}`} className="expedition-item">
              {JSON.stringify(e)}
            </div> )}
          { this.state.expeditions.length === 0 &&
            <span className="empty">No expeditions!</span> }
        </div>
        <Link to={`/projects/${projectSlug}/new-expedition`}>Show new expedition modal</Link>

        <h2>Edit project</h2>
        <ProjectForm
          name={project ? project.name : undefined}
          description={project ? project.description : undefined}
          onSave={this.onProjectSave.bind(this)} />
      </div>
    )
  }
}
