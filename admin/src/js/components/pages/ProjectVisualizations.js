/* @flow */
import React, { Component } from 'react'
import { Switch, Route, Link, NavLink, Redirect } from 'react-router-dom';
import Dropdown, { DropdownTrigger, DropdownContent } from 'react-simple-dropdown';
import type { Lens, Lens_ } from 'safety-lens'
import { get, set } from 'safety-lens'
import { prop } from 'safety-lens/es2015'

import type {
  Match as RouterMatch,
  Location as RouterLocation,
  RouterHistory
} from 'react-router-dom';

import type {Attr, Target, ProjectData, Collection, Filter, FilterFn} from '../../types/CollectionTypes';
import type {Decorator, PointDecorator} from "../../types/VizTypes"
import {emptyViz, emptyPointDecorator} from "../../types/VizTypes"
import {VizComponent, PointDecoratorComponent} from "../Decorators.js"
import log from 'loglevel';

import { FKApiClient } from '../../api/api';
import type { APIProject } from '../../api/types';

import { StringFilterComponent, NumFilterComponent, DateFilterComponent } from '../forms/FiltersForm'

type Bar = {"foo": string}

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
        name: "expedition",
        options: ["ITO 2015","ITO 2016","ITO 2017"],
        type: "string",
        target: "expedition"
      },
      bindings: [],
      doctypes: [],
      attributes: [
        {
          name: "message",
          options: [],
          type: "string",
          target: "attribute"
        },
        {
          name: "user",
          options: ["@eric","@othererik","@gabriel"],
          type: "string",
          target: "attribute"
        },
        {
          name: "humidity",
          options: [],
          type: "num",
          target: "attribute"
        },
        {
          name: "created",
          options: [],
          type: "date",
          target: "attribute"
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
      <div className="visualizations">
        <VizComponent initial_state={emptyViz(first_attr)} project_data={this.projectData}/>
      </div>
    )
  }
}
