// @flow weak

import React, { Component } from 'react'

import { FormContainer } from '../containers/FormContainer';
import { errorsFor, slugify } from '../../common/util';

import type { APIErrors, APIUser, APIFieldkitInput, APINewFieldkitInput, APIMutableInput, APINewTwitterInput } from '../../api/types';

type Props = {
  fieldkitInput?: ?APIFieldkitInput,
  users: APIUser[],
  inputType: string,

  cancelText?: string;
  saveText?: ?string;
  onCancel?: () => void;
  onSave: (i: APINewFieldkitInput | APINewTwitterInput) => Promise<?APIErrors>;
}

export class InputForm extends Component {
  props: Props;

  /*flow-include
  state: {
    ...$Exact<APIFieldkitInput>,  
    users: APIUser[],
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
      users: [],
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
      users: nextProps.users,
      saveDisabled: false,
      errors: null
    });
  }

  async save() {
    const { name } = this.state;
    const errors = await this.props.onSave({name: name});
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