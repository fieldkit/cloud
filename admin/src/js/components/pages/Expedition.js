// @flow weak

import React, { Component } from 'react'
import { Route, Link } from 'react-router-dom'
import ReactModal from 'react-modal'

import { ProjectExpeditionForm } from '../forms/ProjectExpeditionForm';
import { FKApiClient } from '../../api/api';

import type { APIProject, APIExpedition, APINewExpedition } from '../../api/types';

type Props = {
  project: APIProject;
  expedition: APIExpedition;
  onUpdate: (newSlug: ?string) => void;

  match: Object;
  location: Object;
  history: Object;
}

export class Expedition extends Component {
  props: Props;

  async onExpeditionSave(expedition: APINewExpedition) {
    const expeditionRes = await FKApiClient.get().updateExpedition(this.props.expedition.id, expedition);
    if (expeditionRes.type !== 'ok') {
      return expeditionRes.errors;
    }

    if (expeditionRes.slug != this.props.expedition.slug && expeditionRes.payload) {
      this.props.onUpdate(expeditionRes.payload.slug);
    } else {
      this.props.onUpdate();
    }
  }

  render() {
    const { project, expedition } = this.props;
    const projectSlug = project.slug;
    const expeditionSlug = expedition.slug;

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
