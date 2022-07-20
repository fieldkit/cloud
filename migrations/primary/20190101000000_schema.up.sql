CREATE SCHEMA IF NOT EXISTS fieldkit;

CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE fieldkit.user (
	id serial PRIMARY KEY,
	name varchar(256) NOT NULL,
	username varchar(40) NOT NULL UNIQUE,
	email varchar(254) NOT NULL UNIQUE,
	password bytea NOT NULL,
	valid boolean NOT NULL DEFAULT false,
	bio varchar NOT NULL
);

CREATE TABLE fieldkit.validation_token (
	token bytea PRIMARY KEY,
	user_id integer REFERENCES fieldkit.user (id) ON DELETE CASCADE  NOT NULL,
	expires timestamp NOT NULL
);

CREATE TABLE fieldkit.refresh_token (
	token bytea PRIMARY KEY,
	user_id integer REFERENCES fieldkit.user (id) ON DELETE CASCADE NOT NULL,
	expires timestamp NOT NULL
);

CREATE TABLE fieldkit.invite_token (
	token bytea PRIMARY KEY
);

-- project

CREATE TABLE fieldkit.project (
	id serial PRIMARY KEY,
	name varchar(100) NOT NULL,
	slug varchar(100) NOT NULL,
	description text NOT NULL DEFAULT '',
	goal varchar(100) NOT NULL DEFAULT '',
	location varchar(100) NOT NULL DEFAULT '',
	tags varchar(100) NOT NULL DEFAULT '',
	start_time timestamp,
	end_time timestamp,
	private boolean NOT NULL DEFAULT false
);

CREATE UNIQUE INDEX ON fieldkit.project (slug);

CREATE TABLE fieldkit.project_user (
	project_id integer REFERENCES fieldkit.project (id) NOT NULL,
	user_id integer REFERENCES fieldkit.user (id) NOT NULL,
	PRIMARY KEY (project_id, user_id)
);

-- expedition

CREATE TABLE fieldkit.expedition (
	id serial PRIMARY KEY,
	project_id integer REFERENCES fieldkit.project (id) NOT NULL,
	name varchar(100) NOT NULL,
	slug varchar(100) NOT NULL,
	description text NOT NULL DEFAULT ''
);

CREATE UNIQUE INDEX ON fieldkit.expedition (project_id, slug);

-- team

CREATE TABLE fieldkit.team (
	id serial PRIMARY KEY,
	expedition_id integer REFERENCES fieldkit.expedition (id) ON DELETE CASCADE NOT NULL,
	name varchar(256) NOT NULL,
	slug varchar(100) NOT NULL,
	description text NOT NULL DEFAULT '',
	UNIQUE (expedition_id, slug)
);

CREATE TABLE fieldkit.team_user (
	team_id integer REFERENCES fieldkit.team (id) ON DELETE CASCADE NOT NULL,
	user_id integer REFERENCES fieldkit.user (id) ON DELETE CASCADE NOT NULL,
	role varchar(100) NOT NULL,
	PRIMARY KEY (team_id, user_id)
);

-- source

CREATE TABLE fieldkit.source (
	id serial PRIMARY KEY,
	expedition_id integer REFERENCES fieldkit.expedition (id) NOT NULL,
	name varchar(256) NOT NULL,
	team_id int REFERENCES fieldkit.team (id),
	user_id int REFERENCES fieldkit.user (id),
	active boolean NOT NULL DEFAULT false,
	visible boolean NOT NULL DEFAULT true
);

CREATE TABLE fieldkit.source_token (
	id serial PRIMARY KEY,
	token bytea NOT NULL UNIQUE,
	expedition_id integer REFERENCES fieldkit.expedition (id) NOT NULL
);

-- twitter

CREATE TABLE fieldkit.twitter_oauth (
	source_id int REFERENCES fieldkit.source (id) ON DELETE CASCADE PRIMARY KEY,
	request_token varchar NOT NULL UNIQUE,
	request_secret varchar NOT NULL
);

CREATE TABLE fieldkit.twitter_account (
	id bigint PRIMARY KEY,
	screen_name varchar(15) NOT NULL,
	access_token varchar NOT NULL,
	access_secret varchar NOT NULL
);

CREATE TABLE fieldkit.source_twitter_account (
	source_id int REFERENCES fieldkit.source (id) ON DELETE CASCADE PRIMARY KEY,
	twitter_account_id bigint REFERENCES fieldkit.twitter_account (id) ON DELETE CASCADE NOT NULL
);

-- schema

CREATE TABLE fieldkit.schema (
	id serial PRIMARY KEY,
	project_id integer REFERENCES fieldkit.project (id),
	json_schema jsonb NOT NULL
);

CREATE UNIQUE INDEX ON fieldkit.schema ((json_schema->'id'));

-- device

CREATE TABLE fieldkit.device (
    source_id integer REFERENCES fieldkit.source (id) ON DELETE CASCADE PRIMARY KEY,
    key varchar NOT NULL,
    token bytea NOT NULL
);

CREATE UNIQUE INDEX ON fieldkit.device (key);
CREATE UNIQUE INDEX ON fieldkit.device (token);

-- device_firmware

CREATE TABLE fieldkit.firmware (
  id serial PRIMARY KEY,
  time timestamp NOT NULL,
  module varchar(64) NOT NULL,
  profile varchar(64) NOT NULL,
  etag varchar(64) NOT NULL,
  url varchar NOT NULL,
  meta json NOT NULL
);

CREATE TABLE fieldkit.device_firmware (
  id serial PRIMARY KEY,
  device_id integer REFERENCES fieldkit.device (source_id) ON DELETE CASCADE,
  time timestamp NOT NULL,
  module varchar(64) NOT NULL,
  profile varchar(64) NOT NULL,
  etag varchar(64) NOT NULL,
  url varchar NOT NULL
);

-- device_stream

CREATE TABLE fieldkit.device_stream (
  id serial PRIMARY KEY,
  time timestamp NOT NULL,
  stream_id varchar(64) NOT NULL,
  firmware varchar(64) NOT NULL,
  device_id varchar(64) NOT NULL,
  size integer NOT NULL,
  file_id varchar(32) NOT NULL,
  url varchar NOT NULL,
  meta json NOT NULL,
  flags integer[]
);

CREATE UNIQUE INDEX ON fieldkit.device_stream (stream_id);

CREATE INDEX ON fieldkit.device_stream (device_id, time);

CREATE TABLE fieldkit.device_stream_location (
  id serial PRIMARY KEY,
  device_id varchar(64) NOT NULL,
  timestamp timestamp NOT NULL,
  location geometry(POINT, 4326) NOT NULL
);

CREATE INDEX ON fieldkit.device_stream_location (device_id, timestamp);

CREATE UNIQUE INDEX ON fieldkit.device_stream_location (device_id, timestamp, location);

CREATE TABLE fieldkit.device_notes (
  id serial PRIMARY KEY,
  device_id varchar(64) NOT NULL,
  time timestamp NOT NULL,
  name varchar,
  notes varchar
);

CREATE INDEX ON fieldkit.device_notes (device_id, time);

-- device_schema

CREATE TABLE fieldkit.device_schema (
    id serial PRIMARY KEY,
    device_id integer REFERENCES fieldkit.device (source_id) ON DELETE CASCADE,
    schema_id integer REFERENCES fieldkit.schema (id) ON DELETE CASCADE,
    key varchar
);

CREATE UNIQUE INDEX ON fieldkit.device_schema (device_id, schema_id);

-- device_location

CREATE TABLE fieldkit.device_location (
  id serial PRIMARY KEY,
  timestamp timestamp NOT NULL,
  device_id integer REFERENCES fieldkit.device (source_id) NOT NULL,
  location geometry(POINT, 4326) NOT NULL
);

CREATE INDEX ON fieldkit.device_location (device_id, timestamp);

-- data

CREATE TABLE fieldkit.countries (
  "gid" serial PRIMARY KEY,
  "fips" varchar(2),
  "iso2" varchar(2),
  "iso3" varchar(3),
  "un" int2,
  "name" varchar(50),
  "area" int4,
  "pop2005" int8,
  "region" int2,
  "subregion" int2,
  "lon" float8,
  "lat" float8
);

SELECT AddGeometryColumn('fieldkit', 'countries', 'geom', '0', 'MULTIPOLYGON', 2);

CREATE INDEX ON fieldkit.countries USING GIST (geom);

-- user

DO
$body$
BEGIN
   IF NOT EXISTS (SELECT * FROM pg_catalog.pg_user WHERE usename = 'server') THEN
      CREATE ROLE server LOGIN PASSWORD 'changeme';
   END IF;
END
$body$;

-- grants

GRANT USAGE ON SCHEMA fieldkit TO server;
GRANT ALL ON ALL TABLES IN SCHEMA fieldkit TO server;
GRANT ALL ON ALL SEQUENCES IN SCHEMA fieldkit TO server;
