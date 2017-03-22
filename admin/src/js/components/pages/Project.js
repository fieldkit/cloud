// @flow weak

import React, { Component } from 'react'
import { Route, Link } from 'react-router-dom'
import ReactModal from 'react-modal';

import { ProjectForm } from '../forms/ProjectForm';
import { ProjectExpeditionForm } from '../forms/ProjectExpeditionForm';
import { AdministratorForm } from '../forms/AdministratorForm';
import { FKApiClient } from '../../api/api';

import type { APIProject, APIExpedition, APINewProject, APINewExpedition, APINewAdministrator } from '../../api/types';

type Props = {
  project: APIProject;
  expeditions: APIExpedition[];
  onUpdate: (newSlug: ?string) => void;

  match: Object;
  location: Object;
  history: Object;
}

export class Project extends Component {
  props: Props;
  state: {}

  constructor(props: Props) {
    super(props);

    this.state = {
      expeditions: []
    };
  }

  async onExpeditionCreate(e: APINewExpedition) {
    const { match, project } = this.props;

    const expeditionRes = await FKApiClient.get().createExpedition(project.id, e);
    if (expeditionRes.type === 'ok') {
      this.props.history.push(`${match.url}/expeditions/${e.slug}`);
    } else {
      return expeditionRes.errors;
    }
  }

  async onAdministratorAdd(e: APINewAdministrator) {
    // requires: (projectId: number, values: APINewAdministrator)

    // const { match } = this.props;
    // const memberRes = await FKApiClient.get().addMember(teamId, e);
    // if (memberRes.type === 'ok') {
    //   await this.loadTeams();
    //   this.props.history.push(`${match.url}`);
    // } else {
    //   return memberRes.errors;
    // }
  }  

  async onProjectSave(project: APINewProject) {
    // TODO: this isn't implemented on the backend yet!

    const projectRes = await FKApiClient.get().updateProject(this.props.project.id, project);
    if (projectRes.type !== 'ok') {
      return projectRes.errors;
    }

    if (projectRes.slug != this.props.project.slug && projectRes.payload) {
      this.props.onUpdate(projectRes.payload.slug);
    } else {
      this.props.onUpdate();
    }
  }

  render () {
    const { match, project } = this.props;
    const projectSlug = project.slug;

    return (
      <div className="project">
        <Route path={`${match.url}/new-expedition`} render={() =>
          <ReactModal isOpen={true} contentLabel="New expedition form">
            <h1>Create a new expedition</h1>
            <ProjectExpeditionForm
              projectSlug={projectSlug}
              onCancel={() => this.props.history.push(match.url)}
              onSave={this.onExpeditionCreate.bind(this)} />
          </ReactModal> } />

        <Route path={`${match.url}/add-administrator`} render={props =>
          <ReactModal isOpen={true} contentLabel="Add Users">
            <h1>Add Users</h1>
            <AdministratorForm
              project={project}
              onCancel={() => this.props.history.push(`${match.url}`)}
              onSave={this.onAdministratorAdd.bind(this)} />
          </ReactModal> } />          

        <div id="expeditions">
          <h4>Expeditions</h4>
          { this.props.expeditions.map((e, i) =>
            <div key={`expedition-${i}`} className="expedition-item">
              <Link to={`${match.url}/expeditions/${e.slug}`}>{e.name}</Link>
            </div> )}
          { this.props.expeditions.length === 0 &&
            <span className="empty">No expeditions!</span> }
        </div>
        <Link to={`${match.url}/new-expedition`}>Show new expedition modal</Link>

        <h2>Edit project</h2>
        <ProjectForm
          name={project ? project.name : undefined}
          slug={project ? project.slug : undefined}
          description={project ? project.description : undefined}
          onSave={this.onProjectSave.bind(this)} />
        <h3>Users</h3>
        <Link className="button secondary" to={`${match.url}/add-administrator`}>Add Users</Link>
      </div>
    )
  }
}
