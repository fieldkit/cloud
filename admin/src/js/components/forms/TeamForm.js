// @flow weak

import React, { Component } from 'react'

import { FormContainer } from '../containers/FormContainer';
import { errorsFor, slugify } from '../../common/util';

import type { APIErrors, APITeam, APINewTeam } from '../../api/types';

type Props = {
  team?: ?APITeam,

  cancelText?: string;
  saveText?: ?string;
  onCancel?: () => void;
  onSave: (e: APINewTeam, teamId?: number) => Promise<?APIErrors>;
}

export class TeamForm extends Component {
  props: Props;
  /*flow-include
  state: {
    ...$Exact<APITeam>,
    slugHasChanged: boolean,    
    saveDisabled: boolean,
    errors: ?APIErrors
  };
  */
  constructor(props: Props) {
    super(props)
    this.state = {
      name: '',
      slug: '',
      description: '',
      ...props.team,
      slugHasChanged: false,
      saveDisabled: false,
      errors: null
    }
  }

  componentWillReceiveProps(nextProps: Props) {
    this.setState({
      name: '',
      slug: '',
      description: '',
      ...nextProps.team,
      slugHasChanged: false,
      errors: null
    });
  }

  async save(teamId?: number) {
    const errors = await this.props.onSave({
        name: this.state.name,
        slug: this.state.slug,
        description: this.state.description
      }, teamId);

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
    return (
      <FormContainer
        onSave={this.save.bind(this)}
        onCancel={this.props.onCancel}
        saveText={this.props.saveText}
        cancelText={this.props.cancelText}>

        <div className="form-group">
          <label htmlFor="name">Name</label>
          <input type="text" name="name" className="lg" value={this.state.name} onChange={this.handleNameChange.bind(this)} />
          { errorsFor(this.state.errors, 'name') }
        </div>

        <div className="form-group">
          <label htmlFor="description">Description</label>
          <input type="text" name="description" className="lg" value={this.state.description} onChange={this.handleInputChange.bind(this)} />
          { errorsFor(this.state.errors, 'description') }
        </div>
      </FormContainer>
    )
  }
}