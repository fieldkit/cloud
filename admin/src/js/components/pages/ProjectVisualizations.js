/* @flow */
import React, { Component } from 'react';
import log from 'loglevel';
import { FKApiClient } from '../../api/api';
import VizForm from '../forms/VizForm.js';
import type {
  Match as RouterMatch,
  Location as RouterLocation,
  RouterHistory
} from 'react-router-dom';
import type { APIProject } from '../../api/types';
import type { ProjectData } from '../../types/CollectionTypes';
import type { Viz, PointDecorator } from '../../types/VizTypes';

function emptyViz(a: Attr): Viz{
  return {
    grouping_operation: {
      operation: 'equal',
      parameter: null,
      source_attribute: a
    },
    selection_operations: [],
    decorator: emptyPointDecorator(),
    source_collections: []
  }
}

function emptyPointDecorator(): PointDecorator{
  return {
    points: {
      color: {
        type: 'constant',
        colors: [{location: 0, color: '#ff0000'}],
        data_key: null,
        bounds: null
      },
      size: {
        type: 'constant',
        data_key: null,
        bounds: [15,15]
      },
      sprite: 'circle.png'
    },
    title: '',
    type: 'point'
  }
}

type Props = {
  match: RouterMatch;
  location: RouterLocation;
  history: RouterHistory;
}

export class ProjectVisualizations extends Component {
  projectData: ProjectData;
  props: Props;
  state: {
    activeProject: ?APIProject
  }

  constructor(props: Props) {

    super(props);

    this.projectData = {
      expeditions: {
        name: 'expedition',
        options: ['ITO 2015','ITO 2016','ITO 2017'],
        type: 'string',
        target: 'expedition'
      },
      bindings: [],
      doctypes: [],
      attributes: [
        {
          name: 'message',
          options: [],
          type: 'string',
          target: 'attribute'
        },
        {
          name: 'user',
          options: ['@eric','@othererik','@gabriel'],
          type: 'string',
          target: 'attribute'
        },
        {
          name: 'humidity',
          options: [],
          type: 'num',
          target: 'attribute'
        },
        {
          name: 'created',
          options: [],
          type: 'date',
          target: 'attribute'
        }
      ]
    }
    this.state = {
      activeProject: null
    }

    this.loadActiveProject(this.projectSlug());
  }

  componentWillReceiveProps(nextProps: Props) {
  }

  projectSlug(): ?string {
    return this.props.match.params.projectSlug;
  }

  async loadActiveProject(projectSlug: ?string ) {
    log.debug('main -> loadActiveProject', projectSlug);
    if (projectSlug) {
      const projectRes = await FKApiClient.get().getProjectBySlug(projectSlug);
      if (projectRes.type === 'ok') {
        this.setState({ activeProject: projectRes.payload });
      }
    } else {
      this.setState({ activeProject: null });
    }
  }

  handleLinkClick() {
    this.refs.dropdown.hide();
  }

  render() {
    const { activeProject } = this.state;
    const first_attr = this.projectData.attributes[0]

    return (
      <div className='visualizations'>
        <VizForm initial_state={emptyViz(first_attr)} project_data={this.projectData}/>
      </div>
    )
  }
}
