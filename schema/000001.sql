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
	role varchar(100) NOT NULL,
	PRIMARY KEY (team_id, user_id)
);

-- input

CREATE TABLE fieldkit.input (
	id serial PRIMARY KEY,
	expedition_id integer REFERENCES fieldkit.expedition (id) NOT NULL,
	active boolean NOT NULL DEFAULT false
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
