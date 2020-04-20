CREATE TABLE fieldkit.station_module (
	id serial PRIMARY KEY,
	station_id INTEGER NOT NULL REFERENCES fieldkit.station(id),
	meta_record_id integer REFERENCES fieldkit.meta_record(id),
	generation bytea NOT NULL
);

CREATE INDEX ON fieldkit.station_module (station_id);

CREATE TABLE fieldkit.station_reading (
	id serial PRIMARY KEY,
	station_id INTEGER NOT NULL REFERENCES fieldkit.station(id),
	module_id INTEGER NOT NULL REFERENCES fieldkit.station_module(id),
	key TEXT NOT NULL,
	value FLOAT NOT NULL
);

CREATE INDEX ON fieldkit.station_reading (station_id, module_id);

CREATE INDEX ON fieldkit.station_reading (key);
