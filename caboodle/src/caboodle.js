/* @flow */

import type {GroupingOperation, SelectionOperation} from '../../admin/src/js/types/VizTypes'
import type {Attr} from '../../admin/src/js/types/CollectionTypes'
import R from 'ramda'

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

export const selectionFactories = {
  "avg": avgSelection,
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
  const selection_fn:SelectionFn<number> = selection_factory(selection)

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
