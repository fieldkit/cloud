CREATE MATERIALIZED VIEW fieldkit.aggregated_sensor_updated AS
SELECT
	station_id,
	sensor_id,
	module_id,
	MAX(time) AS time
FROM fieldkit.aggregated_10s AS agg
GROUP by station_id, sensor_id, module_id;