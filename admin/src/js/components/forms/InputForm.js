// @flow weak

import React, { Component } from 'react'

import { FormContainer } from '../containers/FormContainer';
import { errorsFor, slugify } from '../../common/util';

import type { APIErrors, APIInput, APINewInput, APIInputType } from '../../api/types';

type Props = {
  projectSlug: string,
  input?: ?APIInput;

  cancelText?: string;
  saveText?: ?string;
  onCancel?: () => void;
  onSave: (i: APINewInput) => Promise<?APIErrors>;
}

const INPUT_TYPES = [
  { value: "webhook", display: "Webhook" },
  { value: "twitter", display: "Twitter" },
]

function initialStateFromProps(props: Props) {
  return {
    saveDisabled: false,
    errors: null,

    name: '',
    type: INPUT_TYPES[0].value,
    active: true,

    ...props.input
  }
}

export class InputForm extends Component {
  props: Props;
  state: {
    name: string,
    type: APIInputType,
    active: boolean,

    saveDisabled: boolean,
    errors: ?APIErrors
  };

  constructor(props: Props) {
    super(props)
    this.state = initialStateFromProps(props);
  }

  componentWillReceiveProps(nextProps: Props) {
    this.state = initialStateFromProps(nextProps);
  }

  async save() {
    const errors = await this.props.onSave({
      name: this.state.name,
      type: this.state.type,
      active: this.state.active
    });

    if (errors) {
      this.setState({ errors });
    }
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

        <div className="form_group">
          <label htmlFor="name">Name:</label>
          <input type="text" name="name" className='lg' value={this.state.name} onChange={this.handleInputChange.bind(this)} />
          { errorsFor(this.state.errors, 'name') }
        </div>

        <div className="form_group">
          <label htmlFor="type">Type:</label>
          <select name="type" className='lg' value={this.state.type} onChange={this.handleInputChange.bind(this)}>
            { INPUT_TYPES.map((o, i) => <option key={i} value={o.value}>{o.display}</option>) }
          </select>
          { errorsFor(this.state.errors, 'type') }
        </div>

        <div className="form_group">
          <label htmlFor="active">Active:</label>
          <input type="checkbox" name="active" className='lg' checked={this.state.active} onChange={this.handleInputChange.bind(this)} />
          { errorsFor(this.state.errors, 'active') }
        </div>
      </FormContainer>
    );
  }
}
