CREATE UNIQUE INDEX sensor_index_idx ON fieldkit.sensor_data (time, station_id, module_id, sensor_id);