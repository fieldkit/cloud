CREATE MATERIALIZED VIEW fieldkit.sensor_data_1m
WITH (timescaledb.continuous) AS
SELECT
    time_bucket('1 minute', "time") AS bucket_time,
    station_id,
    module_id,
    sensor_id,
    COUNT(*) AS bucket_samples,
    MIN(time) AS data_start,
    MAX(time) AS data_end,
    AVG(value) AS avg_value,
    MIN(value) AS min_value,
    MAX(value) AS max_value,
    LAST(value, time) AS last_value
FROM fieldkit.sensor_data
GROUP BY bucket_time, station_id, module_id, sensor_id
WITH NO DATA;
