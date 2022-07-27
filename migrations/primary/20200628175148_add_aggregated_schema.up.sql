CREATE TABLE fieldkit.aggregated_sensor (
	id SERIAL PRIMARY KEY,
	key TEXT NOT NULL
);

CREATE TABLE fieldkit.aggregated_1m (
	id SERIAL PRIMARY KEY,
	time TIMESTAMP NOT NULL,
	station_id integer NOT NULL REFERENCES fieldkit.station(id),
	sensor_id integer NOT NULL REFERENCES fieldkit.aggregated_sensor(id),
	location geometry(POINT, 4326),
	value FLOAT NOT NULL
);

CREATE TABLE fieldkit.aggregated_10m (
	id SERIAL PRIMARY KEY,
	time TIMESTAMP NOT NULL,
	station_id integer NOT NULL REFERENCES fieldkit.station(id),
	sensor_id integer NOT NULL REFERENCES fieldkit.aggregated_sensor(id),
	location geometry(POINT, 4326),
	value FLOAT NOT NULL
);

CREATE TABLE fieldkit.aggregated_30m (
	id SERIAL PRIMARY KEY,
	time TIMESTAMP NOT NULL,
	station_id integer NOT NULL REFERENCES fieldkit.station(id),
	sensor_id integer NOT NULL REFERENCES fieldkit.aggregated_sensor(id),
	location geometry(POINT, 4326),
	value FLOAT NOT NULL
);

CREATE TABLE fieldkit.aggregated_1h (
	id SERIAL PRIMARY KEY,
	time TIMESTAMP NOT NULL,
	station_id integer NOT NULL REFERENCES fieldkit.station(id),
	sensor_id integer NOT NULL REFERENCES fieldkit.aggregated_sensor(id),
	location geometry(POINT, 4326),
	value FLOAT NOT NULL
);

CREATE TABLE fieldkit.aggregated_6h (
	id SERIAL PRIMARY KEY,
	time TIMESTAMP NOT NULL,
	station_id integer NOT NULL REFERENCES fieldkit.station(id),
	sensor_id integer NOT NULL REFERENCES fieldkit.aggregated_sensor(id),
	location geometry(POINT, 4326),
	value FLOAT NOT NULL
);

CREATE TABLE fieldkit.aggregated_12h (
	id SERIAL PRIMARY KEY,
	time TIMESTAMP NOT NULL,
	station_id integer NOT NULL REFERENCES fieldkit.station(id),
	sensor_id integer NOT NULL REFERENCES fieldkit.aggregated_sensor(id),
	location geometry(POINT, 4326),
	value FLOAT NOT NULL
);

CREATE TABLE fieldkit.aggregated_24h (
	id SERIAL PRIMARY KEY,
	time TIMESTAMP NOT NULL,
	station_id integer NOT NULL REFERENCES fieldkit.station(id),
	sensor_id integer NOT NULL REFERENCES fieldkit.aggregated_sensor(id),
	location geometry(POINT, 4326),
	value FLOAT NOT NULL
);

CREATE UNIQUE INDEX ON fieldkit.aggregated_sensor (key);

CREATE UNIQUE INDEX ON fieldkit.aggregated_1m  (time, station_id, sensor_id);
CREATE UNIQUE INDEX ON fieldkit.aggregated_10m (time, station_id, sensor_id);
CREATE UNIQUE INDEX ON fieldkit.aggregated_30m (time, station_id, sensor_id);
CREATE UNIQUE INDEX ON fieldkit.aggregated_1h  (time, station_id, sensor_id);
CREATE UNIQUE INDEX ON fieldkit.aggregated_6h  (time, station_id, sensor_id);
CREATE UNIQUE INDEX ON fieldkit.aggregated_12h (time, station_id, sensor_id);
CREATE UNIQUE INDEX ON fieldkit.aggregated_24h (time, station_id, sensor_id);
