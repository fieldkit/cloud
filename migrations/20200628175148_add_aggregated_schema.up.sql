CREATE TABLE fieldkit.aggregated_sensor (
	id SERIAL PRIMARY KEY,
	key TEXT NOT NULL
);

CREATE TABLE fieldkit.aggregated_minutely (
	id SERIAL PRIMARY KEY,
	time TIMESTAMP NOT NULL,
	station_id integer NOT NULL REFERENCES fieldkit.station(id),
	sensor_id integer NOT NULL REFERENCES fieldkit.aggregated_sensor(id),
	location geometry(POINT, 4326),
	value FLOAT NOT NULL
);

CREATE TABLE fieldkit.aggregated_hourly (
	id SERIAL PRIMARY KEY,
	time TIMESTAMP NOT NULL,
	station_id integer NOT NULL REFERENCES fieldkit.station(id),
	sensor_id integer NOT NULL REFERENCES fieldkit.aggregated_sensor(id),
	location geometry(POINT, 4326),
	value FLOAT NOT NULL
);

CREATE TABLE fieldkit.aggregated_daily (
	id SERIAL PRIMARY KEY,
	time TIMESTAMP NOT NULL,
	station_id integer NOT NULL REFERENCES fieldkit.station(id),
	sensor_id integer NOT NULL REFERENCES fieldkit.aggregated_sensor(id),
	location geometry(POINT, 4326),
	value FLOAT NOT NULL
);

CREATE UNIQUE INDEX ON fieldkit.aggregated_sensor (key);

CREATE UNIQUE INDEX ON fieldkit.aggregated_minutely (time, station_id, sensor_id);
CREATE UNIQUE INDEX ON fieldkit.aggregated_hourly (time, station_id, sensor_id);
CREATE UNIQUE INDEX ON fieldkit.aggregated_daily (time, station_id, sensor_id);
