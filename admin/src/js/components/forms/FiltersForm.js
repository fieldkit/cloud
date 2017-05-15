/* @flow */

import React, { Component } from 'react'
import { FormItem } from './FormItem'
import { FormSelectItem } from './FormSelectItem'
import type { APIErrors } from '../../api/types';
import { RemoveIcon } from '../icons/Icons'

import type {StringFilter, DateFilter, NumFilter} from '../Collection'

export class StringFilterComponent extends Component {
  props: {
    data: StringFilter;
    creator: Object;
    errors: ?APIErrors;        
  }

  render() {
    const operations = ["contains","does not contain","matches","exists"].map((o,i) => { return { value: o, text: o } })
    const { data, creator, errors } = this.props
    let value_field 
    if(data.operation === "exists"){
      value_field = null
    } else if(data.options.length > 0){
      let options = data.options.map((name,i) => {
        return { value: name, text: name }
      })
      value_field = (
        <FormSelectItem
          labelText={'Value'}
          name={'value'}
          value={data.query}
          inline={true}
          firstOptionText={'Select'}
          options={options}
          errors={errors}
          onChange={(e) => creator.updateFilter(data,{"query":e.target.value})}
        />
      )
    } else {
      value_field = (
        <FormItem
          labelText={'Value'}
          name={'value'}
          value={data.query}
          inline={true}
          errors={errors}
          onChange={(e) => creator.updateFilter(data,{"query":e.target.value})}
        />
      )
    }
    
    return (
      <div className="fk-filter fk-guidfilter" key={data.id}>
        <div className="filter-body">
          <FormSelectItem
            labelText={'Condition'}
            name={'condition'}
            value={data.operation}
            inline={true}
            firstOptionText={'Select'}
            options={operations}
            errors={errors}
            onChange={(e) => creator.updateFilter(data,{"operation":e.target.value})}
          />
          {value_field}
          <div className="bt-icon medium"  onClick={() => creator.deleteFilter(data)}>
            <RemoveIcon />
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

  render() {
    const operations = ["GT","LT","EQ","notch"].map((o,i) => { return { value: o, text: o } })
    const { creator, data, errors } = this.props
    
    return (
      <div className="fk-filter fk-guidfilter">
        <div className="filter-body">
          <FormSelectItem
            labelText={'Condition'}
            name={'condition'}
            value={data.operation}
            inline={true}
            firstOptionText={'Select'}
            options={operations}
            errors={errors}
            onChange={(e) => creator.updateFilter(data,{"operation":e.target.value})}
          />
          <FormItem
            labelText={'Value'}
            name={'value'}
            type={'number'}
            value={data.query}
            inline={true}
            errors={errors}
            onChange={(e) => creator.updateFilter(data,{"query":e.target.value})}
          />
          <div className="bt-icon medium"  onClick={() => creator.deleteFilter(data)}>
            <RemoveIcon />
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

  render() {
    const operations = ["before","after","within"].map((o,i) => { return { value: o, text: o } })
    const { data, creator, errors } = this.props
    
    return (
      <div className="fk-filter fk-guidfilter">
        <div className="filter-body">
          <FormSelectItem
            labelText={'Condition'}
            name={'condition'}
            value={data.operation}
            inline={true}
            firstOptionText={'Select'}
            options={operations}
            errors={errors}
            onChange={(e) => creator.updateFilter(data,{"operation":e.target.value})}
          />
          <FormItem
            labelText={data.operation === 'within' ? 'Start Date' : 'Date'}
            name={'date'}
            type={'date'}
            value={data.date}
            inline={true}
            errors={errors}
            onChange={(e) => creator.updateFilter(data,{"date":e.target.value})}
          />
          { data.operation === 'within' &&
            <span> and
              <FormItem
                labelText={'End Date'}
                name={'date'}
                type={'date'}
                value={data.within}
                inline={true}
                errors={errors}
                onChange={(e) => creator.updateFilter(data,{"date":e.target.value})}
              />
            </span>
          }
          <div className="bt-icon medium"  onClick={() => creator.deleteFilter(data)}>
            <RemoveIcon />
          </div>
        </div>
      </div>
    )
  }
}
