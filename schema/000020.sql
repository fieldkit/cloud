-- ingestion

CREATE TABLE fieldkit.ingestion (
  id serial PRIMARY KEY,
  time timestamp NOT NULL,
  upload_id varchar(64) NOT NULL,
  user_id integer NOT NULL,
  device_id bytea NOT NULL,
  size integer NOT NULL,
  url varchar NOT NULL,
  blocks int8range NOT NULL,
  flags integer[] NOT NULL DEFAULT '{}',
  processing_attempted timestamp,
  processing_completed timestamp,
  type varchar,
  errors boolean
);

CREATE INDEX ON fieldkit.ingestion (time, user_id);

CREATE TABLE fieldkit.provision (
  id serial PRIMARY KEY,
  time timestamp NOT NULL,
  generation bytea NOT NULL,
  ingestion_id INTEGER NOT NULL,
  device_id bytea NOT NULL
);

CREATE INDEX ON fieldkit.provision (time, generation);

CREATE TABLE fieldkit.meta_record (
  id serial PRIMARY KEY,
  provision_id integer NOT NULL,
  time timestamp NOT NULL,
  number integer NOT NULL,
  raw json NOT NULL
);

CREATE INDEX ON fieldkit.meta_record (time, provision_id, number);

CREATE UNIQUE INDEX ON fieldkit.meta_record (provision_id, number);

CREATE TABLE fieldkit.data_record (
  id serial PRIMARY KEY,
  provision_id integer NOT NULL,
  time timestamp NOT NULL,
  number integer NOT NULL,
  location geometry(POINT, 4326),
  raw json NOT NULL
);

CREATE INDEX ON fieldkit.data_record (time, provision_id, number);

CREATE UNIQUE INDEX ON fieldkit.data_record (provision_id, number);
