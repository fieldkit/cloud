// @flow weak

import React, { Component } from 'react'

import { FormContainer } from '../containers/FormContainer';
import { errorsFor, slugify } from '../../common/util';

import type { APIErrors, APIFieldkitInput, APINewFieldkitInput } from '../../api/types';

type Props = {
  fieldkitInput?: ?APIFieldkitInput,

  cancelText?: string;
  saveText?: ?string;
  onCancel?: () => void;
  onSave: (e?: APINewFieldkitInput) => Promise<?APIErrors>;
}

export class FieldkitInputForm extends Component {
  props: Props;

  /*flow-include
  state: {
    ...$Exact<APIFieldkitInput>,  
    saveDisabled: boolean,
    errors: ?APIErrors
  };
  */  
  constructor(props: Props) {
    super(props)
    this.state = {
      team_id: -1,
      user_id: -1,
      name: '',
      ...props.fieldkitInput,
      saveDisabled: false,
      errors: null
    }
  }

  componentWillReceiveProps(nextProps: Props) {
    this.setState({
      team_id: -1,
      user_id: -1,
      name: '',
      ...nextProps.fieldkitInput,
      saveDisabled: false,
      errors: null
    });
  }

  async save(fieldkitInput?: APINewFieldkitInput) {
    const errors = await this.props.onSave(fieldkitInput);
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

  render () {
    return (
      <FormContainer
        onSave={this.save.bind(this)}
        onCancel={this.props.onCancel}
        saveText={this.props.saveText}
        cancelText={this.props.cancelText}>

        <div className="form-group">
          <label htmlFor="name">Name</label>
          <input type="text" name="name" className="lg" value={this.state.name} onChange={this.handleInputChange.bind(this)} />
          { errorsFor(this.state.errors, 'name') }
        </div>
      </FormContainer>
    )
  }
}