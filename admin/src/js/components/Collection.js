/* @flow */

import type {Collection, Attr, StringAttr, StringFilter, NumAttr, NumFilter, DateAttr, DateFilter} from '../types/CollectionTypes'

export function stringifyOptions(attr: Attr): string {
  if(attr.type === "string"){
    return attr.options.join(", ")
  } else if (attr.type === "num"){
    return attr.options.map(n => n.toString()).join(", ")
  } else if (attr.type === "date"){
    return attr.options.map(n => n.toString()).join(", ")
  }

  return ""
}

export function cloneCollection(c: Collection): Collection {
  return {
    name: c.name,
    id: c.id,
    filters: c.filters.slice(0),
    guid_filters: c.guid_filters.slice(0),
    string_filters: c.string_filters.slice(0),
    num_filters: c.num_filters.slice(0),
    geo_filters: c.geo_filters.slice(0),
    date_filters: c.date_filters.slice(0),
    mods: c.mods.slice(0),
    string_mods: c.string_mods.slice(0),
    num_mods: c.num_mods.slice(0),
    geo_mods: c.geo_mods.slice(0)
  }
}

function getNextId(c: Collection): number{
  const max = Math.max(...c.filters);
  return max > -1 ? max + 1 : 0;
}

export function emptyStringFilter(c: Collection, attr: StringAttr): StringFilter {
  let next_id = getNextId(c);
  return {
    id: next_id,
    attribute: attr.name,
    operation: "contains",
    query: "",
    options: attr.options,
    type: "string",
    target: attr.target
  }
}

export function emptyNumFilter(c: Collection, attr: NumAttr): NumFilter {
  let next_id = getNextId(c);
  return {
    id: next_id,
    attribute: attr.name,
    operation: "GT",
    query: 3,
    type: "num",
    target: attr.target
  }
}

export function emptyDateFilter(c: Collection, attr: DateAttr): DateFilter {
  let next_id = getNextId(c);
  let f = {
    id: next_id,
    attribute: attr.name,
    operation: "before",
    date: 0,
    within: 0,
    type: "date",
    target: attr.target
  }
  return f
}


