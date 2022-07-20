CREATE TABLE fieldkit.station_interestingness (
    id SERIAL PRIMARY KEY,
    station_id INTEGER NOT NULL REFERENCES fieldkit.station (id),
    window_seconds INTEGER NOT NULL,
    interestingness FLOAT NOT NULL,
    reading_sensor_id INTEGER NOT NULL REFERENCES fieldkit.aggregated_sensor (id),
    reading_module_id INTEGER NOT NULL REFERENCES fieldkit.station_module (id),
    reading_value FLOAT NOT NULL,
    reading_time TIMESTAMP NOT NULL
);

CREATE UNIQUE INDEX station_interestingness_idx ON fieldkit.station_interestingness (station_id, window_seconds);

ALTER TABLE fieldkit.aggregated_sensor ADD COLUMN interestingness_priority INTEGER;

UPDATE fieldkit.aggregated_sensor SET interestingness_priority = 0 WHERE id IN (44);