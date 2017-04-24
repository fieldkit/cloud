/* @flow */

import React, { Component } from 'react'
import Dropdown, { DropdownTrigger, DropdownContent } from 'react-simple-dropdown';
import { FormItem } from './forms/FormItem'
import type { APIErrors } from '../api/types';

import type {StringFilter, DateFilter, NumFilter} from './Collection'

export class StringFilterComponent extends Component {
  props: {
    data: StringFilter;
    creator: Object;
    errors: ?APIErrors;        
  }

  constructor(props: StringFilter) {
    super(props);
  }

  render() {
    const operations = ["contains","does not contain","matches","exists"].map((o,i) => <option value={o} key={i}>{o}</option>)
    const { data, creator, errors } = this.props
    let value_field 
    if(data.operation === "exists"){
      value_field = null
    } else if(data.options.length > 0){
      let options = data.options.map((name,i) => {
        return (
          <option value={name} key={i}>{name}</option>
        )
      })
      value_field = (
        <div>
          <span className="filter-body-label">Value: </span>
          <select className="value-body-select" value={data.query} onChange={(e) => creator.updateFilter(data,{"query":e.target.value})}>
            {options}
          </select>
        </div>
      )
    } else {
      value_field = (
        <FormItem
          labelText={'Value'}
          name={'value'}
          className={'value'}
          value={data.query}
          errors={errors}
          onChange={(e) => creator.updateFilter(data,{"query":e.target.value})}
        />
      )
    }
    
    return (
      <div className="fk-filter fk-guidfilter" key={data.id}>
        <div className="filter-body">
          <div>
            <span className="filter-body-label">Condition: </span>
            <select className="filter-body-select" value={data.operation} onChange={(e) => creator.updateFilter(data,{"operation":e.target.value})}>
              {operations}
            </select>
          </div>
          <div>
            {value_field}
          </div>
          <div className="filter-body-buttons">
            <button className="filter-body-cancel" onClick={() => creator.deleteFilter(data)}>Delete</button>
          </div>
        </div>
      </div>
    )
  }
}

export class NumFilterComponent extends Component {
  props: {
    data: NumFilter;
    creator: Object;
    errors: ?APIErrors;
  }

  constructor(props: NumFilter) {
    super(props);
  }

  render() {
    const operations = ["GT","LT","EQ","notch"].map((o,i) => <option value={o} key={i}>{o.toUpperCase()}</option>)
    const { creator, data, errors } = this.props
    
    return (
      <div className="fk-filter fk-guidfilter">
        <div className="filter-body">
          <div>
            <span className="filter-body-label">Condition: </span>
            <select className="filter-body-select" value={data.operation} onChange={(e) => creator.updateFilter(data,{"operation":e.target.value})}>
              {operations}
            </select>
          </div>
          <FormItem
            labelText={'Value'}
            name={'value'}
            className={'value'}
            type={'number'}
            value={data.query}
            errors={errors}
            onChange={(e) => creator.updateFilter(data,{"query":e.target.value})}
          />
          <div className="filter-body-buttons">
            <button className="filter-body-cancel" onClick={() => creator.deleteFilter(data)}>Delete</button>
          </div>
        </div>
      </div>
    )
  }
}

export class DateFilterComponent extends Component {
  props: {
    data: DateFilter;
    creator: Object;
    errors: ?APIErrors;        
  }

  constructor(props: NumFilter) {
    super(props);
  }

  render() {
    const operations = ["before","after","within"].map((o,i) => <option value={o} key={i}>{o}</option>)
    const { data, creator, errors } = this.props
    
    return (
      <div className="fk-filter fk-guidfilter">
        <div className="filter-body">
          <div>
            <span className="filter-body-label">Condition: </span>
            <select className="filter-body-select" value={data.operation} onChange={(e) => creator.updateFilter(data,{"operation":e.target.value})}>
              {operations}
            </select>
          </div>
          <FormItem
            labelText={data.operation === 'within' && 'Start Date' || 'Date'}
            name={'date'}
            type={'date'}
            value={data.date}
            errors={errors}
            onChange={(e) => creator.updateFilter(data,{"date":e.target.value})}
          />
          { data.operation === 'within' &&
            <FormItem
              labelText={'End Date'}
              name={'date'}
              type={'date'}
              value={data.within}
              errors={errors}
              onChange={(e) => creator.updateFilter(data,{"date":e.target.value})}
            />        
          }
          <div className="filter-body-buttons">
            <button className="filter-body-cancel" onClick={() => creator.deleteFilter(data)}>Delete</button>
          </div>
        </div>
      </div>
    )
  }
}
