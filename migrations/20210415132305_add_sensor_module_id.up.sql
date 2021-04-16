CREATE TABLE fieldkit.aggregated_bymod_10s (
	id SERIAL PRIMARY KEY,
	time TIMESTAMP NOT NULL,
	station_id integer NOT NULL REFERENCES fieldkit.station(id),
	sensor_id integer NOT NULL REFERENCES fieldkit.aggregated_sensor(id),
	module_id integer NOT NULL REFERENCES fieldkit.station_module(id),
	nsamples integer NOT NULL,
	location geometry(POINT, 4326),
	value FLOAT NOT NULL
);

CREATE TABLE fieldkit.aggregated_bymod_1m (
	id SERIAL PRIMARY KEY,
	time TIMESTAMP NOT NULL,
	station_id integer NOT NULL REFERENCES fieldkit.station(id),
	sensor_id integer NOT NULL REFERENCES fieldkit.aggregated_sensor(id),
	module_id integer NOT NULL REFERENCES fieldkit.station_module(id),
	nsamples integer NOT NULL,
	location geometry(POINT, 4326),
	value FLOAT NOT NULL
);

CREATE TABLE fieldkit.aggregated_bymod_10m (
	id SERIAL PRIMARY KEY,
	time TIMESTAMP NOT NULL,
	station_id integer NOT NULL REFERENCES fieldkit.station(id),
	sensor_id integer NOT NULL REFERENCES fieldkit.aggregated_sensor(id),
	module_id integer NOT NULL REFERENCES fieldkit.station_module(id),
	nsamples integer NOT NULL,
	location geometry(POINT, 4326),
	value FLOAT NOT NULL
);

CREATE TABLE fieldkit.aggregated_bymod_30m (
	id SERIAL PRIMARY KEY,
	time TIMESTAMP NOT NULL,
	station_id integer NOT NULL REFERENCES fieldkit.station(id),
	sensor_id integer NOT NULL REFERENCES fieldkit.aggregated_sensor(id),
	module_id integer NOT NULL REFERENCES fieldkit.station_module(id),
	nsamples integer NOT NULL,
	location geometry(POINT, 4326),
	value FLOAT NOT NULL
);

CREATE TABLE fieldkit.aggregated_bymod_1h (
	id SERIAL PRIMARY KEY,
	time TIMESTAMP NOT NULL,
	station_id integer NOT NULL REFERENCES fieldkit.station(id),
	sensor_id integer NOT NULL REFERENCES fieldkit.aggregated_sensor(id),
	module_id integer NOT NULL REFERENCES fieldkit.station_module(id),
	nsamples integer NOT NULL,
	location geometry(POINT, 4326),
	value FLOAT NOT NULL
);

CREATE TABLE fieldkit.aggregated_bymod_6h (
	id SERIAL PRIMARY KEY,
	time TIMESTAMP NOT NULL,
	station_id integer NOT NULL REFERENCES fieldkit.station(id),
	sensor_id integer NOT NULL REFERENCES fieldkit.aggregated_sensor(id),
	module_id integer NOT NULL REFERENCES fieldkit.station_module(id),
	nsamples integer NOT NULL,
	location geometry(POINT, 4326),
	value FLOAT NOT NULL
);

CREATE TABLE fieldkit.aggregated_bymod_12h (
	id SERIAL PRIMARY KEY,
	time TIMESTAMP NOT NULL,
	station_id integer NOT NULL REFERENCES fieldkit.station(id),
	sensor_id integer NOT NULL REFERENCES fieldkit.aggregated_sensor(id),
	module_id integer NOT NULL REFERENCES fieldkit.station_module(id),
	nsamples integer NOT NULL,
	location geometry(POINT, 4326),
	value FLOAT NOT NULL
);

CREATE TABLE fieldkit.aggregated_bymod_24h (
	id SERIAL PRIMARY KEY,
	time TIMESTAMP NOT NULL,
	station_id integer NOT NULL REFERENCES fieldkit.station(id),
	sensor_id integer NOT NULL REFERENCES fieldkit.aggregated_sensor(id),
	module_id integer NOT NULL REFERENCES fieldkit.station_module(id),
	nsamples integer NOT NULL,
	location geometry(POINT, 4326),
	value FLOAT NOT NULL
);

CREATE UNIQUE INDEX ON fieldkit.aggregated_bymod_10s (time, station_id, sensor_id, module_id);
CREATE UNIQUE INDEX ON fieldkit.aggregated_bymod_1m  (time, station_id, sensor_id, module_id);
CREATE UNIQUE INDEX ON fieldkit.aggregated_bymod_10m (time, station_id, sensor_id, module_id);
CREATE UNIQUE INDEX ON fieldkit.aggregated_bymod_30m (time, station_id, sensor_id, module_id);
CREATE UNIQUE INDEX ON fieldkit.aggregated_bymod_1h  (time, station_id, sensor_id, module_id);
CREATE UNIQUE INDEX ON fieldkit.aggregated_bymod_6h  (time, station_id, sensor_id, module_id);
CREATE UNIQUE INDEX ON fieldkit.aggregated_bymod_12h (time, station_id, sensor_id, module_id);
CREATE UNIQUE INDEX ON fieldkit.aggregated_bymod_24h (time, station_id, sensor_id, module_id);
