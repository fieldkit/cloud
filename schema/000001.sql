DROP SCHEMA IF EXISTS fieldkit CASCADE;
CREATE SCHEMA fieldkit;

CREATE TABLE fieldkit.user (
	id serial PRIMARY KEY,
	username varchar(40) NOT NULL,
	email varchar(254) NOT NULL,
	password bytea NOT NULL,
	valid bool NOT NULL DEFAULT false


	-- login_attempts smallint NOT NULL DEFAULT 0,
	-- login_locked_until timestamp
);

CREATE UNIQUE INDEX ON fieldkit.user (username);
CREATE UNIQUE INDEX ON fieldkit.user (email);

CREATE TABLE fieldkit.validation_token (
	token bytea PRIMARY KEY,
	user_id integer REFERENCES fieldkit.user (id) NOT NULL,
	expires timestamp NOT NULL
);

CREATE TABLE fieldkit.refresh_token (
	token bytea PRIMARY KEY,
	user_id integer REFERENCES fieldkit.user (id) NOT NULL,
	expires timestamp NOT NULL
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
	expedition_id integer REFERENCES fieldkit.expedition (id) NOT NULL,
	name varchar(100) NOT NULL,
	slug varchar(100) NOT NULL,
	description text NOT NULL DEFAULT ''
);

CREATE UNIQUE INDEX ON fieldkit.team (expedition_id, slug);

CREATE TABLE fieldkit.team_user (
	team_id integer REFERENCES fieldkit.team (id) NOT NULL,
	user_id integer REFERENCES fieldkit.user (id) NOT NULL,
	PRIMARY KEY (team_id, user_id)
);

-- team

CREATE TABLE fieldkit.input (
	id serial PRIMARY KEY,
	name varchar(100) NOT NULL
);

CREATE TABLE fieldkit.input_expedition (
	input_id integer REFERENCES fieldkit.input (id) NOT NULL,
	expedition_id integer REFERENCES fieldkit.expedition (id) NOT NULL,
	PRIMARY KEY (input_id, expedition_id)
);

-- schema

CREATE TABLE fieldkit.schema (
	id serial PRIMARY KEY,
	schema jsonb NOT NULL
);

