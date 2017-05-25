/* @flow */
import React, { Component } from 'react';
import { set } from 'safety-lens';
import { prop } from 'safety-lens/es2015';
import { FormItem } from '../forms/FormItem';
import { FormSelectItem } from '../forms/FormSelectItem';
import type { Lens_ } from 'safety-lens';
import type { APIErrors } from '../../api/types';
import type { Viz, SelectionOperation, Op } from '../../types/VizTypes';

export const _selectionOperationName: Lens_<SelectionOperation,string> = prop("value_name");
export const _selectionOperationSource: Lens_<SelectionOperation,Attr> = prop("source_attribute");
export const _selectionOperationOp: Lens_<SelectionOperation,Op> = prop("operation");

export default class VizSelectionForm extends Component {
  props: {
    data: Viz,
    errors: ?APIErrors,
    initial_state: SelectionOperation,
    creator: Object,
  };

  state: SelectionOperation;

  updateSelectionName: Object => void;
  updateSelectionSource: Object => void;
  updateSelectionOperation: Object => void;

  constructor(props: *) {
    super(props);
    this.state = props.initial_state;

    this.updateSelectionName = this.updateSelectionName.bind(this);
    this.updateSelectionSource = this.updateSelectionSource.bind(this);
    this.updateSelectionOperation = this.updateSelectionOperation.bind(this);
  }

  updateSelectionName(e: Object) {
    const value = e.target.value;
    const new_state = set(_selectionOperationName, value, this.state);
    this.setState(new_state);
  }

  updateSelectionSource(e: Object) {
    let attr_name = e.target.value;
    let attr = this.props.creator
      .getCollectionAttributes()
      .find(a => a.name === attr_name);
    if (attr) {
      const new_state = set(_selectionOperationSource, attr, this.state);
      this.setState(new_state);
    }
  }

  updateSelectionOperation(e: Object) {
    const value = e.target.value;
    const new_state = set(_selectionOperationOp, value, this.state);
    this.setState(new_state);
  }

  render() {
    const { creator, errors } = this.props;
    const opts = [
      { text: 'Average', value: 'avg' },
      { text: 'Maximum', value: 'max' },
      { text: 'Minimum', value: 'min' },
      { text: 'Median', value: 'median' },
      { text: 'First', value: 'first' },
      { text: 'Last', value: 'last' },
      { text: 'Sum', value: 'sum' },
      { text: 'Count', value: 'count' },
      { text: 'Match', value: 'match_count' },
    ];
    const attrs = creator.getCollectionAttributes().map(a => {
      return { text: a.name, value: a.name };
    });
    return (
      <div>
        <FormItem
          labelText={'Name'}
          name={'name'}
          value={this.state.value_name}
          inline={false}
          errors={errors}
          onChange={this.updateSelectionName}
        />
        <FormSelectItem
          labelText={'Selection Operation'}
          name={'selection_operation'}
          value={this.state.operation}
          inline={false}
          firstOptionText={'Select'}
          options={opts}
          errors={errors}
          onChange={this.updateSelectionOperation}
        />
        <FormSelectItem
          labelText={'Selection Attribute'}
          name={'selection_attribute'}
          value={this.state.source_attribute.name}
          inline={false}
          firstOptionText={'Select'}
          options={attrs}
          errors={errors}
          onChange={this.updateSelectionSource}
        />
        <div>
          <button onClick={() => creator.addSelection(this.state)}>Save</button>
          <button onClick={() => creator.setState({ modal_open: false })}>
            Cancel
          </button>
        </div>
      </div>
    );
  }
}
