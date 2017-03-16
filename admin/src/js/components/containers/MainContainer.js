// @flow

import React, { Component } from 'react'
import { Link, NavLink } from 'react-router-dom';

import fieldkitLogo from '../../../img/logos/fieldkit-logo-red.svg';
import placeholderImage from '../../../img/profile_placeholder.svg'
import '../../../css/main.css'

type Breadcrumb = {
  url: string;
  text: string;
}

export class MainContainer extends Component {
  props: {
    breadcrumbs?: Breadcrumb[];
    projectSlug?: string;
    expeditionName?: string;
    expeditionSlug?: string;
    userPhoto?: string;
    children?: any;
  }

  render() {
    let {
      projectSlug,
      expeditionName,
      expeditionSlug,
      // userPhoto
    } = this.props;

    const hasExpedition = expeditionName && expeditionSlug && projectSlug;

    // Appease flow
    projectSlug = projectSlug || '';
    expeditionName = expeditionName || '';
    expeditionSlug = expeditionSlug || '';

    return (
      <div className="main">
        <div className="left">
          <div className="logo-area">
            <img src={fieldkitLogo} alt="fieldkit logo" />
          </div>

          { hasExpedition &&
            <div className="expedition-sidebar">
              <div className="expedition-name">
                <span>{expeditionName}</span>
                {/* TODO: use image icon */}
                <Link to={`https://${projectSlug}.fieldkit.org/${expeditionSlug}`}>GO</Link>
              </div>
              <div className="settings">
                {/* TODO: use image icon */}
                <Link to={`/projects/${projectSlug}/expeditions/${expeditionSlug}`}>Settings</Link>
              </div>
              <div className="nav">
                <NavLink to={`/projects/${projectSlug}/expeditions/${expeditionSlug}/datasources`}>Data Sources</NavLink>
                <NavLink to={`/projects/${projectSlug}/expeditions/${expeditionSlug}/teams`}>Teams</NavLink>
                <NavLink to={`/projects/${projectSlug}/expeditions/${expeditionSlug}/website`}>Website</NavLink>
              </div>
            </div> }

          <footer>
            <Link to="/help">Help</Link> {}- <Link to="/contact">Contact Us</Link> - <Link to="/privacy">Privacy Policy</Link>
          </footer>
        </div>

        <div className="right">
          <div className="nav">
            <div className="breadcrumbs">
              { this.props.breadcrumbs &&
                this.props.breadcrumbs.map((b, i) => <span> / <Link to={b.url}>{b.text}</Link></span>) }
            </div>
            <div className="profile-image">
              <img src={placeholderImage} alt="profile" />
            </div>
          </div>

          <div className="contents">
            { this.props.children }
          </div>
        </div>
      </div>
    )
  }
}
