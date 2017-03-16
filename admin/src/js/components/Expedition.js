// @flow weak

import React, { Component } from 'react'
import { Route, Link } from 'react-router-dom'

import { MainContainer } from './containers/MainContainer';
import { ProjectExpeditionForm } from './forms/ProjectExpeditionForm';
import { FKApiClient } from '../api/api';

type Props = {
	match: Object;
	location: Object;
	history: Object;
}

export class Expedition extends Component {
	
	props: Props;
	state: {
		project: ?Object,
		expedition: Object
	}

	constructor(props: Props) {
		super(props);

		this.state = {
			project: null,
			expedition: null
		};

		this.loadData();
	}

	expeditionSlug() {
		return this.props.match.params.expeditionSlug;
		this.props.match
	}

	projectSlug() {
		return this.props.match.params.projectSlug;
	}

	async loadData() {
		const projectSlug = this.projectSlug();
		const expeditionSlug = this.expeditionSlug();

		const projectRes = await FKApiClient.get().getProjectBySlug(projectSlug);
		if (projectRes.type === 'ok') {
			this.setState({ project: projectRes.payload })
		}

		const expeditionsRes = await FKApiClient.get().getExpeditionsByProjectSlug(projectSlug);
		if (expeditionsRes.type === 'ok') {
			const expedition = expeditionsRes.payload.find(function(obj){
				return obj.slug === expeditionSlug;
			});
			this.setState({ expedition: expedition })
		}
	}

	async onExpeditionSave(name: string, description: string) {
		// TODO: this isn't implemented on the backend yet!
	}	

	render() {
		const { project } = this.state;
		const { expedition } = this.state;
		const projectSlug = this.projectSlug();
		const expeditionSlug = this.expeditionSlug();

		return (
			<MainContainer
				breadcrumbs={[	{ url: '/', text: 'Projects'},
								{ url: `/projects/${projectSlug}`, text: project ? project.name : 'Current Project' },
								{ url: `/projects/${projectSlug}/expeditions/${expeditionSlug}`, text: expedition ? expedition.name : 'Current Expedition' }
								]}
			>
				<div className="expedition">
					<h1>Expedition Settings</h1>
					<ProjectExpeditionForm 
						name={expedition ? expedition.name : undefined}
						projectSlug={projectSlug}
						description={expedition ? expedition.description : undefined}
						onSave={this.onExpeditionSave.bind(this)}
					/>
				</div>
			</MainContainer>
		)
	}
}