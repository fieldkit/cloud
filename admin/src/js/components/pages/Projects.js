// @flow weak

import React, { Component } from 'react'
import { Route, Link } from 'react-router-dom'
import ReactModal from 'react-modal';

import { ProjectForm } from '../forms/ProjectForm';
import { FKApiClient } from '../../api/api';
import { joinPath } from '../../common/util';
import type { APIProject, APINewProject } from '../../api/types';

import '../../../css/projects.css'

type Props = {
  projects: APIProject[],

  match: Object;
  location: Object;
  history: Object;
}

export class Projects extends Component {
  props: Props;

  async onProjectCreate(p: APINewProject) {
    const project = await FKApiClient.get().createProject(p);
    if (project.type === 'ok') {
      this.props.history.push("/");
    } else {
      return project.errors;
    }
  }

  render () {
    const { projects, match } = this.props;

    return (
      <div className="projects">
        <Route path={joinPath(match.url, 'new-project')} render={() =>
          <ReactModal isOpen={true} contentLabel="New project form">
            <h1>Create a new project</h1>
            <ProjectForm
              onCancel={() => this.props.history.push("/")}
              onSave={this.onProjectCreate.bind(this)} />
          </ReactModal> } />

        <div id="projects">
        { projects.map((p, i) =>
          <div key={`project-${i}`} className="project-item">
            <Link to={joinPath(match.url, 'projects', p.slug, 'projectexpeditions')}>{p.name}</Link>
          </div> )}
        { projects.length === 0 &&
          <span className="empty">No projects!</span> }
        </div>

        <Link to={joinPath(match.url, 'new-project')}>Show new project modal</Link>
      </div>
    )
  }
}
