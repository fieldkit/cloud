// @flow weak

import React, { Component } from 'react'
import { Route, Link } from 'react-router-dom'
import ReactModal from 'react-modal';

import { ProjectForm } from './forms/ProjectForm';
import { FKApiClient } from '../api/api';

import '../../css/home.css'

type Props = {
  newProjectModal?: boolean;
  match: Object;
  location: Object;
  history: Object;
}

export class Home extends Component {
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
    if (projectsRes.type === 'ok') {
      this.setState({ projects: projectsRes.payload })
    }
  }

  async onProjectCreate(name: string, description: string) {
    const project = await FKApiClient.get().createProject(name, description);
    if (project.type === 'ok') {
      await this.loadData();
      this.props.history.push("/");
    } else {
      return project.errors;
    }
  }

  render () {
    return (
      <div className="main">
        <Route path="/new-project" render={() =>
          <ReactModal
            isOpen={true}
            contentLabel="New project form">
            <ProjectForm
              onCancel={() => this.props.history.push("/")}
              onSave={this.onProjectCreate.bind(this)} />
          </ReactModal> } />

        <div id="projects">
        { this.state.projects.map((p, i) =>
          <div key={`project-${i}`} className="project">
            {JSON.stringify(p)}
          </div> )}
        { this.state.projects.length == 0 &&
          <span className="empty">No projects!</span> }
        </div>

        <Link to="/new-project">Show new project modal</Link>
      </div>
    )
  }
}
