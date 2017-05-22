/* @flow */

import {Location} from './proto/flow'
import type {GroupingOperation, SelectionOperation, Viz} from '../../admin/src/js/types/VizTypes'
import type {Attr} from '../../admin/src/js/types/CollectionTypes'
import R from 'ramda'
require('es6-promise').polyfill();
import 'isomorphic-fetch';

/** The type of a single collection of homogenous FieldKit data */
export type Stream = Array<Object>

/** The type of one or more streams the have had a grouping operation applied */
export type GroupStream = Array<Array<Object>>

/** The type of a function for grouping multiple streams into a single stream */
type GroupingFn = (Stream[]) => GroupStream

/** The type of a function for selecting a GroupingFn, given a GroupingOperation */
type GroupingFactory = (GroupingOperation) => GroupingFn

/** The type of a function that transforms a group stream into an output stream,
 *  to be consumed by a visualization
 */
type StreamTransformer = (GroupStream) => Stream

/** The type of function that picks out a single value from a group of documents */
type SelectionFn<A> = (Object[]) => A

/** The type of a function for selecting a SelectionFn, given a SelectionOperation */
type SelectionFactory<A> = (SelectionOperation) => SelectionFn<A>

/** The type of the accumulator for a stream transformation */
type TxData = {"data": Array<Object>, "output": Object}
type SelectionTx = (TxData) => TxData

/** The different types of valid simple selections */
export type SimpleSelectionType = "number" | "string" | "location"

/** A map for specifying the desired types of output for a map-generated selection*/
export type SelectionMap = {[string]: SimpleSelectionType}

/** Transforms a SelectionMap into an array of selections, in which the value key
 *  will match the source_key, and the output type will correspond to a mapping in the SelectionMap.
 *  
 *  This is to be used to easily generate simple selections, when not using the FieldKit admin to generate
 *  selections. Especially useful when using caboodle without the FieldKit frontend.
 *
 *  @param {SelectionMap} sm - The map used to specify the source_keys to select and the desired type for the selection on that key
 *  @returns {SelectionOperation[]} - The resulting selection operations, to be fed into a transform
 */

export function generateSelectionsFromMap(sm: SelectionMap): SelectionOperation[]{
  return Object.keys(sm).map((sk,i) => {
    let op;
    let type = sm[sk]

    if(type === "number"){
      op = "simple_num"
    } else if (type === "string"){
      op = "simple_string"
    } else if (type === "location"){
      op = "simple_location"
    } else {
      throw `Simple Selection type ${type} not recognized`
    }

    return {
      id: i,
      value_name: sk,
      source_attribute: sk,
      operation: op
    }
  })
}


/**
 * Transforms a generic stream of data into a GeoJSON object of Points
 * @param {Stream} s - The stream to be transformed
 * @param {string} location_key - The key of the object to be used as the coordinates for the Point 
 * @returns {Object} - A GeoJSON Feature Collection of Points Features
 */
export function streamToGeoJSON(s: Stream, location_key: string): Object{
  const features = s.filter(d => d[location_key] && d[location_key].hasOwnProperty("latitude"))
                    .map((d) => {
                      return {
                        "type":"Feature",
                        "properties":d,
                        "geometry":{
                          "type":"Point",
                          "coordinates":[
                            d[location_key].latitude,
                            d[location_key].longitude
                          ]
                        }
                      }
                    })
  return {
    "type": "FeatureCollection",
    "features": features
  } 
}

/**
 * Produces a function that performs an equal grouping operation on multiple streams
 * @param {GroupingOperation} grouping - The GroupingOperation that specifies the parameters for the equal grouping
 */
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

// TODO
export const withinGrouping: GroupingFactory = (grouping) => {
  return (streams) => {
    return []
  }
}

// TODO
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

/**
 * Translate a GroupingOperation into a Grouping Function
 * @param {GroupingOperation} grouping - The operation to transform
 * @result - The functional equivalent of the operating
 */
export function getGroupingFn(grouping: GroupingOperation): GroupingFn{
  return groupingFactories[grouping.operation](grouping) 
}

/** Generate a function that selects the average value in a group 
 *  @param {SelectionOperation} s - The operation that specifies which attribute we are averaging
 *  @result {SelectionFn<number>} - The function that will perform the averaging on a group
 */
export function avgSelection(s: SelectionOperation): SelectionFn<number>{
  return (group) => {
    const vals = group.map(d => d[s.source_attribute])
                      .filter(d => typeof(d) === "number")
    return vals.reduce((m,i) => m + i ,0) / vals.length
  }
}

/** Generate a function that counts the members in a group 
 *  @param {SelectionOperation} s - The selection operation
 *  @result {SelectionFn<number>} - The function that will count the members of a group
 */
export function countSelection(s: SelectionOperation): SelectionFn<number>{
  return (group) => {
    return group.length
  }
}

/** Generate a function that selects the first valid string value at a given attribute in a group
 *  @param {SelectionOperation} s - The operation that specifies which attribute we are selecting
 *  @result {SelectionFn<number>} - The function that will perform the selection
 */
export function simpleStringSelection(s: SelectionOperation): SelectionFn<string>{
  return (group) => {
    const first_string_datum = group.find(d => d[s.source_attribute] && typeof(d[s.source_attribute]) == "string" && d[s.source_attribute].length > 0)
    return first_string_datum ? first_string_datum[s.source_attribute] : ""
  }
}

/** Generate a function that selects the first valid number value at a given attribute in a group
 *  @param {SelectionOperation} s - The operation that specifies which attribute we are selecting
 *  @result {SelectionFn<number>} - The function that will perform the selection
 */
export function simpleNumSelection(s: SelectionOperation): SelectionFn<number>{
  return (group) => {
    const first_num_datum = group.find(d => d[s.source_attribute] && typeof(d[s.source_attribute]) == "number" && d[s.source_attribute] !== 0)
    return first_num_datum ? first_num_datum[s.source_attribute] : 0
  }
}

/** Generate a function that selects the first valid location value at a given attribute in a group
 *  @param {SelectionOperation} s - The operation that specifies which attribute we are selecting
 *  @result {SelectionFn<number>} - The function that will perform the selection
 */
export function simpleLocationSelection(s: SelectionOperation): SelectionFn<Location>{
  return (group) => {
    const first_loc_datum = group.find(d => d[s.source_attribute] && d[s.source_attribute].hasOwnProperty("latitude") && d[s.source_attribute].latitude !== 0)
    return first_loc_datum ? first_loc_datum[s.source_attribute] : new Location({latitude: 0, longitude: 0, altitude: 0})
  }
}


/** Generate a function that copies the entire group into a value of the selection
 *  @param {SelectionOperation} s - Ignored
 *  @result {SelectionFn<number>} - The function that will perform the selection
 */
export function wholeGroupSelection(s: SelectionOperation): SelectionFn<Array<Object>>{
  return (group) => {
    return group
  }
}

export const selectionFactories = {
  "avg": avgSelection,
  "simple_string": simpleStringSelection,
  "simple_num": simpleNumSelection,
  "simple_location": simpleLocationSelection,
  "whole_group": wholeGroupSelection,
  "max": countSelection,
  "min": countSelection,
  "median": countSelection,
  "first": countSelection,
  "last": countSelection,
  "sum": countSelection,
  "count": countSelection
}

/**
 * Converts a selection operation into a function that transforms an accumulator by that operation
 * @param {SelectionOperation} selection - The operation to use for the transformation
 * @result {SelectionTx} - A function for transforming an accumulator over a stream of group
 */
export function makeSelection(selection: SelectionOperation): SelectionTx{
  const selection_factory = selectionFactories[selection.operation]
  const selection_fn:SelectionFn<*> = selection_factory(selection)

  return ({data, output}) => {
    output[selection.value_name] = selection_fn(data)
    return {data, output}
  }
}

/**
 *  Convert multiple selection operations into a single function that transforms a stream of groups
 *  into a single stream of objects, which each have a single field and value corresponding to the result 
 *  of a selection operation
 *
 *  @params {SelectionOperation[]} selections - The selections to perform
 *  @result {StreamTransformer} - A function that transforms a grouped stream into a simple stream
 */
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

/**
 * Transforms multiple streams into a simple stream
 *
 * @param {Stream[]} streams - The input streams to transform
 * @param {GroupingOperation} grouping - The instructions for how to group the streams
 * @params {SelectionOperations[]} selections - The selections to output as fields of the objects in the final stream
 */
export function transform(streams: Stream[], grouping: GroupingOperation, selections: SelectionOperation[]): Stream{
  const grouping_fn = getGroupingFn(grouping)
  const selection_fn = getStreamTransformer(selections)
  const new_stream = selection_fn(grouping_fn(streams))
  return new_stream
}

/** The type of a cached subscription of a consumer with the broker */
type Registration = {
  source_collections: string[],
  callback: (Stream) => void,
  processor: (Stream[]) => Stream,
  id: number
}

/** A broker to manage multiple requests for streaming data. Controls live-updating and cacheing. 
 * @constructor
 */
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

  /** Get an new incremental id for a registration */
  getNewRegId(): number{
    if(this.registrations.length === 0){
      return 0
    } else {
      return Math.max(...this.registrations.map(r => r.id)) + 1
    }
  }

  /** A convenience method for registering Vizs with a broker */
  registerViz(viz: Viz, callback: (Array<Object>) => void){
    this.register(viz.source_collections, viz.grouping_operation, viz.selection_operations, callback)
  }

  /** Register a request for streaming data 
   *  @param {string[]} source_collections - The collections to request from the api
   *  @param {GroupingOperation} grouping - The grouping to apply to the streams
   *  @param {SelectionOperation[]} selections - The selections to reduce the groups into
   *  @callback {function} callback - The callback to pass the transformed stream into whenever data is updated
   */
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

  /** Request data for each registration */
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


