// @flow

import React, { Component } from 'react'
import { Switch, Route, Link, NavLink, Redirect } from 'react-router-dom';

import { FKApiClient } from '../api/api';
import type { APIUser, APIProject, APIExpedition } from '../api/types';

import { Projects } from './pages/Projects';
import { Project } from './pages/Project';
import { Expedition } from './pages/Expedition'
import { Teams } from './pages/Teams';

import fieldkitLogo from '../../img/logos/fieldkit-logo-red.svg';
import placeholderImage from '../../img/profile_placeholder.svg'
import '../../css/main.css'

type Props = {
  match: Object;
  location: Object;
  history: Object;
}

export class Main extends Component {
  props: Props;
  state: {
    loading: boolean,
    redirectTo: ?string,
    user: ?APIUser,
    project: ?APIProject,
    expedition: ?APIExpedition
  }

  constructor(props: Props) {
    super(props);

    this.state = {
      loading: true,
      redirectTo: null,
      user: null,
      project: null,
      expedition: null
    }

    const {
      projectSlug,
      expeditionSlug
    } = props.match.params;

    Promise.all([
      this.loadUser(),
      this.loadProject(projectSlug),
      this.loadExpedition(projectSlug, expeditionSlug)
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
      promises.push(this.loadProject(newProjectSlug));
      stateChange.project = null;
    } else if (!newProjectSlug) {
      stateChange.project = null;
    }

    if (expeditionSlug != newExpeditionSlug) {
      promises.push(this.loadExpedition(newProjectSlug, newExpeditionSlug));
      stateChange.expedition = null;
    } else if (!newExpeditionSlug) {
      stateChange.expedition = null;
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
    const userRes = await FKApiClient.get().getCurrentUser();
    if (userRes.type == 'ok' && userRes.payload) {
      this.setState({ user: userRes.payload });
    }
  }

  async loadProject(projectSlug: ?string = this.projectSlug()) {
    if (projectSlug) {
      const projectRes = await FKApiClient.get().getProjectBySlug(projectSlug);
      if (projectRes.type == 'ok' && projectRes.payload) {
        this.setState({ project: projectRes.payload });
      }
    } else {
      this.setState({ project: null });
    }
  }

  async loadExpedition(projectSlug: ?string = this.projectSlug(), expeditionSlug: ?string = this.expeditionSlug()) {
    if (projectSlug && expeditionSlug) {
      const expRes = await FKApiClient.get().getExpeditionBySlugs(projectSlug, expeditionSlug);
      if (expRes.type == 'ok' && expRes.payload) {
        this.setState({ expedition: expRes.payload });
      }
    } else {
      this.setState({ expedition: null });
    }
  }

  onProjectUpdate(newSlug: ?string = null) {
    if (newSlug) {
      this.setState({ redirectTo: `/projects/${newSlug}`})
    } else {
      this.setState({ loading: true })
      this.loadProject(this.projectSlug()).then(() => this.setState({ loading: false }));
    }
  }

  onExpeditionUpdate(newSlug: ?string = null) {
    if (newSlug) {
      const projectSlug = this.projectSlug();
      if (projectSlug) {
        this.setState({ redirectTo: `/projects/${projectSlug}/expeditions/${newSlug}`})
      } else {
        this.setState({ redirectTo: '/' });
      }
    } else {
      this.setState({ loading: true })
      this.loadExpedition().then(() => this.setState({ loading: false }));
    }
  }

  render() {
    if (this.state.redirectTo) {
      return <Redirect to={this.state.redirectTo} />;
    }

    const {
      project,
      expedition
    } = this.state;

    const breadcrumbs = [];
    if (project && this.projectSlug()) {
      breadcrumbs.push(
        <Link key={0} to={`/projects/${project.slug}`}>{project.name}</Link>
      );
      if (expedition && this.expeditionSlug()) {
        breadcrumbs.push(
          <Link key={1} to={`/projects/${project.slug}/expeditions/${expedition.slug}`}>{expedition.name}</Link>
        );
      }
    }

    return (
      <div className="main">
        <div className="left">
          <div className="logo-area">
            <img src={fieldkitLogo} alt="fieldkit logo" />
          </div>

          { project && expedition &&
            <div className="expedition-sidebar">
              <div className="expedition-name">
                <span>{expedition.name}</span>
                {/* TODO: use image icon */}
                <a href={`https://${project.slug}.fieldkit.org/${expedition.slug}`}
                  alt="go to expedition"
                  target="_blank">GO</a>
              </div>
              <div className="settings">
                {/* TODO: use image icon */}
                <Link to={`/projects/${project.slug}/expeditions/${expedition.slug}`}>Settings</Link>
              </div>
              <div className="nav">
                <NavLink to={`/projects/${project.slug}/expeditions/${expedition.slug}/datasources`}>Data Sources</NavLink>
                <NavLink to={`/projects/${project.slug}/expeditions/${expedition.slug}/teams`}>Teams</NavLink>
                <NavLink to={`/projects/${project.slug}/expeditions/${expedition.slug}/website`}>Website</NavLink>
              </div>
            </div> }

          <footer>
            <Link to="/help">Help</Link> {}- <Link to="/contact">Contact Us</Link> - <Link to="/privacy">Privacy Policy</Link>
          </footer>
        </div>

        <div className="right">
          <div className="nav">
            <div className="breadcrumbs">
              { breadcrumbs.length > 0 && breadcrumbs.reduce((prev, curr) => [prev, ' / ', curr]) }
            </div>
            <div className="profile-image">
              <img src={placeholderImage} alt="profile" />
            </div>
          </div>

          <div className="contents">
            <Switch>
              <Route path="/projects/:projectSlug/expeditions/:expeditionSlug/teams" render={props => {
                if (project && expedition) {
                  return (
                    <Teams
                      project={project}
                      expedition={expedition}
                      {...props} />
                  )
                } else {
                  return <div></div>
                }
              }} />
              <Route path="/projects/:projectSlug/expeditions/:expeditionSlug" render={props => {
                if (project && expedition) {
                  return (
                    <Expedition
                      project={project}
                      expedition={expedition}
                      {...props} />
                  )
                } else {
                  return <div></div>
                }
              }} />
              <Route path="/projects/:projectSlug" render={props => {
                if (project) {
                  return (
                    <Project
                      project={project}
                      onProjectUpdate={this.onProjectUpdate.bind(this)}
                      {...props} />
                  )
                } else {
                  return <div></div>
                }
              }} />
              <Route path="/" render={props => <Projects {...props} />} />
            </Switch>
          </div>
        </div>
      </div>
    )
  }
}
