/* @flow */

import React, { Component } from 'react'
import { FormContainer } from '../containers/FormContainer';
import type { APIErrors, APICollection, APINewCollection } from '../../api/types';

type Props = {
  collection?: ?APICollection,

  cancelText?: string;
  saveText?: ?string;
  onCancel?: () => void;
  onSave: (e: APINewCollection, teamId?: number) => Promise<?APIErrors>;
}


export class PCForm extends Component {
  attributes: { [string]: Attr };
  props: Props;
  state: {
    errors: ?APIErrors
  }

  constructor(props: Props) {
    super(props);
  }

  componentWillReceiveProps(nextProps: Props) {
  }

  save(){

  }

  render() {
    
    return (
      <FormContainer
        onSave={this.save.bind(this)}
        onCancel={this.props.onCancel}
        saveText={this.props.saveText}
        cancelText={this.props.cancelText}>

      </FormContainer>
    )
  }
}
