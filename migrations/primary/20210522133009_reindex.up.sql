CREATE UNIQUE INDEX aggregated_10s_idx ON fieldkit.aggregated_10s USING btree (station_id ASC, module_id ASC, sensor_id ASC, time ASC);
ANALYZE;

CREATE UNIQUE INDEX aggregated_1m_idx ON fieldkit.aggregated_1m USING btree (station_id ASC, module_id ASC, sensor_id ASC, time ASC);
ANALYZE;

CREATE UNIQUE INDEX aggregated_10m_idx ON fieldkit.aggregated_10m USING btree (station_id ASC, module_id ASC, sensor_id ASC, time ASC);
ANALYZE;

CREATE UNIQUE INDEX aggregated_30m_idx ON fieldkit.aggregated_30m USING btree (station_id ASC, module_id ASC, sensor_id ASC, time ASC);
ANALYZE;

CREATE UNIQUE INDEX aggregated_1h_idx ON fieldkit.aggregated_1h USING btree (station_id ASC, module_id ASC, sensor_id ASC, time ASC);
ANALYZE;

CREATE UNIQUE INDEX aggregated_6h_idx ON fieldkit.aggregated_6h USING btree (station_id ASC, module_id ASC, sensor_id ASC, time ASC);
ANALYZE;

CREATE UNIQUE INDEX aggregated_12h_idx ON fieldkit.aggregated_12h USING btree (station_id ASC, module_id ASC, sensor_id ASC, time ASC);
ANALYZE;

CREATE UNIQUE INDEX aggregated_24h_idx ON fieldkit.aggregated_24h USING btree (station_id ASC, module_id ASC, sensor_id ASC, time ASC);
ANALYZE;


DROP INDEX IF EXISTS aggregated_bymod_10s_time_station_id_sensor_id_module_id_idx;
DROP INDEX IF EXISTS aggregated_bymod_1m_time_station_id_sensor_id_module_id_idx;
DROP INDEX IF EXISTS aggregated_bymod_10m_time_station_id_sensor_id_module_id_idx;
DROP INDEX IF EXISTS aggregated_bymod_30m_time_station_id_sensor_id_module_id_idx;
DROP INDEX IF EXISTS aggregated_bymod_1h_time_station_id_sensor_id_module_id_idx;
DROP INDEX IF EXISTS aggregated_bymod_6h_time_station_id_sensor_id_module_id_idx;
DROP INDEX IF EXISTS aggregated_bymod_12h_time_station_id_sensor_id_module_id_idx;
DROP INDEX IF EXISTS aggregated_bymod_24h_time_station_id_sensor_id_module_id_idx;
