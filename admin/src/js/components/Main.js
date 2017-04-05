// @flow

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

export class Main extends Component {
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
        this.props.history.replace(`/projects/${newSlug}/expeditions/${newSlug}`);
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

    return (
      <div className="main">
        <div className="left">
          <div className="logo-area">
            <Link to="/"><img src={fieldkitLogo} alt="fieldkit logo" /></Link>
          </div>
          <div className="sidebar">
            <div className="sidebar-section project-section">
              <h5>Projects</h5>
              <Link to={`/`}>
                <div className="bt-icon medium">
                  <HamburgerIcon />
                </div>
                All
              </Link>
              { activeProject &&
                <div>
                  <p className="project-name">
                    <div className="bt-icon medium">
                      <ArrowDownIcon />
                    </div>
                    <span>{activeProject.name}</span>
                    <a className="bt-icon small" href={`//${activeProject.slug}.${fkHost()}/`} alt="go to project`" target="_blank">
                      <OpenInNewIcon />
                    </a>
                  </p>
                  <div className="sidebar-nav">
                    <NavLink exact to={`/projects/${activeProject.slug}`}>Expeditions</NavLink>
                    <NavLink to={`/projects/${activeProject.slug}/settings`}>Settings</NavLink>
                  </div>
                </div> }
            </div>
            {activeProject && activeExpedition &&
              <div className="sidebar-section expedition-section">
                <h5>Expedition</h5>
                  <p className="expedition-name">
                    <span>{activeExpedition.name}</span>
                    <a className="bt-icon medium" href={`//${activeProject.slug}.${fkHost()}/${activeExpedition.slug}`}
                      alt="go to expedition"
                      target="_blank">
                      <OpenInNewIcon />
                    </a>
                  </p>
                  <div className="sidebar-nav">
                    <NavLink to={`/projects/${activeProject.slug}/expeditions/${activeExpedition.slug}/datasources`}>Data Sources</NavLink>
                    <NavLink to={`/projects/${activeProject.slug}/expeditions/${activeExpedition.slug}/teams`}>Teams</NavLink>
                    <NavLink exact to={`/projects/${activeProject.slug}/expeditions/${activeExpedition.slug}`}>Settings</NavLink>
                    {/* <NavLink to={`/projects/${activeProject.slug}/expeditions/${activeExpedition.slug}/website`}>Website</NavLink> */}
                  </div>
              </div>
              }
          </div>


          <footer>
            <Link to="/help">Help</Link> {}- <Link to="/contact">Contact Us</Link> - <Link to="/privacy">Privacy Policy</Link>
          </footer>
        </div>

        <div className="right">
          <Dropdown className="account-dropdown" ref="dropdown">
            <DropdownTrigger className="trigger">
              <div className="user-avatar small">
                <img src={placeholderImage} alt="profile" />
              </div>
            </DropdownTrigger>
            <DropdownContent className="dropdown-contents">
              <div className="header">
                Signed in as <strong>{user ? user.username : ''}</strong>
              </div>
              <div className="nav">
                <Link to="/profile" onClick={this.handleLinkClick.bind(this)}>Profile</Link>
                <Link to="/logout">Logout</Link>
              </div>
            </DropdownContent>
          </Dropdown>

          <div className="contents">
            <Switch>
              <RouteOrLoading
                path="/projects/:projectSlug/expeditions/:expeditionSlug/teams"
                component={Teams}
                required={[activeProject, activeExpedition]}
                project={activeProject}
                expedition={activeExpedition} />
              <RouteOrLoading
                path="/projects/:projectSlug/expeditions/:expeditionSlug/datasources"
                component={DataSources}
                required={[activeProject, activeExpedition]}
                project={activeProject}
                expedition={activeExpedition} />
              <RouteOrLoading
                path="/projects/:projectSlug/expeditions/:expeditionSlug"
                component={Expedition}
                required={[activeProject, activeExpedition]}
                project={activeProject}
                expedition={activeExpedition}
                onUpdate={this.onExpeditionUpdate.bind(this)} />
              <RouteOrLoading
                path="/projects/:projectSlug/settings"
                component={ProjectSettings}
                required={[activeProject]}
                project={activeProject}
                user={user}
                onUpdate={this.onProjectUpdate.bind(this)} />
              <RouteOrLoading
                path="/projects/:projectSlug"
                component={ProjectExpeditions}
                required={[activeProject, expeditions]}
                project={activeProject}
                expeditions={expeditions}
                onUpdate={this.onProjectUpdate.bind(this)}
                onExpeditionCreate={this.onExpeditionCreate.bind(this)} />
              <RouteOrLoading
                path="/profile"
                component={Profile}
                required={[user]}
                user={user}
                onUpdate={this.onUserUpdate.bind(this)} />
              <RouteOrLoading
                path="/"
                component={Projects}
                required={[projects]}
                projects={projects}
                onProjectCreate={this.onProjectCreate.bind(this)} />
            </Switch>
          </div>
        </div>
      </div>
    )
  }
}
