// @flow weak

import React, { Component } from 'react'
import { Route, Link } from 'react-router-dom'
import ReactModal from 'react-modal';

import { ProjectForm } from './forms/ProjectForm';
import { ProjectExpeditionForm } from './forms/ProjectExpeditionForm';
import { FKApiClient } from '../api/api';

import '../../css/home.css'

type Props = {
  match: Object;
  location: Object;
  history: Object;
}

export class Project extends Component {
  props: Props;
  state: {
    project: ?Object,
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
    const projectSlug = this.props.match.params.projectSlug;

    const projectRes = await FKApiClient.get().getProjectBySlug(projectSlug);
    if (projectRes.type === 'ok') {
      this.setState({ project: projectRes.payload })
    }

    const expeditionsRes = await FKApiClient.get().getExpeditionsByProjectSlug(projectSlug);
    if (expeditionsRes.type === 'ok') {
      this.setState({ expeditions: expeditionsRes.payload || [] })
    }
  }

  async onExpeditionCreate(name: string, description: string) {

  }

  async onProjectSave(name: string, description: string) {
    // TODO: this isn't implemented on the backend yet!

    // const project = await FKApiClient.get().saveProject(slug, description);
    // if (project.type === 'ok') {
    //   await this.loadData();
    //   this.props.history.push("/");
    // } else {
    //   return project.errors;
    // }
  }

  render () {
    const { project } = this.state;
    const projectSlug = this.props.match.params.projectSlug;

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
          { this.state.expeditions.length == 0 &&
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
