// @flow weak

import React, { Component } from 'react'
import slug from 'slug';

import { FormContainer } from '../containers/FormContainer';
import { errorsFor } from '../../common/util';
import type { ErrorMap } from '../../common/util';

type Props = {
  name?: string,
  path?: string,
  description?: string,

  cancelText?: string;
  saveText?: ?string;
  onCancel?: () => void;
  onSave: (n: string, d: string) => Promise<?ErrorMap>;
}

export class ProjectForm extends Component {
  props: Props;
  state: {
    name: string,
    path: string,
    description: string,
    saveDisabled: boolean,
    errors: ?ErrorMap
  };

  constructor(props: Props) {
    super(props)
    this.state = {
      name: this.props.name || '',
      path: this.props.path || '',
      description: this.props.description || '',
      saveDisabled: false,
      errors: null
    }
  }

  async save() {
    const errors = await this.props.onSave(this.state.name, this.state.description);
    if (errors) {
      this.setState({ errors });
    }
  }

  handleNameChange(event) {
    const v = event.target.value;
    this.setState({
      name: v,
      path: slug(v).toLowerCase()
    });
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

        <h1>Create a new project</h1>

        <div className="form_group">
          <label htmlFor="name">Pick a name for your project:</label>
          <input type="text" name="name" className='lg' value={this.state.name} onChange={this.handleNameChange.bind(this)} />
          { errorsFor(this.state.errors, 'name') }
        </div>

        <div className="url-preview">
          <p className="label">Your project will be available at the following address:</p>
          <p className="url">
            {/* TODO: replace with something that handles alternative domains */}
            {`https://${this.state.path}.fieldkit.org/`}
          </p>
          { errorsFor(this.state.errors, 'path') }
        </div>

        <div className="form_group">
          <label htmlFor="description">Enter a description for your project:</label>
          <input type="text" name="description" className='lg' value={this.state.description} onChange={this.handleInputChange.bind(this)} />
          { errorsFor(this.state.errors, 'description') }
        </div>
      </FormContainer>
    );
  }
}
