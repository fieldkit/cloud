// @flow weak

import React, { Component } from 'react'
import { Route, Link } from 'react-router-dom'
import ReactModal from 'react-modal';

import { ProjectForm } from '../forms/ProjectForm';
import { ProjectExpeditionForm } from '../forms/ProjectExpeditionForm';
import { AdministratorForm } from '../forms/AdministratorForm';
import { FKApiClient } from '../../api/api';

import removeImg from '../../../img/icons/icon-remove-small.png'
import '../../../css/projects.css'

import type { APIProject, APINewProject, APINewAdministrator, APIAdministrator } from '../../api/types';

type Props = {
  project: APIProject;
  administrators: APIAdministrator[];
  onUpdate: (newSlug: ?string) => void;

  match: Object;
  location: Object;
  history: Object;
}

export class ProjectSettings extends Component {
  props: Props;
  state: {
    administrators: APIAdministrator[]
  }

  constructor(props: Props) {
    super(props);

    this.state = {
      administrators: []
    };

    this.loadAdministrators();
  }

  async loadAdministrators() {
    const { project } = this.props;
    const administratorsRes = await FKApiClient.get().getAdministrators(project.id);  
    if (administratorsRes.type === 'ok' && administratorsRes.payload) {
      const administrators = administratorsRes.payload.administrators;
      this.setState({administrators: administrators});
    }
  }

  async onAdministratorAdd(e: APINewAdministrator) {
    const { project, match } = this.props;
    const administratorRes = await FKApiClient.get().addAdministrator(project.id, e);
    if (administratorRes.type === 'ok') {
      await this.loadAdministrators();
      this.props.history.push(`${match.url}`);
    } else {
      return administratorRes.errors;
    }
  }

  async onAdministratorDelete(e: APIAdministrator) {
    const { match } = this.props;
    const administratorRes = await FKApiClient.get().deleteAdministrator(e.project_id, e.user_id);
    if (administratorRes.type === 'ok') {
      await this.loadAdministrators();
      this.props.history.push(`${match.url}`);
    } else {
      return administratorRes.errors;
    }
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

        <Route path={`${match.url}/add-administrator`} render={props =>
          <ReactModal isOpen={true} contentLabel="Add Users">
            <h2>Add Users</h2>
            <AdministratorForm
              project={project}
              administrators={this.state.administrators}
              onCancel={() => this.props.history.push(`${match.url}`)}
              onSave={this.onAdministratorAdd.bind(this)} />
          </ReactModal> } />

        <h1>Project Settings</h1>
        <div className="row">
          <div className="two-columns">
            <h3>Main</h3>
            <ProjectForm
              name={project ? project.name : undefined}
              slug={project ? project.slug : undefined}
              description={project ? project.description : undefined}

              onSave={this.onProjectSave.bind(this)} />
          </div>
          <div className="two-columns">        
            <h3>Users</h3>
            <p>
              Users you add to this project have administrative rights. They can create a new expedition, change its settings and add members to it. If you’re trying to add a member to an expedition instead, select an expedition and then go to <i>Teams</i>.
            </p>

            <table className="administrators-table">
              <thead>
                <tr>
                  <th></th>
                  <th>Users ({ this.state.administrators.length })</th>
                  <th></th>
                </tr>
              </thead>
              <tbody>
              { this.state.administrators.map(administrator =>
                <tr>
                  <td>
                    <div className="user-avatar">
                      <img />
                    </div>                    
                  </td>
                  <td>
                    {administrator.user_id}
                  </td>
                  <td>
                    <div className="bt-icon" onClick={this.onAdministratorDelete.bind(this, administrator)}>
                      <img src={removeImg} alt="external link" />
                    </div>
                  </td>
                </tr> )}
              </tbody>
            </table>

            <Link className="button secondary" to={`${match.url}/add-administrator`}>Add Users</Link>
          </div>
        </div>
      </div>
    )
  }
}
