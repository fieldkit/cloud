/* @flow */

export type Collection = {
    name: string;
    id: string;
    filters: number[];
    guid_filters: GuidFilter[];      
    string_filters: StringFilter[];
	num_filters: NumFilter[];
	geo_filters: GeoFilter[];
	date_filters: DateFilter[];
	mods: number[];
	string_mods: StringMod[];
	num_mods: NumMod[];
	geo_mods: GeoMod[];
}

export type StringAttr = {
    name: string;
    options: string[];
    type: "string";
}

export type NumAttr = {
    name: string;
    options: number[];
    type: "num";
}

export type DateAttr = {
    name: string;
    options: number[];
    type: "date";
}

export type Attr = StringAttr | NumAttr | DateAttr

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
        type: "string"
    }
}

export function emptyNumFilter(c: Collection, attr: NumAttr): NumFilter {
    let next_id = getNextId(c);
    return {
        id: next_id,
        attribute: attr.name,
        operation: "GT",
        query: 3,
        type: "num"
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
        type: "date"
    }
    return f
}

export type GuidFilter = {
    id: number;
    attribute: string;
    operation: "contains";
    query: number[];
    type: "guid";
}

export type StringFilter = {
    id: number;
    attribute: string;
    operation: "contains" | "does not contain" | "matches" | "exists";
    query: string;
    options: string[];
    type: "string";
}

export type NumFilter = {
    id: number;
    attribute: string;
    operation: "GT" | "LT" | "EQ";
    query: number;
    type: "num";
}

export type GeoFilter = {
    id: number;
    attribute: string;
    operation: "within" | "not within";
    query: Object;
    type: "geo";
}

export type DateFilter = {
    id: number;
    attribute: string;
    operation: "before" | "after" | "within";
    date: number;
    within: number;
    type: "date";
}

export type Filter = NumFilter | StringFilter | DateFilter
export type FilterFn = (filter: StringFilter, update: $Shape<StringFilter>) => void |
                       (filter: NumFilter, update: $Shape<NumFilter>) => void |
                       (filter: DateFilter, update: $Shape<DateFilter>) => void

export type StringMod = {
    id: number;
    attribute: string;
    operation: "gsub";
    query: string;
    type: "string";
}
export type NumMod = {
    id: number;
    attribute: string;
    operation: "round";
    type: "num";
}
export type GeoMod = {
    id: number;
    attribute: string;
    operation: "copy from" | "jitter";
    source_collection: string;
    filters: number[];
    type: "geo";
}
