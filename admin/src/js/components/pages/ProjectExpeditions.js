// @flow weak

import React, { Component } from 'react'
import { Route, Link } from 'react-router-dom'
import ReactModal from 'react-modal';

import { ProjectForm } from '../forms/ProjectForm';
import { ProjectExpeditionForm } from '../forms/ProjectExpeditionForm';
import { FKApiClient } from '../../api/api';

import type { APIProject, APIExpedition, APINewProject, APINewExpedition, } from '../../api/types';

type Props = {
  project: APIProject;
  expeditions: APIExpedition[];
  onUpdate: (newSlug: ?string) => void;

  match: Object;
  location: Object;
  history: Object;
}

export class ProjectExpeditions extends Component {
  props: Props;
  state: {
    expeditions: APIExpedition[],
    expeditionPictures: {[expeditionId: number]: string}    
  }

  constructor(props: Props) {
    super(props);

    this.state = {
      expeditions: [],
      expeditionPictures: {}
    };
    this.loadExpeditionPictures();
  }

  onComponentWillReceiveProps(nextProps: Props) {
    this.setState({expeditions: nextProps.expeditions});
    this.loadExpeditionPictures();
  }

  async loadExpeditionPictures (){
    const { expeditions } = this.props;
    const { expeditionPictures } = this.state;

    for (const expedition of expeditions) {
      const expeditionRes = await FKApiClient.get().expeditionPictureUrl(expedition.id);
      if (expeditionRes) {
        expeditionPictures[expedition.id] = expeditionRes;
      }
    }
    this.setState({expeditionPictures: expeditionPictures});        
  }

  async onExpeditionCreate(e: APINewExpedition) {
    const { match, project } = this.props;

    const expeditionRes = await FKApiClient.get().createExpedition(project.id, e);
    if (expeditionRes.type === 'ok') {
      this.props.history.push(`${match.url}/expeditions/${e.slug}/datasources`);
    } else {
      return expeditionRes.errors;
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
    const { expeditionPictures } = this.state;    

    return (
      <div className="project">
        <Route path={`${match.url}/new-expedition`} render={() =>
          <ReactModal isOpen={true} contentLabel="New expedition form" className="modal" overlayClassName="modal-overlay">
            <h2>Create New Expedition</h2>
            <ProjectExpeditionForm
              projectSlug={projectSlug}
              onCancel={() => this.props.history.push(match.url)}
              onSave={this.onExpeditionCreate.bind(this)} />
          </ReactModal> } />
        
        <h1>Expeditions</h1>
        <div id="expeditions" className="gallery-list">
          { this.props.expeditions.map((e, i) =>
            <Link to={`/projects/${projectSlug}/expeditions/${e.slug}/datasources`} key={`expedition-${i}`} className="gallery-list-item expeditions"
              style={{
                backgroundImage: 'url(' + expeditionPictures[e.id] + ')'
              }}>
              <h4>{e.name}</h4>
            </Link> )}
          { this.props.expeditions.length === 0 &&
            <span className="empty">You don't have any expeditions yet.</span> }
        </div>
        <Link to={`${match.url}/new-expedition`} className="button">Create new expedition</Link>
      </div>
    )
  }
}
