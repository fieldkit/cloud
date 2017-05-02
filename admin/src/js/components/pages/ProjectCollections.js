// @flow weak

import React, { Component } from 'react'
import { Route, Link } from 'react-router-dom'
import ReactModal from 'react-modal';

import { ProjectCollectionsForm } from '../forms/ProjectCollectionsForm';
import { FKApiClient } from '../../api/api';
import { FormContainer } from '../containers/FormContainer';
import Collapse, { Panel } from 'rc-collapse';
import 'rc-collapse/assets/index.css';

import type { APIProject, APIExpedition, APICollection, APINewCollection } from '../../api/types';

type Props = {
  project: APIProject;
  expeditions: APIExpedition[],

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
    const collectionsRes = await FKApiClient.get().getCollectionsByProjectSlug(this.props.project.slug);
    if (collectionsRes.type === 'ok' && collectionsRes.payload) {
      const { collections } = collectionsRes.payload;
      this.setState({ collections });
    }
  }

  async onCollectionCreate(c: APINewCollection) {
    const { project, match } = this.props;
    const collectionRes = await FKApiClient.get().createCollection(project.id, c);
    if (collectionRes.type === 'ok') {
      await this.loadCollections();
      this.props.history.push(`${match.url}`);
    } else {
      return collectionRes.errors;
    }
  }

  async onCollectionUpdate(collectionId: number, c: APINewCollection) {
    const collectionRes = await FKApiClient.get().updateCollection(collectionId, c);
    if(collectionRes.type === 'ok' && collectionRes.payload) {
      await this.loadCollections();
    } else {
      return collectionRes.errors;
    }
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

  getCollectionById (id: number): ?APICollection {
    const { collections } = this.state;
    return collections.find(collection => collection.id === id);
  }  

  render() {
    const { match, expeditions } = this.props;
    const { collections, collectionDeletion } = this.state;

    return (
      <div className="collections">

        <Route path={`${match.url}/new-collection`} render={() =>
          <ReactModal isOpen={true} contentLabel="Create New Collection" className="modal" overlayClassName="modal-overlay">
            <h2>Create New Collection</h2>
            <ProjectCollectionsForm
              onCancel={() => this.props.history.push(`${match.url}`)}
              onSave={this.onCollectionCreate.bind(this)} />
          </ReactModal> } />

        <Route path={`${match.url}/:collectionId/edit`} render={props =>
          <ReactModal isOpen={true} contentLabel="Edit Collection" className="modal" overlayClassName="modal-overlay">
            <ProjectCollectionsForm
              collection={this.getCollectionById(parseInt(props.match.params.collectionId))}
              onCancel={() => this.props.history.push(`${match.url}`)}
              onSave={this.onCollectionUpdate.bind(this, parseInt(props.match.params.collectionId))} />
          </ReactModal> } />          

        { collectionDeletion &&
          <ReactModal isOpen={true} contentLabel="Delete Collection" className="modal" overlayClassName="modal-overlay">
            <h2>Delete Team</h2>
            <FormContainer
              onSave={this.confirmCollectionDelete.bind(this)}
              onCancel={this.confirmCollectionDelete.bind(this)}
              saveText="Confirm"
            >
              <div>{collectionDeletion.contents}</div>
            </FormContainer>
          </ReactModal> }

        <h1>Collections</h1>
        <div id="collections">
          { collections.map((collection, i) =>
            <div className="accordion-row-header">
              <div className="accordion-row-header-contents">
                <h4>{collection.name}</h4>
                <div className="nav">
                  <Link className="button secondary" to={`${match.url}/${collection.id}/edit`}>Edit</Link>
                  <button className="secondary" onClick={this.startCollectionDelete.bind(this, collection)}>Delete</button>
                </div>
              </div>
            </div>
          )}
          { collections.length === 0 &&
            <span className="empty">This project doesn't have any collections yet.</span> }
        </div>
        <Link className="button" to={`${match.url}/new-collection`}>Create New Collection</Link>
      </div>
    )
  }
}
