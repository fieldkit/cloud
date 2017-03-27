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

import { RouteOrLoading } from './shared/RequiredRoute';
import { Projects } from './pages/Projects';
import { Project } from './pages/Project';
import { Expedition } from './pages/Expedition'
import { Teams } from './pages/Teams';
import { DataSources } from './pages/DataSources';
import { Profile } from './pages/Profile';

import fieldkitLogo from '../../img/logos/fieldkit-logo-red.svg';
import placeholderImage from '../../img/profile_placeholder.svg'
import externalLinkImg from '../../img/icons/icon-external-link.png'
import gearImg from '../../img/icons/icon-gear.png'
import '../../css/main.css'

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
      promises.push(this.loadProjects());
      stateChange.activeProject = null;
    } else if (!newProjectSlug) {
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
    if (newSlug) {
      this.setState({ redirectTo: `/projects/${newSlug}`})
    } else {
      this.setState({ loading: true })
      await Promise.all([
        this.loadProjects(),
        this.loadActiveProject(this.projectSlug())
      ]);
      this.setState({ loading: false });
    }
  }

  async onExpeditionUpdate(newSlug: ?string = null) {
    if (newSlug) {
      const projectSlug = this.projectSlug();
      if (projectSlug) {
        this.setState({ redirectTo: `/projects/${projectSlug}/expeditions/${newSlug}`})
      } else {
        this.setState({ redirectTo: '/' });
      }
    } else {
      this.setState({ loading: true });
      await Promise.all([
        this.loadExpeditions(this.projectSlug()),
        this.loadActiveExpedition(this.projectSlug(), this.expeditionSlug())
      ]);
      this.setState({ loading: false });
    }
  }

  handleLinkClick() {
    this.refs.dropdown.hide();
  }

  render() {
    if (this.state.redirectTo) {
      return <Redirect to={this.state.redirectTo} />;
    }

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

          { activeProject &&
            <div className="sidebar">
              <div className="sidebar-section project-nav">
                <div className="name project-name">
                  <Link to={`/projects/${activeProject.slug}`}>{activeProject.name}</Link>
                  <a className="out-link" href={`https://${activeProject.slug}.fieldkit.org/`} alt="go to project`" target="_blank">
                    <img src={externalLinkImg} alt="external link" />
                  </a>
                </div>
                <div className="settings">
                  <Link to={`/projects/${activeProject.slug}/settings`}>
                    <img src={gearImg} alt="settings gear" />
                  </Link>
                </div>
                { !activeExpedition &&
                  <div className="sidebar-nav">
                    <NavLink to={`/projects/${activeProject.slug}`}>Expeditions</NavLink>
                  </div> }
              </div>
              { activeExpedition &&
                <div className="sidebar-section expedition-nav">
                  <div className="name expedition-name">
                    <Link to={`/projects/${activeProject.slug}/expeditions/${activeExpedition.slug}`}>{activeExpedition.name}</Link>
                    <a className="out-link" href={`https://${activeProject.slug}.fieldkit.org/${activeExpedition.slug}`}
                      alt="go to expedition"
                      target="_blank">
                      <img src={externalLinkImg} alt="external link" />
                    </a>
                  </div>
                  <div className="settings">
                    <Link to={`/projects/${activeProject.slug}/expeditions/${activeExpedition.slug}`}>
                      <img src={gearImg} alt="settings gear"/>
                    </Link>
                  </div>
                  <div className="sidebar-nav">
                    <NavLink to={`/projects/${activeProject.slug}/expeditions/${activeExpedition.slug}/datasources`}>Data Sources</NavLink>
                    <NavLink to={`/projects/${activeProject.slug}/expeditions/${activeExpedition.slug}/teams`}>Teams</NavLink>
                    {/* <NavLink to={`/projects/${activeProject.slug}/expeditions/${activeExpedition.slug}/website`}>Website</NavLink> */}
                  </div>
              </div> }
              { activeProject && !activeExpedition && projects &&
                <div className="project-select">
                  <ul>
                    { projects.map((p, i) => (
                      <li key={i}><Link to={`/projects/${p.slug}`}>{p.name}</Link></li>
                    ))}
                  </ul>
                </div> }
              { activeProject && activeExpedition && expeditions &&
                <div className="expedition-select">
                  <ul>
                    { expeditions.map((e, i) => (
                      <li key={i}><Link to={`/projects/${activeProject.slug}/expeditions/${e.slug}`}>{e.name}</Link></li>
                    ))}
                  </ul>
                </div> }
            </div> }

          <footer>
            <Link to="/help">Help</Link> {}- <Link to="/contact">Contact Us</Link> - <Link to="/privacy">Privacy Policy</Link>
          </footer>
        </div>

        <div className="right">
          <div className="nav">
            <Dropdown className="account-dropdown" ref="dropdown">
              <DropdownTrigger className="trigger">
                <img src={placeholderImage} alt="profile" />
              </DropdownTrigger>
              <DropdownContent className="dropdown-contents">
                <div className="header">
                  Signed in as <strong>{user ? user.username : ''}</strong>
                </div>
                <ul>
                  <li><Link to="/profile" onClick={this.handleLinkClick.bind(this)}>Profile</Link></li>
                  <li><Link to="/logout">Logout</Link></li>
                </ul>
              </DropdownContent>
            </Dropdown>
          </div>

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
                path="/projects/:projectSlug"
                component={Project}
                required={[activeProject, expeditions]}
                project={activeProject}
                expeditions={expeditions}
                onUpdate={this.onProjectUpdate.bind(this)}
                onExpeditionCreate={this.onExpeditionCreate.bind(this)} />
              <RouteOrLoading
                path="/profile"
                component={Profile}
                required={[user]}
                user={user} />
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
