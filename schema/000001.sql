CREATE EXTENSION IF NOT EXISTS postgis;

DROP SCHEMA IF EXISTS fieldkit CASCADE;
CREATE SCHEMA fieldkit;

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
	description text NOT NULL DEFAULT ''
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
	active boolean NOT NULL DEFAULT false
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

-- raw_message

CREATE TABLE fieldkit.raw_message (
    id serial PRIMARY KEY,
    time timestamp,
    origin_id varchar NOT NULL, -- SQS Id or some other provider generated identifier.
    data json NOT NULL
);

CREATE UNIQUE INDEX ON fieldkit.raw_message (origin_id);

-- device

CREATE TABLE fieldkit.device (
    source_id integer REFERENCES fieldkit.source (id) ON DELETE CASCADE PRIMARY KEY,
    key varchar NOT NULL,
    token bytea NOT NULL
);

CREATE UNIQUE INDEX ON fieldkit.device (key);

CREATE UNIQUE INDEX ON fieldkit.device (token);

CREATE TABLE fieldkit.device_schema (
    id serial PRIMARY KEY,
    device_id integer REFERENCES fieldkit.device (source_id) ON DELETE CASCADE,
    schema_id integer REFERENCES fieldkit.schema (id) ON DELETE CASCADE,
    key varchar
);

CREATE UNIQUE INDEX ON fieldkit.device_schema (device_id, schema_id);

CREATE TABLE fieldkit.device_location (
  id serial PRIMARY KEY,
  timestamp timestamp NOT NULL,
  device_id integer REFERENCES fieldkit.device (source_id) NOT NULL,
  location geometry(POINT, 4326) NOT NULL
);

CREATE INDEX ON fieldkit.device_location (device_id, timestamp);

-- records

CREATE TABLE fieldkit.record (
	id bigserial PRIMARY KEY,
	source_id int REFERENCES fieldkit.source (id) NOT NULL,
	schema_id int REFERENCES fieldkit.schema (id) NOT NULL,
	team_id int REFERENCES fieldkit.team (id),
	user_id int REFERENCES fieldkit.user (id),
	insertion timestamp NOT NULL DEFAULT now(),
	timestamp timestamp NOT NULL,
	location geometry(POINT, 4326) NOT NULL,
	visible boolean NOT NULL DEFAULT true,
	fixed boolean NOT NULL DEFAULT true,
	data jsonb NOT NULL
);

CREATE INDEX ON fieldkit.record (timestamp, source_id);
CREATE INDEX ON fieldkit.record USING GIST (location);

-- user

DO
$body$
BEGIN
   IF NOT EXISTS (
      SELECT *
      FROM   pg_catalog.pg_user
      WHERE  usename = 'server') THEN

      CREATE ROLE server LOGIN PASSWORD 'changeme';
   END IF;
END
$body$;

-- grants

GRANT USAGE ON SCHEMA fieldkit TO server;
GRANT ALL ON ALL TABLES IN SCHEMA fieldkit TO server;
GRANT ALL ON ALL SEQUENCES IN SCHEMA fieldkit TO server;
