CREATE TABLE fieldkit.station_activity (
	id serial PRIMARY KEY,
    created_at timestamp NOT NULL,
	station_id INTEGER REFERENCES fieldkit.station(id) NOT NULL
);

CREATE INDEX ON fieldkit.station_activity (station_id, created_at);

CREATE TABLE fieldkit.station_deployed (
	deployed_at timestamp NOT NULL,
	location geometry(POINT, 4326) NOT NULL
) INHERITS (fieldkit.station_activity);

CREATE TABLE fieldkit.station_ingestion (
	uploader_id INTEGER REFERENCES fieldkit.user(id) NOT NULL,
	data_ingestion_id INTEGER REFERENCES fieldkit.ingestion(id) NOT NULL,
	data_records INTEGER NOT NULL,
	errors BOOLEAN NOT NULL
) INHERITS (fieldkit.station_activity);

CREATE TABLE fieldkit.project_follower (
	id serial PRIMARY KEY,
    created_at timestamp NOT NULL DEFAULT now(),
	project_id INTEGER REFERENCES fieldkit.project(id) NOT NULL,
	follower_id INTEGER REFERENCES fieldkit.user(id) NOT NULL
);

CREATE UNIQUE INDEX ON fieldkit.project_follower (project_id, follower_id);
CREATE        INDEX ON fieldkit.project_follower (created_at);

CREATE TABLE fieldkit.project_activity (
	id serial PRIMARY KEY,
    created_at timestamp NOT NULL,
	project_id INTEGER REFERENCES fieldkit.project(id) NOT NULL
);

CREATE INDEX ON fieldkit.project_activity (project_id, created_at);

CREATE TABLE fieldkit.project_update (
	author_id INTEGER REFERENCES fieldkit.user(id) NOT NULL,
    body text
) INHERITS (fieldkit.project_activity);

/* PG is unable to do this kind of FK to a base table, for now and
** that's really annoying. So, instead I'm going to query stations
** directly and get rid of this level of indirection. */

/*
CREATE TABLE fieldkit.project_station_activity (
	station_activity_id INTEGER REFERENCES fieldkit.station_activity(id) NOT NULL
) INHERITS (fieldkit.project_activity);
*/

CREATE OR REPLACE VIEW fieldkit.project_and_station_activity AS
    SELECT sa.created_at, ps.project_id, sa.station_id, sa.id AS station_activity_id, NULL AS project_activity_id
	  FROM fieldkit.project_station AS ps
      JOIN fieldkit.station_activity AS sa ON (ps.station_id = sa.station_id)
    UNION
	SELECT pa.created_at, pa.project_id, NULL, NULL, pa.id AS project_activity_id
	  FROM fieldkit.project_activity AS pa;


/**
 * One of these per ingestion, this makes things idempotent.
 */
CREATE UNIQUE INDEX ON fieldkit.station_ingestion(data_ingestion_id);
