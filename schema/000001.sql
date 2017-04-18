CREATE EXTENSION IF NOT EXISTS postgis;

DROP SCHEMA IF EXISTS fieldkit CASCADE;
CREATE SCHEMA fieldkit;

CREATE TABLE fieldkit.user (
	id serial PRIMARY KEY,
	name varchar(256) NOT NULL,
	username varchar(40) NOT NULL UNIQUE,
	email varchar(254) NOT NULL UNIQUE,
	password bytea NOT NULL,
	valid bool NOT NULL DEFAULT false,
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

-- input

CREATE TABLE fieldkit.input (
	id serial PRIMARY KEY,
	expedition_id integer REFERENCES fieldkit.expedition (id) NOT NULL,
	name varchar(256) NOT NULL,
	team_id int REFERENCES fieldkit.team (id), 
	user_id int REFERENCES fieldkit.user (id),
	active boolean NOT NULL DEFAULT false
);

CREATE TABLE fieldkit.input_token (
	id serial PRIMARY KEY,
	token bytea NOT NULL UNIQUE,
	expedition_id integer REFERENCES fieldkit.expedition (id) NOT NULL
);

-- twitter

CREATE TABLE fieldkit.twitter_oauth (
	input_id int REFERENCES fieldkit.input (id) ON DELETE CASCADE PRIMARY KEY,
	request_token varchar NOT NULL UNIQUE,
	request_secret varchar NOT NULL
);

CREATE TABLE fieldkit.twitter_account (
	id bigint PRIMARY KEY,
	screen_name varchar(15) NOT NULL,
	access_token varchar NOT NULL,
	access_secret varchar NOT NULL
);

CREATE TABLE fieldkit.input_twitter_account (
	input_id int REFERENCES fieldkit.input (id) ON DELETE CASCADE PRIMARY KEY,
	twitter_account_id bigint REFERENCES fieldkit.twitter_account (id) ON DELETE CASCADE NOT NULL
);

-- schema

CREATE TABLE fieldkit.schema (
	id serial PRIMARY KEY,
	project_id integer REFERENCES fieldkit.project (id),
	json_schema jsonb NOT NULL
);

CREATE UNIQUE INDEX ON fieldkit.schema ((json_schema->'id'));

-- fieldkit

CREATE TYPE fieldkit_binary_field AS ENUM ('varint', 'uvarint', 'float32', 'float64');

CREATE TABLE fieldkit.fieldkit_binary (
	input_id int REFERENCES fieldkit.input (id) ON DELETE CASCADE PRIMARY KEY,
	schema_id int REFERENCES fieldkit.schema (id) ON DELETE CASCADE NOT NULL,
	id smallint NOT NULL,
	fields fieldkit_binary_field[] NOT NULL,
	mapper jsonb NOT NULL,
	longitude varchar,
	latitude varchar,
	UNIQUE (input_id, id)
);

CREATE TABLE fieldkit.input_fieldkit (
	input_id int REFERENCES fieldkit.input (id) ON DELETE CASCADE PRIMARY KEY
);

-- documents

CREATE TABLE fieldkit.document (
	id bigserial PRIMARY KEY,
	input_id int REFERENCES fieldkit.input (id) NOT NULL,
	schema_id int REFERENCES fieldkit.schema (id) NOT NULL,
	team_id int REFERENCES fieldkit.team (id), 
	user_id int REFERENCES fieldkit.user (id),
	timestamp timestamp NOT NULL,
	location geometry(POINT, 4326) NOT NULL,
	data jsonb NOT NULL
);

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
