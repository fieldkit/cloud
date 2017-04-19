/* @flow */

import React, { Component } from 'react'
import { Switch, Route, Link, NavLink, Redirect } from 'react-router-dom';
import Dropdown, { DropdownTrigger, DropdownContent } from 'react-simple-dropdown';

import type {
  Match as RouterMatch,
  Location as RouterLocation,
  RouterHistory
} from 'react-router-dom';

import log from 'loglevel';

import { FKApiClient } from '../api/api';
import type { APIUser, APIProject, APIExpedition } from '../api/types';
import { fkHost } from '../common/util';

import { RouteOrLoading } from './shared/RequiredRoute';
import { Projects } from './pages/Projects';
import { ProjectExpeditions } from './pages/ProjectExpeditions';
import { ProjectSettings } from './pages/ProjectSettings';
import { Expedition } from './pages/Expedition'
import { Teams } from './pages/Teams';
import { DataSources } from './pages/DataSources';
import { Profile } from './pages/Profile';

import fieldkitLogo from '../../img/logos/fieldkit-logo-red.svg';
import placeholderImage from '../../img/profile_placeholder.svg'
import externalLinkImg from '../../img/icons/icon-external-link.png'
import '../../css/main.css'

import { HamburgerIcon, OpenInNewIcon, ArrowDownIcon } from './icons/Icons'

type Props = {
  match: RouterMatch;
  location: RouterLocation;
  history: RouterHistory;
}

type Collection = {
    name: string;
    id: string;
    filters: number[];
    guid_filters: GuidFilter[];      
    string_filters: StringFilter[];
	num_filters: NumFilter[];
	geo_filters: GeoFilter[];
	date_filters: DateFilter[];
	mods: number[];
	string_mods: StringMod[];
	num_mods: NumMod[];
	geo_mods: GeoMod[];
}

type GuidFilter = {
    id: number;
    attribute: string;
    operation: "contains";
    query: number[];
}

type StringFilter = {
    id: number;
    attribute: string;
    operation: "contains" | "does not contain" | "matches" | "exists";
    query: string;
}

type NumFilter = {
    id: number;
    attribute: string;
    operation: "GT" | "LT" | "EQ" | "notch";
    query: number;
}

type GeoFilter = {
    id: number;
    attribute: string;
    operation: "within" | "not within";
    query: Object
}
type DateFilter = {
    id: number;
    attribute: string;
    operation: "before" | "after" | "within";
    date: number;
    within: number;
}
type StringMod = {
    id: number;
    attribute: string;
    operation: "gsub";
    query: string;
}
type NumMod = {
    id: number;
    attribute: string;
    operation: "round";
}
type GeoMod = {
    id: number;
    attribute: string;
    operation: "copy from" | "jitter";
    source_collection: string;
    filters: number[];
}

export class StringFilterComponent extends Component {
    props: {
        data: StringFilter;
        options: [string, string][];
    }

    constructor(props: StringFilter) {
        super(props);
    }

    render() {
        const operations = ["contains","does not contain","matches","exists"].map(o => <option value={o}>{o.toUpperCase()}</option>)
        const data = this.props.data
        let value_field 
        if(this.props.options.length > 0){
            let options = this.props.options.map(([id,name]) => {
                return (
                    <option value={name}>{name}</option>
                )
            })
            value_field = (
                <select className="value-body-select" value={data.query}>
                    {options}
                </select>
            )
        } else {
            value_field = <input className="filter-body-input" value={data.query}/>
        }
        
        return (
            <div className="fk-filter fk-guidfilter">
                <div className="filter-title-bar">
                    <span>{data.attribute}</span>
                    <div className="filter-title-controls">
                        <span className="filter-icon"></span>
                        <span className="filter-closer"></span>
                    </div>
                </div>
                <div className="filter-body">
                    <div>
                        <span className="filter-body-label">Operation: </span>
                        <select className="filter-body-select" value={data.operation}>
                            {operations}
                        </select>
                    </div>
                    <div>
                        <span className="filter-body-label">Value: </span>
                        {value_field}
                    </div>
                    <div className="filter-body-buttons">
                        <button className="filter-body-cancel">Delete</button>
                        <button className="filter-body-save">Save</button>
                    </div>
                </div>
            </div>
        )
    }
}

export class Creator extends Component {
  props: Props;
  state: {
    loading: boolean,
    redirectTo: ?string,
    user: ?APIUser,
    projects: ?APIProject[],
    expeditions: ?APIExpedition[],
    activeProject: ?APIProject,
    activeExpedition: ?APIExpedition
  }

  constructor(props: Props) {
    super(props);

    this.state = {
      loading: true,
      redirectTo: null,
      user: null,
      projects: [],
      expeditions: [],
      activeProject: null,
      activeExpedition: null
    }

    Promise.all([
      this.loadUser(),
      this.loadProjects(),
      this.loadExpeditions(this.projectSlug()),
      this.loadActiveProject(this.projectSlug()),
      this.loadActiveExpedition(this.projectSlug(), this.expeditionSlug()),
    ]).then(() => this.setState({ loading: false }));
  }

  componentWillReceiveProps(nextProps: Props) {
    const promises = [];
    const stateChange = {};

    const {
      projectSlug,
      expeditionSlug
    } = this.props.match.params;

    const {
      projectSlug: newProjectSlug,
      expeditionSlug: newExpeditionSlug
    } = nextProps.match.params;

    if (projectSlug != newProjectSlug) {
      promises.push(this.loadActiveProject(newProjectSlug));
      promises.push(this.loadExpeditions(newProjectSlug));
      stateChange.activeProject = null;
    } else if (!newProjectSlug) {
      promises.push(this.loadProjects());
      stateChange.activeProject = null;
    }

    if (expeditionSlug != newExpeditionSlug) {
      promises.push(this.loadActiveExpedition(newProjectSlug, newExpeditionSlug));
      promises.push(this.loadExpeditions(newProjectSlug));
      stateChange.activeExpedition = null;
    } else if (!newExpeditionSlug) {
      stateChange.activeExpedition = null;
    }

    if (promises.length > 0 || Object.keys(stateChange).length > 0) {
      this.setState({ loading: true, ...stateChange });
      Promise.all(promises).then(() => { this.setState({ loading: false }); });
    }

  }

  projectSlug(): ?string {
    return this.props.match.params.projectSlug;
  }

  expeditionSlug(): ?string {
    return this.props.match.params.expeditionSlug;
  }

  async onUserUpdate (user: APIUser) {
    log.debug('main -> onUserUpdate');
    this.setState({user: user});
  }

  async loadUser() {
    log.debug('main -> loadUser');
    const userRes = await FKApiClient.get().getCurrentUser();
    if (userRes.type === 'ok') {
      this.setState({ user: userRes.payload });
    } else {
      this.setState({ user: null });
    }
  }

  async loadProjects() {
    log.debug('main -> loadProjects');
    const projectRes = await FKApiClient.get().getCurrentUserProjects();
    if (projectRes.type === 'ok') {
      this.setState({ projects: projectRes.payload.projects });
    } else {
      this.setState({ projects: null });
    }
  }

  async loadExpeditions(projectSlug: ?string) {
    log.debug('main -> loadExpeditions', projectSlug);
    if (projectSlug) {
      const expRes = await FKApiClient.get().getExpeditionsByProjectSlug(projectSlug);
      if (expRes.type === 'ok') {
        this.setState({ expeditions: expRes.payload.expeditions });
      }
    } else {
      this.setState({ expeditions: null });
    }
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

  async loadActiveExpedition(projectSlug: ?string, expeditionSlug: ?string) {
    log.debug('main -> loadActiveExpedition', projectSlug, expeditionSlug);
    if (projectSlug && expeditionSlug) {
      const expRes = await FKApiClient.get().getExpeditionBySlugs(projectSlug, expeditionSlug);
      if (expRes.type === 'ok') {
        this.setState({ activeExpedition: expRes.payload });
      }
    } else {
      this.setState({ activeExpedition: null });
    }
  }

  onProjectCreate() {
    this.loadProjects();
    this.props.history.push('/');
  }

  onExpeditionCreate() {
    const projectSlug = this.projectSlug();
    if (projectSlug) {
      this.loadExpeditions(projectSlug);
      this.props.history.push(`/projects/${projectSlug}`);
    }
  }

  async onProjectUpdate(newSlug: ?string = null) {
    log.debug('onProjectUpdate', newSlug);
    if (newSlug) {
      this.props.history.replace(`/projects/${newSlug}/settings`);
    } else {
      await Promise.all([
        this.loadProjects(),
        this.loadActiveProject(this.projectSlug())
      ]);
    }
  }

  async onExpeditionUpdate(newSlug: ?string = null) {
    if (newSlug) {
      const projectSlug = this.projectSlug();
      if (projectSlug) {
        this.props.history.replace(`/projects/${projectSlug}/expeditions/${newSlug}`);
      } else {
        this.props.history.replace('/');
      }
    } else {
      await Promise.all([
        this.loadExpeditions(this.projectSlug()),
        this.loadActiveExpedition(this.projectSlug(), this.expeditionSlug())
      ]);
    }
  }

  handleLinkClick() {
    this.refs.dropdown.hide();
  }

  render() {
    const {
      user,
      activeProject,
      activeExpedition
    } = this.state;

    let {
      projects,
      expeditions
    } = this.state;

    if (projects && activeProject) {
      projects = projects.filter(p => p.id !== activeProject.id);
    }

    if (expeditions && activeExpedition) {
      expeditions = expeditions.filter(e => e.id !== activeExpedition.id);
    }

    const test_string_filter = {
        id: 0,
        attribute: "username",
        operation: "matches",
        query: "bar"
    }

    return (
        <div>
            <StringFilterComponent data={test_string_filter} options={[]}/>
            <StringFilterComponent data={test_string_filter} options={[["foo","bar"]]}/>
        </div>
    )
  }
}
