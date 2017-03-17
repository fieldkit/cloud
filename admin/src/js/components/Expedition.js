// @flow weak

import React, { Component } from 'react'
import { Route, Link } from 'react-router-dom'

import { ProjectExpeditionForm } from './forms/ProjectExpeditionForm';
import { FKApiClient } from '../api/api';

import type { APIProject } from '../api/types';
import type { APIExpedition } from '../api/types';

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

  constructor(props: Props) {
    super(props);
  }

  async onExpeditionSave(name: string, description: string) {
    const expedition = await FKApiClient.get().updateExpedition(this.props.expedition.id, { name, description });
    if (expedition.type === 'ok') {
      await this.loadData();
      this.props.history.push("/");
    } else {
      return expedition.errors;
    }

    if (slug != this.props.expedition.slug) {
      this.props.onExpeditionUpdate(slug);
    } else {
      this.props.onProjectUpdate();
    }
  } 

  render() {
    const { project } = this.props;
    const projectSlug = project.slug;
    const { expedition } = this.props;
    const expeditionSlug = project.slug;

    return (

      <div className="expedition">
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