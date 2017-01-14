-- user

CREATE SCHEMA IF NOT EXISTS admin;

DROP TABLE IF EXISTS admin.user CASCADE;
CREATE TABLE admin.user (
	id bytea PRIMARY KEY,
	username varchar(80) NOT NULL,
	email varchar(80) NOT NULL,
	password bytea NOT NULL,
	first_name varchar(80) NOT NULL,
	last_name varchar(80) NOT NULL,
	valid bool NOT NULL DEFAULT false
);

CREATE UNIQUE INDEX ON admin.user (username);
CREATE UNIQUE INDEX ON admin.user (email);


-- user validation

DROP TABLE IF EXISTS admin.user_validation_token CASCADE;
CREATE TABLE admin.user_validation_token (
	id bytea PRIMARY KEY,
	user_id bytea REFERENCES admin.user (id) NOT NULL,
	expires timestamp NOT NULL
);

CREATE UNIQUE INDEX ON admin.user_validation_token (user_id);

-- invite

DROP TABLE IF EXISTS admin.invite CASCADE;
CREATE TABLE admin.invite (
	id bytea PRIMARY KEY,
	expires timestamp NOT NULL
);


-- project

DROP TABLE IF EXISTS admin.project CASCADE;
CREATE TABLE admin.project (
	id bytea PRIMARY KEY,
	name varchar(80) NOT NULL,
	slug varchar(80) NOT NULL
);

CREATE UNIQUE INDEX ON admin.project (slug);


-- project user

DROP TYPE IF EXISTS role_type CASCADE;
CREATE TYPE role_type AS ENUM ('member', 'administrator', 'owner');

DROP TABLE IF EXISTS admin.project_user CASCADE;
CREATE TABLE admin.project_user (
	project_id bytea REFERENCES admin.project (id) NOT NULL,
	user_id bytea REFERENCES admin.user (id) NOT NULL,
	role role_type NOT NULL
);

CREATE INDEX ON admin.project_user (project_id);
CREATE INDEX ON admin.project_user (user_id);
CREATE UNIQUE INDEX ON admin.project_user (project_id, user_id);
CREATE INDEX ON admin.project_user (project_id, role);


-- expedition

DROP TABLE IF EXISTS admin.expedition CASCADE;
CREATE TABLE admin.expedition (
	id bytea PRIMARY KEY,
	project_id bytea REFERENCES admin.project (id) NOT NULL,
	name varchar(80) NOT NULL,
	slug varchar(80) NOT NULL
);

CREATE INDEX ON admin.expedition (name);
CREATE UNIQUE INDEX ON admin.expedition (project_id, slug);


-- team

DROP TABLE IF EXISTS admin.team CASCADE;
CREATE TABLE admin.team (
	id bytea PRIMARY KEY,
	expedition_id bytea REFERENCES admin.expedition (id) NOT NULL,
	leader_user_id bytea REFERENCES admin.user (id),
	name varchar(80) NOT NULL,
	slug varchar(80) NOT NULL
);

CREATE INDEX ON admin.team (expedition_id);
CREATE INDEX ON admin.team (name);
CREATE UNIQUE INDEX ON admin.team (expedition_id, slug);

DROP TABLE IF EXISTS admin.team_user CASCADE;
CREATE TABLE admin.team_user (
	team_id bytea REFERENCES admin.team (id) NOT NULL,
	user_id bytea REFERENCES admin.user (id) NOT NULL
);

CREATE INDEX ON admin.team_user (team_id);
CREATE INDEX ON admin.team_user (user_id);
CREATE UNIQUE INDEX ON admin.team_user (team_id, user_id);


-- input

DROP TABLE IF EXISTS admin.input CASCADE;
CREATE TABLE admin.input (
	id bytea PRIMARY KEY,
	expedition_id bytea REFERENCES admin.expedition (id) NOT NULL
);

CREATE INDEX ON admin.input (expedition_id);


-- request

CREATE SCHEMA IF NOT EXISTS data;

DROP TYPE IF EXISTS format_type CASCADE;
CREATE TYPE format_type AS ENUM ('fieldkit', 'csv', 'json');

DROP TABLE IF EXISTS data.request CASCADE;
CREATE TABLE data.request (
	id bytea,
	input_id bytea REFERENCES admin.input (id) NOT NULL,
	format format_type NOT NULL,
	checksum bytea NOT NULL,
	data bytea NOT NULL
);

CREATE INDEX ON data.request (input_id);
CREATE UNIQUE INDEX ON data.request (id, input_id);
