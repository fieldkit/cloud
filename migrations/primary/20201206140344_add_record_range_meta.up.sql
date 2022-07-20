CREATE EXTENSION IF NOT EXISTS btree_gist;

CREATE TABLE fieldkit.record_range_meta (
	id serial PRIMARY KEY,
	station_id integer NOT NULL REFERENCES fieldkit.station(id),
	start_time timestamp NOT NULL,
	end_time timestamp NOT NULL,
	flags integer NOT NULL,

	EXCLUDE USING gist (station_id WITH =, tsrange(start_time, end_time) WITH &&)
);
