CREATE TABLE fieldkit.aggregated_10s (
	id SERIAL PRIMARY KEY,
	time TIMESTAMP NOT NULL,
	station_id integer NOT NULL REFERENCES fieldkit.station(id),
	sensor_id integer NOT NULL REFERENCES fieldkit.aggregated_sensor(id),
	nsamples integer NOT NULL,
	location geometry(POINT, 4326),
	value FLOAT NOT NULL
);

CREATE UNIQUE INDEX ON fieldkit.aggregated_10s  (time, station_id, sensor_id);
