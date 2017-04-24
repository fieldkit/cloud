/* @flow */

import React, { Component } from 'react'
import { Switch, Route, Link, NavLink, Redirect } from 'react-router-dom';
import Dropdown, { DropdownTrigger, DropdownContent } from 'react-simple-dropdown';
import { AddIcon } from './icons/Icons'
import type { APIErrors } from '../api/types';

import type {
  Match as RouterMatch,
  Location as RouterLocation,
  RouterHistory
} from 'react-router-dom';

import type {Attr, Collection, Filter, FilterFn, StringFilter, DateFilter, NumFilter} from './Collection';
import {cloneCollection, emptyStringFilter, emptyNumFilter, emptyDateFilter, stringifyOptions} from './Collection';

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

import { StringFilterComponent, NumFilterComponent, DateFilterComponent } from './Filters'

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

function emptyCollection() : Collection {
    return {
        name: "",
        id: "",
        filters: [],
        guid_filters: [],
        string_filters: [],
        num_filters: [],
        geo_filters: [],
        date_filters: [],
        mods: [],
        string_mods: [],
        num_mods: [],
        geo_mods: []
    }
}


export class Creator extends Component {
  attributes: { [string]: Attr };
  props: Props;
  state: {
    loading: boolean,
    redirectTo: ?string,
    user: ?APIUser,
    projects: ?APIProject[],
    expeditions: ?APIExpedition[],
    activeProject: ?APIProject,
    activeExpedition: ?APIExpedition,
    collection: Collection,
    errors: ?APIErrors
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
      loading: true,
      redirectTo: null,
      user: null,
      projects: [],
      expeditions: [],
      activeProject: null,
      activeExpedition: null,
      collection: emptyCollection(),
      errors: null
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

  handleLinkClick() {
    this.refs.dropdown.hide();
  }

  addFilter(attr: string) {
    let collection = cloneCollection(this.state.collection);
    const attribute = this.attributes[attr];
    let new_filter;
    if(attribute.type === "string"){
        new_filter = emptyStringFilter(collection, attribute);
        collection.string_filters.push(new_filter)
    } else if (attribute.type === "num"){
        new_filter = emptyNumFilter(collection, attribute);
        collection.num_filters.push(new_filter)
    } else if (attribute.type === "date"){
        new_filter = emptyDateFilter(collection, attribute);
        collection.date_filters.push(new_filter)
    }
    if(new_filter){
        collection.filters.push(new_filter.id)
        this.setState({collection})
    }
  }

  updateFilter:FilterFn = (filter, update) => {
      let new_filter = Object.assign(filter,update)
      let collection = cloneCollection(this.state.collection)
      if(new_filter.type === "string"){
        collection.string_filters = collection.string_filters.filter(f => f.id !== new_filter.id)
        collection.string_filters.push(new_filter)
      } else if (filter.type === "date") {
        collection.date_filters = collection.date_filters.filter(f => f.id !== new_filter.id)
        collection.date_filters.push(new_filter)
      } else if (new_filter.type === "num") {
        collection.num_filters = collection.num_filters.filter(f => f.id !== new_filter.id)
        collection.num_filters.push(new_filter)
      }
      this.setState({collection})
  }

  deleteFilter(f: Filter){
      let collection = cloneCollection(this.state.collection)
      if(f.type === "string"){
        collection.string_filters = collection.string_filters.filter(cf => cf.id !== f.id)
      } else if(f.type === "num"){
        collection.num_filters = collection.num_filters.filter(cf => cf.id !== f.id)
      } else if(f.type === "date"){
        collection.date_filters = collection.date_filters.filter(cf => cf.id !== f.id)
      }
      collection.filters = collection.filters.filter(cf => cf !== f.id) 
      this.setState({collection})
  }

  save(){
    const {collection} = this.state;

    console.log(JSON.stringify(collection,null,' '))
  }

  render() {
    const {
      user,
      activeProject,
      activeExpedition,
      errors
    } = this.state;

    let {
      projects,
      expeditions,
      collection
    } = this.state;

    let component_lookup : { [number]: Object } = {};

    collection.num_filters.forEach((f,i) => {
        component_lookup[f.id] = <NumFilterComponent creator={this} data={f} key={"n-" + i} errors={errors}/>
    })
    collection.string_filters.forEach((f,i) => {
        component_lookup[f.id] = <StringFilterComponent creator={this} data={f} key={"s-" + i} errors={errors}/>
    })
    collection.date_filters.forEach((f,i) => {
        component_lookup[f.id] = <DateFilterComponent creator={this} data={f} key={"d-" + i} errors={errors}/>
    })
    
   let filters_by_attribute = Object.keys(this.attributes).reduce((m,attr_name) => {
     const attr = this.attributes[attr_name]
     const {type,options} = attr
     let options_block = null
     if(options.length > 0){
       options_block = <span className="filter-row-options">{stringifyOptions(attr)}</span> 
     }
     const attribute_row = (
         <div key={"filter-" + attr_name} className="accordion-row-header">
            <h4>{attr_name}</h4>            
            <button className="secondary" onClick={() => this.addFilter(attr_name)}>
              <div className="bt-icon medium">
                <AddIcon />
              </div>
              Add Filter
            </button>
         </div>
     )

     let buttons: Object[] = [attribute_row]
     collection.filters.forEach((f) => {
        const component = component_lookup[f]
        if(component){
            const {attribute} = component.props.data
            if(attr_name === attribute){
                buttons.push(component)
            } 
        }
     })
     m[attr_name] = buttons
     return m
   },{})
     
    const components = Object.entries(filters_by_attribute).map(([attr,components]) => {
        return (
          <div className="accordion-row expanded" key={attr}>
              {components}
          </div>
        )
    })
    
    return (
        <div className="row">
          <div className="accordion">
            {components}
          </div>
          <button onClick={() => this.save()}>Save Filters</button>
        </div>
    )
  }
}
