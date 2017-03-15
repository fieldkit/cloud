// @flow weak

import React, { Component } from 'react'

import { FormContainer } from '../containers/FormContainer';
import { errorsFor, slugify } from '../../common/util';
import type { ErrorMap } from '../../common/util';

type Props = {
  projectSlug: string,
  name?: string,
  description?: string,

  cancelText?: string;
  saveText?: ?string;
  onCancel?: () => void;
  onSave: (n: string, d: string) => Promise<?ErrorMap>;
}

export class ProjectExpeditionForm extends Component {
  props: Props;
  state: {
    name: string,
    path: string,
    description: string,
    saveDisabled: boolean,
    errors: ?ErrorMap
  };

  constructor(props: Props) {
    super(props);
    const name = props.name || '';
    const description = props.description || '';
    this.state = {
      name,
      path: slugify(name),
      description,
      saveDisabled: false,
      errors: null
    }
  }

  componentWillReceiveProps(nextProps: Props) {
    const name = nextProps.name || '';
    const description = nextProps.description || '';
    this.setState({ name, path: slugify(name), description });
  }

  async save() {
    const errors = await this.props.onSave(this.state.name, this.state.description);
    if (errors) {
      this.setState({ errors });
    }
  }

  handleExpeditionNameChange(event) {
    const v = event.target.value;
    this.setState({
      name: v,
      path: slugify(v)
    });
  }

  handleInputChange(event) {
    const target = event.target;
    const value = target.type === 'checkbox' ? target.checked : target.value;
    const name = target.name;

    this.setState({
      [name]: value
    });
  }

  render () {
    return (
      <FormContainer
        onSave={this.save.bind(this)}
        onCancel={this.props.onCancel}
        saveText={this.props.saveText}
        cancelText={this.props.cancelText}>

        <div className="form-group">
          <label htmlFor="name">Pick a name for this new expedition:</label>
          <input type="text" name="name" className="text_lg" value={this.state.name} onChange={this.handleExpeditionNameChange.bind(this)} />
          <div className="hint">You can change this later if you want.</div>
          { errorsFor(this.state.errors, 'name') }
        </div>

        <div className="url-preview">
          <p className="label">Your expedition will be available at the following address:</p>
          <p className="url">
            {/* TODO: replace with something that handles alternative domains */}
            {`https://${this.props.projectSlug}.fieldkit.org/${this.state.path}`}
          </p>
        </div>
      </FormContainer>
    )
  }
}
