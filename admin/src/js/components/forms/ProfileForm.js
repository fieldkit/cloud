// @flow weak

import React, { Component } from 'react'

import { FormContainer } from '../containers/FormContainer';
import { FormItem } from './FormItem';
import { errorsFor } from '../../common/util';

import type { APIErrors, APIBaseUser } from '../../api/types';

type Props = {
  name: string,
  bio: string,

  cancelText?: string;
  saveText?: ?string;
  onCancel?: () => void;
  onSave: (name: string, bio: string) => Promise<?APIErrors>;
}

export class ProfileForm extends Component {
  props: Props;
  state: {
    name: string,
    bio: string,
    saveDisabled: boolean,
    errors: ?APIErrors
  };

  constructor(props: Props) {
    super(props)
    this.state = {
      name: this.props.name || '',
      bio: this.props.bio || '',
      saveDisabled: true,
      errors: null
    }
  }

  componentWillReceiveProps(nextProps: Props) {
    this.setState({
      name: nextProps.name || '',
      bio: nextProps.bio || '',
      saveDisabled: true,
      errors: null
    });
  }

  handleInputChange(event) {
    const target = event.target;
    const value = target.type === 'checkbox' ? target.checked : target.value;
    const name = target.name;
    this.setState({ [name]: value, saveDisabled: false });
  }

  async save() {
    const { name, bio } = this.state;
    const errors = await this.props.onSave(name, bio);
    if (errors) {
      this.setState({ errors });
    }
  }

  render() {
    const { name, bio, saveDisabled, errors } = this.state;

    return (
      <FormContainer
        onSave={this.save.bind(this)}
        onCancel={this.props.onCancel}
        saveText={this.props.saveText}
        cancelText={this.props.cancelText}
        saveDisabled={saveDisabled}>
            <FormItem labelText="Name" name="name" value={name} onChange={this.handleInputChange.bind(this)} errors={errors}  className="lg" />
            <FormItem labelText="Bio" name="bio" value={bio} onChange={this.handleInputChange.bind(this)} errors={errors}  className="lg" />
      </FormContainer>
    );
  }
}
