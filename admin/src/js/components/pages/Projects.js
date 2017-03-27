// @flow weak

import React, { Component } from 'react'
import { Route, Link } from 'react-router-dom'
import ReactModal from 'react-modal';

import { ProjectForm } from '../forms/ProjectForm';
import { FKApiClient } from '../../api/api';
import { joinPath } from '../../common/util';
import type { APIProject, APINewProject } from '../../api/types';

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

        <h1>Projects</h1>
        <div id="projects" className="gallery-list">
        { projects.map((p, i) =>
          <Link to={joinPath(match.url, 'projects', p.slug)} key={`project-${i}`} className="gallery-list-item">
            <h4>{p.name}</h4>
          </Link> )}
        { projects.length === 0 &&
          <span className="empty">No projects!</span> }
        </div>

        <Link className="button" to={joinPath(match.url, 'new-project')}>Create new project</Link>
      </div>
    )
  }
}
