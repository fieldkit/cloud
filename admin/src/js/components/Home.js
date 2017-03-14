// @flow weak

import React, { Component } from 'react'
import { Route, Link } from 'react-router-dom'
import ReactModal from 'react-modal';

import { ProjectForm } from './forms/ProjectForm';

import '../../css/home.css'

type Props = {
  newProjectModal?: boolean;
  match: Object;
  location: Object;
  history: Object;
}

export class Home extends Component {
  props: Props;

  onSave() {
    // TODO: really save
    this.props.history.push("/app")
  }

  render () {
    return (
      <div className="main">
        <Route path="/app/new-project" render={() =>
          <ReactModal
            isOpen={true}
            contentLabel="Minimal Modal Example">
            <ProjectForm
              onCancel={() => this.props.history.push("/app")}
              onSave={this.onSave.bind(this)} />
          </ReactModal> } />

        <Link to="/app/new-project">Show new project modal</Link>
      </div>
    )
  }
}
