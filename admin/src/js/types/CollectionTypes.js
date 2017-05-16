export type Target = "expedition" | "binding" | "doctype" | "attribute"

export type StringAttr = {
  name: string;
  options: string[];
  type: "string";
  target: Target;
}

export type NumAttr = {
  name: string;
  options: number[];
  type: "num";
  target: Target;
}

export type DateAttr = {
  name: string;
  options: number[];
  type: "date";
  target: Target;
}

export type Attr = StringAttr | NumAttr | DateAttr

export type Expedition = StringAttr

export type Binding = StringAttr

export type Doctype = StringAttr

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
  target: Target;
}

export type NumFilter = {
  id: number;
  attribute: string;
  operation: "GT" | "LT" | "EQ";
  query: number;
  type: "num";
  target: Target;
}

export type GeoFilter = {
  id: number;
  attribute: string;
  operation: "within" | "not within";
  query: Object;
  type: "geo";
  target: Target;
}

export type DateFilter = {
  id: number;
  attribute: string;
  operation: "before" | "after" | "within";
  date: number;
  within: number;
  type: "date";
  target: Target;
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

export type ProjectData = {
    expeditions: Expedition,
    bindings: Binding[],
    doctypes: Doctype[],
    attributes: Attr[] 
}
