DROP TABLE fieldkit.module_sensor;
DROP TABLE fieldkit.station_module;

CREATE TABLE fieldkit.station_configuration (
	id serial PRIMARY KEY,
	provision_id INTEGER NOT NULL REFERENCES fieldkit.provision(id),
	meta_record_id INTEGER REFERENCES fieldkit.meta_record(id),
	source_id INTEGER,
    updated_at TIMESTAMP NOT NULL
);

CREATE UNIQUE INDEX ON fieldkit.station_configuration(provision_id, meta_record_id);
CREATE UNIQUE INDEX ON fieldkit.station_configuration(provision_id, source_id);

CREATE TABLE fieldkit.station_module (
	id serial PRIMARY KEY,
	configuration_id INTEGER REFERENCES fieldkit.station_configuration(id) NOT NULL,
	hardware_id BYTEA NOT NULL,
	module_index INTEGER NOT NULL,
	position INTEGER NOT NULL,
	flags INTEGER NOT NULL,
	manufacturer INTEGER NOT NULL,
	kind INTEGER NOT NULL,
	version INTEGER NOT NULL,
	name VARCHAR NOT NULL
);

CREATE UNIQUE INDEX ON fieldkit.station_module(configuration_id, hardware_id);

CREATE TABLE fieldkit.module_sensor (
	id serial PRIMARY KEY,
	module_id INTEGER NOT NULL REFERENCES fieldkit.station_module(id),
	configuration_id INTEGER REFERENCES fieldkit.station_configuration(id),
	sensor_index INTEGER NOT NULL,
	unit_of_measure VARCHAR NOT NULL,
	name VARCHAR NOT NULL,
	reading_last FLOAT,
	reading_time TIMESTAMP
);

CREATE UNIQUE INDEX ON fieldkit.module_sensor(module_id, sensor_index);

CREATE TABLE fieldkit.visible_configuration (
	station_id INTEGER NOT NULL REFERENCES fieldkit.station(id),
	configuration_id INTEGER NOT NULL REFERENCES fieldkit.station_configuration(id),
	PRIMARY KEY (station_id)
);
