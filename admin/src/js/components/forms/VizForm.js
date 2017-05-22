/* @flow */

import { get, set, compose } from 'safety-lens'
import { prop, _1, _2 } from 'safety-lens/es2015'
import React, { Component } from 'react'
import { FormItem } from './FormItem'
import { FormSelectItem } from './FormSelectItem'
import type { APIErrors } from '../../api/types';
import { RemoveIcon } from '../icons/Icons'
import MultiSelect from 'react-select';

import type {Viz, SelectionOperation} from '../../types/VizTypes'
import {_sourceCollections, _groupingOperation, _groupingOperationOp, _groupingOperationAttribute, _groupingOperationParam, _selectionOperationName, _selectionOperationOp, _selectionOperationSource} from '../Decorators.js'

import 'react-select/dist/react-select.css';

export class EditSelectionOperationComponent extends Component {
  props: {
    data: Viz;
    errors: ?APIErrors;
    initial_state: SelectionOperation;
    creator: Object;
  }

  state: SelectionOperation

  updateSelectionName: (Object) => void
  updateSelectionSource: (Object) => void
  updateSelectionOperation: (Object) => void

  constructor(props: *){
    super(props) 
    this.state = props.initial_state

    this.updateSelectionName = this.updateSelectionName.bind(this)
    this.updateSelectionSource = this.updateSelectionSource.bind(this)
    this.updateSelectionOperation = this.updateSelectionOperation.bind(this)
  }


  
  updateSelectionName(e: Object){
    const value = e.target.value
    const new_state = set(_selectionOperationName,value,this.state)
    this.setState(new_state)
  }
  
  updateSelectionSource(e: Object){
    let attr_name = e.target.value
    let attr = this.props.creator.getCollectionAttributes().find(a => a.name === attr_name)
    if(attr){
      const new_state = set(_selectionOperationSource,attr,this.state)
      this.setState(new_state)
    }
  }
  
  updateSelectionOperation(e: Object){
    const value = e.target.value
    const new_state = set(_selectionOperationOp,value,this.state)
    this.setState(new_state)
  }

  render(){
    const {creator, errors} = this.props
    const opts = [
      {text: "Average", value: "avg"},
      {text: "Maximum", value: "max"},
      {text: "Minimum", value: "min"},
      {text: "Median", value: "median"},
      {text: "First", value: "first"},
      {text: "Last", value: "last"},
      {text: "Sum", value: "sum"},
      {text: "Count", value: "count"},
      {text: "Match", value: "match_count"}
    ]
    const attrs = creator.getCollectionAttributes().map((a) => {
      return {text: a.name, value: a.name}
    })
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
            <button onClick={() => creator.setState({modal_open: false})}>Cancel</button>
        </div>
      </div>

    )
  }
}

export class GroupByComponent extends Component {
  props: {
    data: Viz;
    errors: ?APIErrors;        
    creator: Object;
  }

  state: {
    modal_open: boolean;
  }

  updateGroupingOperation: (Object) => void
  updateGroupingAttribute: (Object) => void
  updateGroupingParameter: (Object) => void
  updateSelectedCollections: (Object) => void
  

  constructor(props: *){
    super(props)
    this.updateGroupingOperation = this.updateGroupingOperation.bind(this)
    this.updateGroupingAttribute = this.updateGroupingAttribute.bind(this)
    this.updateGroupingParameter = this.updateGroupingParameter.bind(this)
    this.updateSelectedCollections = this.updateSelectedCollections.bind(this)

  }

  updateSelectedCollections(e: Object){
    console.log(e)
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
    let attr = this.props.creator.getcollectionattributes().find(a => a.name === attr_name)
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

    const collection_options = [
      {value: '1', label: 'Collection 1'},
      {value: '2', label: 'Collection 2'},
      {value: '3', label: 'Collection 3'}
    ]
    const selected_collections = ""
    
    return (
      <div>
        <div>
          <div>
            <MultiSelect
              labelKey={'Source Collections'}
              name={'source_collection'}
              value={selected_collections}
              options={collection_options}
              mutiple={true}
              onChange={this.updateSelectedCollections}
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
      </div>
    )
  }
}
