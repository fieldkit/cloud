-- ingestion

CREATE TABLE fieldkit.ingestion (
  id serial PRIMARY KEY,
  time timestamp NOT NULL,
  upload_id varchar(64) NOT NULL,
  user_id integer NOT NULL,
  device_id bytea NOT NULL,
  generation bytea NOT NULL,
  type varchar NOT NULL,
  size integer NOT NULL,
  url varchar NOT NULL,
  blocks int8range NOT NULL,
  flags integer[] NOT NULL DEFAULT '{}',
  attempted timestamp,
  completed timestamp,
  errors boolean
);

CREATE INDEX ON fieldkit.ingestion (time, user_id);

CREATE TABLE fieldkit.provision (
  id serial PRIMARY KEY,
  created timestamp NOT NULL,
  updated timestamp NOT NULL,
  generation bytea NOT NULL,
  device_id bytea NOT NULL
);

CREATE UNIQUE INDEX ON fieldkit.provision (device_id, generation);

CREATE TABLE fieldkit.meta_record (
  id serial PRIMARY KEY,
  provision_id integer NOT NULL REFERENCES fieldkit.provision(id),
  time timestamp NOT NULL,
  number integer NOT NULL,
  raw json NOT NULL
);

CREATE INDEX ON fieldkit.meta_record (time, provision_id, number);

CREATE UNIQUE INDEX ON fieldkit.meta_record (provision_id, number);

CREATE TABLE fieldkit.data_record (
  id serial PRIMARY KEY,
  provision_id integer NOT NULL REFERENCES fieldkit.provision(id),
  time timestamp NOT NULL,
  number integer NOT NULL,
  meta integer NOT NULL REFERENCES fieldkit.meta_record(id),
  location geometry(POINT, 4326),
  raw json NOT NULL
);

CREATE INDEX ON fieldkit.data_record (time, provision_id, number);

CREATE UNIQUE INDEX ON fieldkit.data_record (provision_id, number);

CREATE AGGREGATE range_merge(anyrange) (
  sfunc = range_merge,
  stype = anyrange
);
