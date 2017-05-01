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

import type {Attr} from '../Collection';
import type {Decorator, PointDecorator} from "../Decorators.js"
import {PointDecoratorComponent, emptyPointDecorator} from "../Decorators.js"
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
  attributes: { [string]: Attr };
  props: Props;
  state: {
    activeProject: ?APIProject
  }

  constructor(props: Props) {
    
    super(props);

    this.attributes = {
        "message": {
            name: "message",
            options: [],
            type: "string"
        },
        "user": {
            name: "user",
            options: ["@eric","@othererik","@gabriel"],
            type: "string"
        },
        "humidity": {
            name: "humidity",
            options: [],
            type: "num"
        },
        "created": {
            name: "created",
            options: [],
            type: "date"
        }
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

    return (
      <div>
        <PointDecoratorComponent initial_state={emptyPointDecorator()} attributes={this.attributes}/>
      </div>
    )
  }
}
