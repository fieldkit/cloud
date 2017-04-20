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

export function emptyStringFilter(c: Collection, attr: string, opts: string[]): StringFilter {
    let next_id = getNextId(c);
    return {
        id: next_id,
        attribute: attr,
        operation: "contains",
        query: "",
        options: opts
    }
}

export function emptyNumFilter(c: Collection, attr: string): NumFilter {
    let next_id = getNextId(c);
    return {
        id: next_id,
        attribute: attr,
        operation: "GT",
        query: 3,
    }
}

export function emptyDateFilter(c: Collection, attr: string): DateFilter {
    let next_id = getNextId(c);
    return {
        id: next_id,
        attribute: attr,
        operation: "before",
        date: 0,
        within: 0
    }
}

export type GuidFilter = {
    id: number;
    attribute: string;
    operation: "contains";
    query: number[];
}

export type StringFilter = {
    id: number;
    attribute: string;
    operation: "contains" | "does not contain" | "matches" | "exists";
    query: string;
    options: string[];
}

export type NumFilter = {
    id: number;
    attribute: string;
    operation: "GT" | "LT" | "EQ" | "notch";
    query: number;
}

export type GeoFilter = {
    id: number;
    attribute: string;
    operation: "within" | "not within";
    query: Object
}

export type DateFilter = {
    id: number;
    attribute: string;
    operation: "before" | "after" | "within";
    date: number;
    within: number;
}
export type StringMod = {
    id: number;
    attribute: string;
    operation: "gsub";
    query: string;
}
export type NumMod = {
    id: number;
    attribute: string;
    operation: "round";
}
export type GeoMod = {
    id: number;
    attribute: string;
    operation: "copy from" | "jitter";
    source_collection: string;
    filters: number[];
}
