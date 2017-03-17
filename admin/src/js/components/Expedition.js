// @flow weak

import React, { Component } from 'react'
import { Route, Link } from 'react-router-dom'

import { ProjectExpeditionForm } from './forms/ProjectExpeditionForm';
import { FKApiClient } from '../api/api';

import type { APIProject } from '../api/types';
import type { APIExpedition } from '../api/types';

import '../../css/home.css'

type Props = {
	project: APIProject;
	expedition: APIExpedition;
	onExpeditionUpdate: (newSlug: ?string) => void;

	match: Object;
	location: Object;
	history: Object;
}

export class Expedition extends Component {
	
	props: Props;
	state: {
		expedition: Object
	}

	constructor(props: Props) {
		super(props);

		this.state = {
			project: null,
			expeditions: []
		};
	}

	async onExpeditionSave(name: string, description: string) {
		// TODO	
	}	

	render() {
		const { project } = this.props;
		const projectSlug = project.slug;
		const { expedition } = this.props;
		const expeditionSlug = project.slug;

		return (

			<div className="expedition">
				<h1>Expedition Settings</h1>
				<ProjectExpeditionForm
					projectSlug={projectSlug}
					onSave={this.onExpeditionSave.bind(this)} />				
			</div>
		)
	}
}