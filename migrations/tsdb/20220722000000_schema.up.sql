CREATE SCHEMA IF NOT EXISTS fieldkit;

CREATE EXTENSION IF NOT EXISTS timescaledb;

CREATE TABLE fieldkit.sensor_data (
    time TIMESTAMPTZ NOT NULL,
    station_id INTEGER NOT NULL /* REFERENCES fieldkit.station (id) */,
    module_id INTEGER NOT NULL /* REFERENCES fieldkit.station_module (id) */,
    sensor_id INTEGER NOT NULL /* REFERENCES fieldkit.aggregated_sensor (id) */,
    value FLOAT NOT NULL
);

SELECT create_hypertable('fieldkit.sensor_data', 'time');

CREATE INDEX sensor_data_idx ON fieldkit.sensor_data (station_id, module_id, sensor_id, time DESC);

/*
INSERT INTO fieldkit.sensor_data SELECT time, station_id, module_id, sensor_id, value FROM fieldkit.aggregated_10s;
*/
