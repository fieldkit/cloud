ALTER TABLE fieldkit.station ADD battery FLOAT;
ALTER TABLE fieldkit.station ADD recording_started_at TIMESTAMP;
ALTER TABLE fieldkit.station ADD memory_used INTEGER;
ALTER TABLE fieldkit.station ADD memory_available INTEGER;
ALTER TABLE fieldkit.station ADD firmware_number INTEGER;
ALTER TABLE fieldkit.station ADD firmware_time INTEGER;

DROP TABLE fieldkit.station_reading;
DROP TABLE fieldkit.station_module;

CREATE TABLE fieldkit.station_module (
	id serial PRIMARY KEY,
	provision_id INTEGER NOT NULL REFERENCES fieldkit.provision(id),
	meta_record_id integer REFERENCES fieldkit.meta_record(id),
	position INTEGER NOT NULL,
	hardware_id bytea NOT NULL,
	manufacturer INTEGER NOT NULL,
	kind INTEGER NOT NULL,
	version INTEGER NOT NULL,
	name varchar NOT NULL
);

CREATE UNIQUE INDEX ON fieldkit.station_module(hardware_id);
CREATE        INDEX ON fieldkit.station_module(provision_id);

CREATE TABLE fieldkit.module_sensor (
	id serial PRIMARY KEY,
	module_id INTEGER NOT NULL REFERENCES fieldkit.station_module(id),
	position INTEGER NOT NULL,
	unit_of_measure varchar NOT NULL,
	name varchar NOT NULL,
	reading float
);

CREATE UNIQUE INDEX ON fieldkit.module_sensor(module_id, position);
CREATE        INDEX ON fieldkit.module_sensor(module_id);
