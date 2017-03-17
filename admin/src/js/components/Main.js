// @flow

import React, { Component } from 'react'
import { Switch, Route, Link, NavLink, Redirect } from 'react-router-dom';

import { FKApiClient } from '../api/api';

import { Projects } from './Projects';
import { Project } from './Project';

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
    user: ?Object,
    project: ?Object,
    expedition: ?Object
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

    Promise.all([
      this.loadUser(),
      this.loadProject(),
      this.loadExpedition()
    ]).then(() => this.setState({ loading: false }));
  }

  projectSlug(): ?string {
    return this.props.match.params.projectSlug;
  }

  expeditionSlug(): ?string {
    return this.props.match.params.expeditionSlug;
  }

  async loadUser() {
    const userRes = await FKApiClient.get().getUser();
    if (userRes.type == 'ok' && userRes.payload) {
      this.setState({ user: userRes.payload.user });
    }
  }

  async loadProject() {
    const projectSlug = this.projectSlug();
    if (projectSlug) {
      const projectRes = await FKApiClient.get().getProjectBySlug(projectSlug);
      console.log(projectRes);
      if (projectRes.type == 'ok' && projectRes.payload) {
        this.setState({ project: projectRes.payload });
      }
    }
  }

  async loadExpedition() {
    const projectSlug = this.projectSlug();
    const expeditionSlug = this.expeditionSlug();
    if (projectSlug && expeditionSlug) {
      const expRes = await FKApiClient.get().getExpeditionBySlugs(projectSlug, expeditionSlug);
      if (expRes.type == 'ok' && expRes.payload) {
        this.setState({ expedition: expRes.payload });
      }
    }
  }

  onProjectUpdate(newSlug: ?string = null) {
    if (newSlug) {
      this.setState({ redirectTo: `/projects/${newSlug}`})
    } else {
      this.loadProject();
    }
  }

  render() {
    if (this.state.redirectTo) {
      return <Redirect to={this.state.redirectTo} />;
    }

    const projectSlug = this.projectSlug();
    const expeditionSlug = this.expeditionSlug();

    const {
      project,
      expedition
    } = this.state;

    const breadcrumbs = [];
    if (project) {
      breadcrumbs.push(
        <Link to={`/projects/${project.slug}`}>{project.name}</Link>
      );
      if (expedition) {
        breadcrumbs.push(
          <Link to={`/projects/${project.slug}/expeditions/${expedition.slug}`}>{expedition.name}</Link>
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
                <Link to={`https://${project.slug}.fieldkit.org/${expedition.slug}`}>GO</Link>
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
