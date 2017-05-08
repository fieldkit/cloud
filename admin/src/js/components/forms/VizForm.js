/* @flow */

import { get, set, compose } from 'safety-lens'
import { prop, _1, _2 } from 'safety-lens/es2015'
import React, { Component } from 'react'
import { FormItem } from './FormItem'
import { FormSelectItem } from './FormSelectItem'
import type { APIErrors } from '../../api/types';
import { RemoveIcon } from '../icons/Icons'

import type {Viz} from '../Decorators.js'
import {_groupingOperation, _groupingOperationOp, _groupingOperationAttribute, _groupingOperationParam} from '../Decorators.js'

export class GroupByComponent extends Component {
  props: {
    data: Viz;
    errors: ?APIErrors;        
    creator: Object;
  }

  updateGroupingOperation: (string) => void
  updateGroupingAttribute: (string) => void
  updateGroupingParameter: (string) => void

  constructor(props: *){
    super(props)
    this.updateGroupingOperation = this.updateGroupingOperation.bind(this)
    this.updateGroupingAttribute = this.updateGroupingAttribute.bind(this)
    this.updateGroupingParameter = this.updateGroupingParameter.bind(this)
  }

  updateGroupingOperation(e: Object){
    let op = e.target.value
    if(op === "equal" || op === "within" || op === "peak"){
      const lens = compose(_groupingOperation, _groupingOperationOp)
      this.props.creator.update(lens,op) 
    }
  }

  updateGroupingAttribute(e: Object){
    let attr_name = e.target.value
    let attr = this.props.creator.getCollectionAttributes().find(a => a.name === attr_name)
    if(attr){
      const lens = compose(_groupingOperation, _groupingOperationAttribute)
      this.props.creator.update(lens,attr)
    }
  }
  
  updateGroupingParameter(e: Object){
    let value = e.target.value
    const lens = compose(_groupingOperation, _groupingOperationParam)
    this.props.creator.update(lens,value)
  }
// NOTE: CHANGE FROM PROJECT ATTRS TO SUM OVER COLLECTIONS 
  render() {
    const {data, errors, creator} = this.props
    const {grouping_operation } = data
    const grouping_ops = [
      {text: "equal", value: "equal"}, 
      {text:"within", value: "within"}, 
      {text: "peak", value: "peak"}
    ]
    const grouping_attrs = this.props.creator.getCollectionAttributes().map((a) => {
      return {text: a.name, value: a.name}
    })
    
    return (
      <div>
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
        {
          grouping_operation.operation === 'within' &&
          <FormItem
            labelText={'Grouping Param'}
            name={'grouping_param'}
            value={grouping_operation.parameter || ""}
            inline={true}
            errors={errors}
            onChange={this.updateGroupingParameter}
          />
        }
      </div>
    )
  }
}
