CREATE TABLE fieldkit.aggregated_sensor_updated_temporary (
    station_id INTEGER NOT NULL REFERENCES fieldkit.station (id),
    sensor_id INTEGER NOT NULL REFERENCES fieldkit.aggregated_sensor (id),
    module_id INTEGER NOT NULL REFERENCES fieldkit.station_module (id),
    time TIMESTAMP NOT NULL
);

CREATE UNIQUE INDEX aggregated_sensor_updated_idx ON fieldkit.aggregated_sensor_updated_temporary (station_id, sensor_id, module_id);

INSERT INTO fieldkit.aggregated_sensor_updated_temporary SELECT station_id, sensor_id, module_id, time FROM fieldkit.aggregated_sensor_updated;

DROP MATERIALIZED VIEW fieldkit.aggregated_sensor_updated;

ALTER TABLE fieldkit.aggregated_sensor_updated_temporary RENAME TO aggregated_sensor_updated;