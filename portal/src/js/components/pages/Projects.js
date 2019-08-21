// @flow

import React, { Component } from 'react'
import { Route, Link } from 'react-router-dom'
import ReactModal from 'react-modal';

import type { Match as RouterMatch, Location as RouterLocation, RouterHistory } from 'react-router-dom';

import { ProjectForm } from '../forms/ProjectForm';
import { FKApiClient } from '../../api/api';
import { joinPath } from '../../common/util';
import type { APIProject, APINewProject } from '../../api/types';

type Props = {
projects: APIProject[];
onProjectCreate: () => void;

match: RouterMatch;
location: RouterLocation;
history: RouterHistory;
}

export class Projects extends Component {
    props: $Exact<Props>;

    async onProjectCreate(p: APINewProject) {
        const {match} = this.props;
        const project = await FKApiClient.get().createProject(p);
        if (project.type === 'ok') {
            this.props.onProjectCreate();
            this.props.history.push(joinPath(match.url, 'projects', p.slug));
        } else {
            return project.errors;
        }
    }

    render() {
        const {projects, match} = this.props;

        return (
            <div className="projects">
                <Route path={ joinPath(match.url, 'new-project') } render={ () => <ReactModal isOpen={ true } contentLabel="New project form" className="modal" overlayClassName="modal-overlay">
                                                                                      <h2>Create New Project</h2>
                                                                                      <ProjectForm onCancel={ () => this.props.history.push("/") } onSave={ this.onProjectCreate.bind(this) } />
                                                                                  </ReactModal> } />
                <h1>Projects</h1>
                <div id="projects" className="gallery-list projects">
                    { projects.map((p, i) => <Link to={ joinPath(match.url, 'projects', p.slug) } key={ `project-${i}` } className="gallery-list-item projects" style={ { backgroundImage: 'url(' + FKApiClient.get().projectPictureUrl(p.id) + ')' } }>
                                                 <h4>{ p.name }</h4>
                                             </Link>) }
                    { projects.length === 0 &&
                      <span className="empty">You don't have any projects yet.</span> }
                </div>
                <Link className="button" to={ joinPath(match.url, 'new-project') }>Create new project</Link>
            </div>
        )
    }
}
