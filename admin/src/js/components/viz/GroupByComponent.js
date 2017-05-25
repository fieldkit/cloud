/* @flow */
import React, { Component } from 'react';
import { get, set, compose } from 'safety-lens';
import { prop, _1, _2 } from 'safety-lens/es2015';
import { FormItem } from '../forms/FormItem';
import { FormSelectItem } from '../forms/FormSelectItem';
import MultiSelect from 'react-select';
import type { Lens_ } from 'safety-lens';
import type { APIErrors } from '../../api/types';
import type { Viz, GroupingOperation, GroupingOperationType } from '../../types/VizTypes';

import 'react-select/dist/react-select.css';

export const _sourceCollections: Lens_<Viz,string[]> = prop('source_collections');
export const _groupingOperation: Lens_<Viz,GroupingOperation> = prop('grouping_operation');
export const _groupingOperationOp: Lens_<GroupingOperation,GroupingOperationType> = prop('operation');
export const _groupingOperationParam: Lens_<GroupingOperation,?number> = prop('parameter');
export const _groupingOperationAttribute: Lens_<GroupingOperation,Attr> = prop('source_attribute');

export default class GroupByComponent extends Component {
  props: {
    data: Viz,
    errors: ?APIErrors,
    creator: Object,
  };

  updateGroupingOperation: Object => void;
  updateGroupingAttribute: Object => void;
  updateGroupingParameter: Object => void;
  updateSelectedCollections: Object => void;

  constructor(props: *) {
    super(props);
    this.updateGroupingOperation = this.updateGroupingOperation.bind(this);
    this.updateGroupingAttribute = this.updateGroupingAttribute.bind(this);
    this.updateGroupingParameter = this.updateGroupingParameter.bind(this);
    this.updateSelectedCollections = this.updateSelectedCollections.bind(this);
  }

  updateSelectedCollections(e: Object) {
    const ids = e.split(',');
    this.props.creator.update(_sourceCollections, ids);
  }

  updateGroupingOperation(e: Object) {
    let op = e.target.value;
    if (op === 'equal' || op === 'within' || op === 'peak') {
      const lens = compose(_groupingOperation, _groupingOperationOp);
      this.props.creator.update(lens, op);
    }
  }

  updateGroupingAttribute(e: Object) {
    let attr_name = e.target.value;
    let attr = this.props.creator
      .getcollectionattributes()
      .find(a => a.name === attr_name);
    if (attr) {
      const lens = compose(_groupingOperation, _groupingOperationAttribute);
      this.props.creator.update(lens, attr);
    }
  }

  updateGroupingParameter(e: Object) {
    let value = e.target.value;
    const lens = compose(_groupingOperation, _groupingOperationParam);
    this.props.creator.update(lens, value);
  }
  // NOTE: CHANGE FROM PROJECT ATTRS TO SUM OVER COLLECTIONS
  render() {
    const { data, errors, creator } = this.props;
    const { grouping_operation, source_collections } = data;
    const grouping_ops = [
      { text: 'equal', value: 'equal' },
      { text: 'within', value: 'within' },
      { text: 'peak', value: 'peak' },
    ];
    const grouping_attrs = this.props.creator
      .getCollectionAttributes()
      .map(a => {
        return { text: a.name, value: a.name };
      });

    const collection_options = [
      { value: 1, label: 'Collection 1' },
      { value: 2, label: 'Collection 2' },
      { value: 3, label: 'Collection 3' },
    ];
    const selected_collections = source_collections.join(',');
    return (
      <div>
        <div>
          <div>
            <label>Source Collections:</label>
            <MultiSelect
              name={'source_collection'}
              value={selected_collections}
              options={collection_options}
              multi={true}
              simpleValue={true}
              joinValues={true}
              onChange={this.updateSelectedCollections}
              multiple={true}
              options={collection_options}
            />
          </div>
          <FormSelectItem
            labelText={'Grouping Operation'}
            name={'grouping_operation'}
            value={grouping_operation.operation}
            inline={true}
            firstOptionText={'Select'}
            options={grouping_ops}
            errors={errors}
            onChange={this.updateGroupingOperation}
          />
          <FormSelectItem
            labelText={'Grouping Attribute'}
            name={'grouping_attribute'}
            value={grouping_operation.source_attribute.name}
            inline={true}
            firstOptionText={'Select'}
            options={grouping_attrs}
            errors={errors}
            onChange={this.updateGroupingAttribute}
          />
          {grouping_operation.operation === 'within' &&
            <FormItem
              labelText={'Grouping Param'}
              name={'grouping_param'}
              value={grouping_operation.parameter || ''}
              inline={true}
              errors={errors}
              onChange={this.updateGroupingParameter}
            />}
        </div>
      </div>
    );
  }
}
