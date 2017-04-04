// @flow weak

import React, { Component } from 'react'

import { FormContainer } from '../containers/FormContainer';
import { FormItem } from './FormItem';
import { errorsFor, slugify, fkHost } from '../../common/util';

import type { APIErrors, APINewExpedition } from '../../api/types';

type Props = {
  projectSlug: string,
  name?: string,
  slug?: string,
  description?: string,

  cancelText?: string;
  saveText?: ?string;
  onCancel?: () => void;
  onSave: (e: APINewExpedition) => Promise<?APIErrors>;
}

export class ProjectExpeditionForm extends Component {
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
    const errors = await this.props.onSave({
      name: this.state.name,
      slug: this.state.slug,
      description: this.state.description
    });
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

  render () {
    const { projectSlug } = this.props;
    const {
      name,
      slug,
      description,
      errors
    } = this.state;

    return (
      <FormContainer
        onSave={this.save.bind(this)}
        onCancel={this.props.onCancel}
        saveText={this.props.saveText}
        cancelText={this.props.cancelText}>

        <FormItem name="name" value={name} labelText="name" className="lg" onChange={this.handleNameChange.bind(this)} errors={errors} />

        <div className="url-preview">
          <p className="label">Your expedition will be available at the following address:</p>
          <p className="url">
            {/* TODO: replace with something that handles alternative domains */}
            {`//${projectSlug}.${fkHost()}/${slug}`}
          </p>
        </div>

        <FormItem name="description" value={description} labelText="description" className="lg" onChange={this.handleInputChange.bind(this)} errors={errors} />
      </FormContainer>
    )
  }
}
