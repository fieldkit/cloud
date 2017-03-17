// @flow weak

import React, { Component } from 'react'

import { FormContainer } from '../containers/FormContainer';
import { errorsFor, slugify } from '../../common/util';

import type { APIErrors } from '../../api/types';

type Props = {
  name?: string,
  slug?: string,
  description?: string,

  cancelText?: string;
  saveText?: ?string;
  onCancel?: () => void;
  onSave: (n: string, s: string, d: string) => Promise<?APIErrors>;
}

export class ProjectForm extends Component {
  props: Props;
  state: {
    name: string,
    slug: string,
    description: string,
    slugHasChanged: boolean,
    saveDisabled: boolean,
    errors: ?APIErrors
  };

  constructor(props: Props) {
    super(props)
    this.state = {
      name: this.props.name || '',
      slug: this.props.slug || '',
      description: this.props.description || '',
      slugHasChanged: false,
      saveDisabled: false,
      errors: null
    }
  }

  componentWillReceiveProps(nextProps: Props) {
    this.setState({
      name: nextProps.name || '',
      slug: nextProps.slug || '',
      description: nextProps.description || '',
      slugHasChanged: false,
      errors: null
    });
  }

  async save() {
    const errors = await this.props.onSave(this.state.name, this.state.slug, this.state.description);
    if (errors) {
      this.setState({ errors });
    }
  }

  handleNameChange(event) {
    const v = event.target.value;
    let stateUpdate = { name: v };
    if (!this.state.slugHasChanged) {
      stateUpdate = { slug: slugify(v), ...stateUpdate };
    }
    this.setState(stateUpdate);
  }

  handleSlugChange(event) {
    const v = event.target.value;
    this.setState({ slug: v, slugHasChanged: true });
  }

  handleInputChange(event) {
    const target = event.target;
    const value = target.type === 'checkbox' ? target.checked : target.value;
    const name = target.name;

    this.setState({ [name]: value });
  }

  render() {
    return (
      <FormContainer
        onSave={this.save.bind(this)}
        onCancel={this.props.onCancel}
        saveText={this.props.saveText}
        cancelText={this.props.cancelText}>

        <div className="form-group">
          <label htmlFor="name">Name</label>
          <input type="text" name="name" className='lg' value={this.state.name} onChange={this.handleNameChange.bind(this)} />
          { errorsFor(this.state.errors, 'name') }
        </div>

        <div className="url-preview">
          <p className="label">Your project will be available at the following address:</p>
          <p className="url">
            <input type="text" name="slug" className='lg slug' value={this.state.slug} onChange={this.handleSlugChange.bind(this)} />
            {/* TODO: replace with something that handles alternative domains */}
            <span className="domain">.fieldkit.org/</span>
          </p>
          { errorsFor(this.state.errors, 'path') }
        </div>

        <div className="form-group">
          <label htmlFor="description">Description</label>
          <input type="text" name="description" className='lg' value={this.state.description} onChange={this.handleInputChange.bind(this)} />
          { errorsFor(this.state.errors, 'description') }
        </div>
      </FormContainer>
    );
  }
}
