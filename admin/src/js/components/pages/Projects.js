// @flow weak

import React, { Component } from 'react'
import { Route, Link } from 'react-router-dom'
import ReactModal from 'react-modal';

import { ProjectForm } from '../forms/ProjectForm';
import { FKApiClient } from '../../api/api';
import type { APINewProject } from '../../api/types';

import '../../../css/projects.css'

type Props = {
  match: Object;
  location: Object;
  history: Object;
}

export class Projects extends Component {
  props: Props;
  state: {
    projects: Object[]
  }

  constructor(props: Props) {
    super(props);

    this.state = {
      projects: []
    };

    this.loadData();
  }

  async loadData() {
    const projectsRes = await FKApiClient.get().getProjects();
    if (projectsRes.type === 'ok' && projectsRes.payload) {
      this.setState({ projects: projectsRes.payload.projects })
    }
  }

  async onProjectCreate(p: APINewProject) {
    const project = await FKApiClient.get().createProject(p);
    if (project.type === 'ok') {
      await this.loadData();
      this.props.history.push("/");
    } else {
      return project.errors;
    }
  }

  render () {
    const { match } = this.props;

    return (
      <div className="projects">
        <Route path={`${match.url}/new-project`} render={() =>
          <ReactModal isOpen={true} contentLabel="New project form">
            <h1>Create a new project</h1>
            <ProjectForm
              onCancel={() => this.props.history.push("/")}
              onSave={this.onProjectCreate.bind(this)} />
          </ReactModal> } />

        <div id="projects">
        { this.state.projects.map((p, i) =>
          <div key={`project-${i}`} className="project-item">
            <Link to={`${match.url}/projects/${p.slug}`}>{p.name}</Link>
          </div> )}
        { this.state.projects.length === 0 &&
          <span className="empty">No projects!</span> }
        </div>

        <Link to={`${match.url}/new-project`}>Show new project modal</Link>
      </div>
    )
  }
}
