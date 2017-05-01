// @flow weak

import React, { Component } from 'react'
import { Route, Link } from 'react-router-dom'
import ReactModal from 'react-modal';

import { ProjectCollectionsForm } from '../forms/ProjectCollectionsForm';
import { FKApiClient } from '../../api/api';
import { FormContainer } from '../containers/FormContainer';
import Collapse, { Panel } from 'rc-collapse';
import 'rc-collapse/assets/index.css';

import type { APIProject, APIExpedition, APIUser, APITeam, APIMember, APICollection, APINewCollection } from '../../api/types';

type Props = {
  project: APIProject;
  expeditions: APIExpedition[],
  teams: APITeam[],
  members: APIMember[],
  users: APIUser[],

  match: Object;
  location: Object;
  history: Object;
}

export class ProjectCollections extends Component {

  props: Props;
  state: {
    collections: APICollection[],
    collectionDeletion: ?{
      contents: React$Element<*>;
      collectionId: number;
    }
  }

  constructor (props: Props) {
    super(props);
    this.state = {
      collections: [],
      collectionDeletion: null
    }

    this.loadCollections();
  }

  async loadCollections() {
    // TO-DO
    // const collectionRes = await FKApiClient.get().getCollectionsBySlugs(this.props.project.slug);
    // if (collectionsRes.type === 'ok' && collectionsRes.payload) {
    //   const { collections } = collectionsRes.payload;
    //   this.setState({ collections });
    // }
  }

  async onCollectionCreate(c: APINewCollection) {
    // TO-DO
    // const { project, match } = this.props;
    // const collectionRes = await FKApiClient.get().createCollection(project.id, c);
    // if (collectionRes.type === 'ok') {
    //   await this.loadCollections();
    //   this.props.history.push(`${match.url}`);
    // } else {
    //   return collectionRes.errors;
    // }
  }

  async onTeamUpdate(collectionId: number, c: APINewCollection) {
    // TO-DO
    // const collectionRes = await FKApiClient.get().updateCollection(collectionId, c);
    // if(collectionRes.type === 'ok' && collectionRes.payload) {
    //   await this.loadCollections();
    // } else {
    //   return collectionRes.errors;
    // }
  }  

  startCollectionDelete(c: APICollection) {
    const collectionId = c.id;
    this.setState({
      collectionDeletion: {
        contents: <span>Are you sure you want to delete the <strong>{c.name}</strong> collection?</span>,
        collectionId
      }
    })
  }

  async confirmCollectionDelete() {
    const { collectionDeletion } = this.state;

    if (collectionDeletion) {
      const { collectionId } = collectionDeletion;

      const collectionRes = await FKApiClient.get().deleteCollection(collectionId);
      if (collectionRes.type === 'ok') {
        await this.loadCollections();
        this.setState({ collectionDeletion: null })
      } else {
        // TODO: show errors somewhere
      }
    }
  }
  
  cancelCollectionDelete() {
    this.setState({ collectionDeletion: null });
  }

  render() {
    const { match } = this.props;
    const { teams, members, users, memberDeletion, teamDeletion } = this.state;

    return (
      <div className="teams">

        <Route path={`${match.url}/new-team`} render={() =>
          <ReactModal isOpen={true} contentLabel="Create New Team" className="modal" overlayClassName="modal-overlay">
            <h2>Create New Team</h2>
            <TeamForm
              onCancel={() => this.props.history.push(`${match.url}`)}
              onSave={this.onTeamCreate.bind(this)} />
          </ReactModal> } />

        <Route path={`${match.url}/:teamId/add-member`} render={props =>
          <ReactModal isOpen={true} contentLabel="Add Members" className="modal" overlayClassName="modal-overlay">
            <h2>Add Member</h2>
            <MemberForm
              teamId={props.match.params.teamId}
              members={this.state.members[props.match.params.teamId]}
              onCancel={() => this.props.history.push(`${match.url}`)}
              onSave={this.onMemberAdd.bind(this)} 
              saveText="Add" />
          </ReactModal> } />      

        <Route path={`${match.url}/:teamId/edit`} render={props =>
          <ReactModal isOpen={true} contentLabel="Edit Team" className="modal" overlayClassName="modal-overlay">
            <TeamForm
              team={this.getTeamById(parseInt(props.match.params.teamId))}
              onCancel={() => this.props.history.push(`${match.url}`)}
              onSave={this.onTeamUpdate.bind(this, parseInt(props.match.params.teamId))} />
          </ReactModal> } />          

        { teamDeletion &&
          <ReactModal isOpen={true} contentLabel="Delete Team" className="modal" overlayClassName="modal-overlay">
            <h2>Delete Team</h2>
            <FormContainer
              onSave={this.confirmTeamDelete.bind(this)}
              onCancel={this.cancelTeamDelete.bind(this)}
              saveText="Confirm"
            >
              <div>{teamDeletion.contents}</div>
            </FormContainer>
          </ReactModal> }

        { memberDeletion &&
          <ReactModal isOpen={true} contentLabel="Remove Member" className="modal" overlayClassName="modal-overlay">
            <h2>Remove Member</h2>
            <FormContainer
              onSave={this.confirmMemberDelete.bind(this)}
              onCancel={this.cancelMemberDelete.bind(this)}
              saveText="Confirm"
            >
              <div>{memberDeletion.contents}</div>
            </FormContainer>
          </ReactModal> }

        <h1>Teams</h1>

        <Collapse className="accordion">
          { teams.map((team, i) =>
            <Panel
              className={"accordion-row"}
              header={
              <div className="accordion-row-header">
                <div className="accordion-row-header-contents">
                  <h4>{team.name}</h4>
                  <div className="nav">
                    <Link className="button secondary" to={`${match.url}/${team.id}/edit`}>Edit</Link>
                    <button className="secondary" onClick={this.startTeamDelete.bind(this, team)}>Delete</button>
                  </div>
                </div>
              </div>}
            >
              <p>{team.description}</p>
              { members[team.id] && members[team.id].length > 0 &&
                <MembersTable
                  teamId={team.id}
                  members={members[team.id]}
                  users={users[team.id]}
                  onDelete={this.startMemberDelete.bind(this)} 
                  onUpdate={this.confirmMemberUpdate.bind(this)}/> }
              { (!members[team.id] || members[team.id].length === 0) &&
                <p className="empty">This team has no members yet.</p> }
              <Link className="button secondary" to={`${match.url}/${team.id}/add-member`}>Add Member</Link>              
            </Panel>
            )
          }
        </Collapse>   
        { teams.length === 0 &&
          <span className="empty">Your expedition doesn't have any teams yet.</span> }
        <Link className="button" to={`${match.url}/new-team`}>Create New Team</Link>
      </div>
    )
  }
}
