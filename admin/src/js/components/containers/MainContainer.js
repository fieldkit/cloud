// @flow

import React, { Component } from 'react'

import fieldkitLogo from '../../../img/logos/fieldkit-logo-red.svg';
import '../../../css/main.css'

export class MainContainer extends Component {
  render() {
    return (
      <div className="main">
        <div className="left">
          <img src={fieldkitLogo} alt="fieldkit logo" />
        </div>
        <div className="nav">
          <div className="breadcrumbs">
            { this.props.breadcrumbs && this.props.breadcrumbs.map((b, i) =>
              <div></div> )}
          </div>
          <div className="profile-image">
          </div>
        </div>
        <div className="contents">
          { this.props.children }
        </div>
      </div>
    )
  }
}
