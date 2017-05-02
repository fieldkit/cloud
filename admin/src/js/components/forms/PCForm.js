/* @flow */

import React, { Component } from 'react'
import { FormContainer } from '../containers/FormContainer';
import type { APIErrors, APICollection, APINewCollection } from '../../api/types';
import { AddIcon } from '../icons/Icons'
import log from 'loglevel';

import type {Attr, Collection, Filter, FilterFn, StringFilter, DateFilter, NumFilter} from '../Collection';
import {cloneCollection, emptyStringFilter, emptyNumFilter, emptyDateFilter, stringifyOptions} from '../Collection';
import { StringFilterComponent, NumFilterComponent, DateFilterComponent } from '../forms/FiltersForm'



type Props = {
  collection?: ?APICollection,

  cancelText?: string;
  saveText?: ?string;
  onCancel?: () => void;
  onSave: (c: APINewCollection, collectionId?: number) => Promise<?APIErrors>;
};

function emptyCollection() : Collection {
  return {
    name: "",
    id: "",
    filters: [],
    guid_filters: [],
    string_filters: [],
    num_filters: [],
    geo_filters: [],
    date_filters: [],
    mods: [],
    string_mods: [],
    num_mods: [],
    geo_mods: []
  }
}

export class PCForm extends Component {
  attributes: { [string]: Attr };
  props: Props;
  state: {
    collection: Collection,
    errors: ?APIErrors
  }

  constructor(props: Props) {
    super(props);

    this.attributes = {
      "message": {
        name: "message",
        options: [],
        type: "string"
      },
      "user": {
        name: "user",
        options: ["@eric","@othererik","@gabriel"],
        type: "string"
      },
      "humidity": {
        name: "humidity",
        options: [],
        type: "num"
      },
      "created": {
        name: "created",
        options: [],
        type: "date"
      }
    }

    this.state = {
      collection: emptyCollection(),
      errors: null
    }
  }

  componentWillReceiveProps(nextProps: Props) {
  }

  addFilter(attr: string) {
    let collection = cloneCollection(this.state.collection);
    const attribute = this.attributes[attr];
    let new_filter;
    if(attribute.type === "string"){
      new_filter = emptyStringFilter(collection, attribute);
      collection.string_filters.push(new_filter)
    } else if (attribute.type === "num"){
      new_filter = emptyNumFilter(collection, attribute);
      collection.num_filters.push(new_filter)
    } else if (attribute.type === "date"){
      new_filter = emptyDateFilter(collection, attribute);
      collection.date_filters.push(new_filter)
    }
    if(new_filter){
      collection.filters.push(new_filter.id)
      this.setState({collection})
    }
  }

  updateFilter:FilterFn = (filter, update) => {
    let new_filter = Object.assign(filter,update)
    let collection = cloneCollection(this.state.collection)
    if(new_filter.type === "string"){
      collection.string_filters = collection.string_filters.filter(f => f.id !== new_filter.id)
      collection.string_filters.push(new_filter)
    } else if (filter.type === "date") {
      collection.date_filters = collection.date_filters.filter(f => f.id !== new_filter.id)
      collection.date_filters.push(new_filter)
    } else if (new_filter.type === "num") {
      collection.num_filters = collection.num_filters.filter(f => f.id !== new_filter.id)
      collection.num_filters.push(new_filter)
    }
    this.setState({collection})
  }

  deleteFilter(f: Filter){
    let collection = cloneCollection(this.state.collection)
    if(f.type === "string"){
      collection.string_filters = collection.string_filters.filter(cf => cf.id !== f.id)
    } else if(f.type === "num"){
      collection.num_filters = collection.num_filters.filter(cf => cf.id !== f.id)
    } else if(f.type === "date"){
      collection.date_filters = collection.date_filters.filter(cf => cf.id !== f.id)
    }
    collection.filters = collection.filters.filter(cf => cf !== f.id) 
    this.setState({collection})
  }

  save() {
    const { collection } = this.state;
    console.log(JSON.stringify(collection,null,' '));
    // TO-DO
    // const errors = await this.props.onSave(collection);
    // if (errors) {
    //   this.setState({ errors });
    // }
  }

  render() {
    const { errors } = this.state;
    let { collection } = this.state;

    let component_lookup : { [number]: Object } = {};

    collection.num_filters.forEach((f,i) => {
        component_lookup[f.id] = <NumFilterComponent creator={this} data={f} key={"n-" + i} errors={errors}/>
    })
    collection.string_filters.forEach((f,i) => {
        component_lookup[f.id] = <StringFilterComponent creator={this} data={f} key={"s-" + i} errors={errors}/>
    })
    collection.date_filters.forEach((f,i) => {
        component_lookup[f.id] = <DateFilterComponent creator={this} data={f} key={"d-" + i} errors={errors}/>
    })
    
    let filters_by_attribute = Object.keys(this.attributes).reduce((m,attr_name) => {
      const attr = this.attributes[attr_name]
      const { type,options } = attr
      let options_block = null
      if(options.length > 0) {
        options_block = <span className="filter-row-options">{stringifyOptions(attr)}</span> 
      }
      const attribute_row = (
        <div key={"filter-" + attr_name} className="accordion-row-header">
          <h4>{attr_name}</h4>            
          <button className="secondary" onClick={() => this.addFilter(attr_name)}>
            <div className="bt-icon medium">
              <AddIcon />
            </div>
            Add Filter
          </button>
        </div>
      )

      let buttons: Object[] = [attribute_row]
      collection.filters.forEach((f) => {
        const component = component_lookup[f]
        if(component){
          const {attribute} = component.props.data
          if(attr_name === attribute){
            buttons.push(component)
          }
        }
      })
      m[attr_name] = buttons
      return m
    },{})
     
    const components = Object.entries(filters_by_attribute).map(([attr,components]) => {
        return (
          <div className="accordion-row expanded" key={attr}>
            {components}
          </div>
        )
    })

    return (
      <FormContainer
        onSave={this.save.bind(this)}
        onCancel={this.props.onCancel}
        saveText={this.props.saveText}
        cancelText={this.props.cancelText}>
        {components}
      </FormContainer>
    )
  }
}
