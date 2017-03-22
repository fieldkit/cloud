// @flow weak

import React, { Component } from 'react'

import { FormContainer } from '../containers/FormContainer';
import { errorsFor, slugify } from '../../common/util';

import type { APIErrors, APIProject, APIUser, APINewAdministrator } from '../../api/types';
import { FKApiClient } from '../../api/api';

type Props = {
  project: APIProject,
//   member?: APIMember,

//   cancelText?: string;
//   saveText?: ?string;
//   onCancel?: () => void;
//   onSave: (teamId: number, e: APINewMember) => Promise<?APIErrors>; 
}

export class AdministratorForm extends Component {
  props: Props;
  state: {
    users: APIUser[],    
    userId: number,
    saveDisabled: boolean,
    errors: ?APIErrors    
  }

  constructor(props: Props) {
    super(props);
    this.state = {
      users: [],
      userId: -1,
      saveDisabled: true,
      errors: null      
    }
  }

  async save() {
  //   const errors = await this.props.onSave(this.props.teamId, {
  //     user_id: parseInt(this.state.userId),
  //     role: this.state.role
  //   });
  //   if (errors) {
  //     this.setState({ errors });
  //   }
  }

  render () {
    
    const { users, userId } = this.state;

    return (
      <FormContainer
        onSave={this.save.bind(this)}>



      </FormContainer>
    )
  }
}