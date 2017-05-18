/* @flow */

import {Location} from './proto/flow'
import type {GroupingOperation, SelectionOperation, Viz} from '../../admin/src/js/types/VizTypes'
import type {Attr} from '../../admin/src/js/types/CollectionTypes'
import R from 'ramda'
require('es6-promise').polyfill();
import 'isomorphic-fetch';

export type Stream = Array<Object>
export type GroupStream = Array<Array<Object>>
type GroupingFn = (Stream[]) => GroupStream
type GroupingFactory = (GroupingOperation) => GroupingFn
type StreamTransformer = (GroupStream) => Stream
type SelectionFn<A> = (Object[]) => A
type SelectionFactory<A> = (SelectionOperation) => SelectionFn<A>
type TxData = {"data": Array<Object>, "output": Object}
type SelectionTx = (TxData) => TxData

export const equalGrouping: GroupingFactory = (grouping) => {
  // TODO: ADD SOME SORTING
  return (streams) => {
    let dict = {}
    streams.forEach((s) => {
      s.forEach((d) => {
        let value = d[grouping.source_attribute]
        if(dict[value]){
          dict[value].push(d)
        } else {
          dict[value] = [d]
        }
      })
    })
    let output = [] 
    Object.keys(dict).forEach(k => output.push(dict[k]))
    return output
  }
}

export const withinGrouping: GroupingFactory = (grouping) => {
  return (streams) => {
    return []
  }
}

export const peakGrouping: GroupingFactory = (grouping) => {
  return (streams) => {
    return []
  }
}

export const groupingFactories = {
  "equal": equalGrouping,
  "within": withinGrouping,
  "peak": peakGrouping
}

export function getGroupingFn(grouping: GroupingOperation): GroupingFn{
  return groupingFactories[grouping.operation](grouping) 
}

export function avgSelection(s: SelectionOperation): SelectionFn<number>{
  return (group) => {
    const vals = group.map(d => d[s.source_attribute])
                      .filter(d => typeof(d) === "number")
    return vals.reduce((m,i) => m + i ,0) / vals.length
  }
}

export function countSelection(s: SelectionOperation): SelectionFn<number>{
  return (group) => {
    return group.length
  }
}

export function simpleStringSelection(s: SelectionOperation): SelectionFn<string>{
  return (group) => {
    const first_string_datum = group.find(d => d[s.source_attribute] && typeof(d[s.source_attribute]) == "string")
    return first_string_datum ? first_string_datum[s.source_attribute] : ""
  }
}

export function simpleNumSelection(s: SelectionOperation): SelectionFn<number>{
  return (group) => {
    const first_num_datum = group.find(d => d[s.source_attribute] && typeof(d[s.source_attribute]) == "number")
    return first_num_datum ? first_num_datum[s.source_attribute] : 0
  }
}

export function rawSelection(s: SelectionOperation): SelectionFn<Array<Object>>{
  return (group) => {
    return group
  }
}

export const selectionFactories = {
  "avg": avgSelection,
  "simple_string": simpleStringSelection,
  "simple_num": simpleNumSelection,
  "raw": rawSelection,
  "max": countSelection,
  "min": countSelection,
  "median": countSelection,
  "first": countSelection,
  "last": countSelection,
  "sum": countSelection,
  "count": countSelection
}

export function makeSelection(selection: SelectionOperation): SelectionTx{
  const selection_factory = selectionFactories[selection.operation]
  const selection_fn:SelectionFn<*> = selection_factory(selection)

  return ({data, output}) => {
    output[selection.value_name] = selection_fn(data)
    return {data, output}
  }
}

export function getStreamTransformer(selections: SelectionOperation[]): StreamTransformer{
  return (stream) => {
    const select = selections.reduce((m,s) => {
      return R.compose(makeSelection(s),m) 
    },R.identity)

    return stream.map((group) => {
      const new_datum = select({data: group,output: {}})
      return new_datum.output
    })
  }
}

export function transform(streams: Stream[], grouping: GroupingOperation, selections: SelectionOperation[]): Stream{
  const grouping_fn = getGroupingFn(grouping)
  const selection_fn = getStreamTransformer(selections)
  const new_stream = selection_fn(grouping_fn(streams))
  return new_stream
}

type Registration = {
  source_collections: string[],
  callback: (Stream) => void,
  processor: (Stream[]) => Stream,
  id: number
}

export class Broker {
  started: boolean;
  registrations: Registration[];
  api_root: string;
  cache: {[number]: Stream}

  constructor(api_root: string){
    this.registrations = [];
    this.started = false;
    this.api_root = api_root;
    this.cache = {}
    return this;
  }

  overwriteRegCache(id: number, data: Stream){
    this.cache[id] = data
  }

  getNewRegId(): number{
    if(this.registrations.length === 0){
      return 0
    } else {
      return Math.max(...this.registrations.map(r => r.id)) + 1
    }
  }

  registerViz(viz: Viz, callback: (Array<Object>) => void){
    this.register(viz.source_collections, viz.grouping_operation, viz.selection_operations, callback)
  }

  register(source_collections: string[], grouping: GroupingOperation, selections: SelectionOperation[], callback: (Array<Object>) => void){
    if(this.started){
      let error = new Error ("You cannot register with a broker once it has been started")
      throw error
    }
    const grouping_fn = getGroupingFn(grouping)
    const selection_fn = getStreamTransformer(selections)
    const id = this.getNewRegId(this.registrations)
    const new_registration = {
      source_collections: source_collections,
      processor: R.compose(selection_fn,grouping_fn),
      id,
      callback
    }

    this.registrations.push(new_registration)
    return new_registration
  }

  getReq(collections: string[]): [string,Object]{
    const url = ""
    const req = {}
    return [url,req]
  }

  request(){
    this.registrations.forEach(({source_collections,processor,callback,id}) => {
      const [url,req] = this.getReq(source_collections)
      fetch(url,req)
        .then((response) => {
          if (response.status >= 200 && response.status < 300) {
            return response
          } else {
            var error = new Error(response.statusText)
            throw error
          }
        })
        .then(response => response.json())
        .then(data => this.handleResponse(id,data,processor,callback))
    })  
  }

  handleResponse(reg_id: number, data: {streams: Stream[]}, processor: (Stream[]) => Stream, callback: (Stream) => void ) {
    let new_data = processor(data.streams)
    this.overwriteRegCache(reg_id, new_data) 
    callback(new_data)
    return new_data
  }
  
  start(){
    if(this.started){
      throw "You cannot start an already-running broker."
    } 
    this.started = true
  }
}
