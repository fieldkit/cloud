CREATE MATERIALIZED VIEW fieldkit.sensor_data_365d
WITH (timescaledb.continuous) AS
SELECT
    time_bucket('365 days', "time") AS bucket_time,
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

CREATE MATERIALIZED VIEW fieldkit.sensor_data_7d
WITH (timescaledb.continuous) AS
SELECT
    time_bucket('7 days', "time") AS bucket_time,
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


CREATE MATERIALIZED VIEW fieldkit.sensor_data_24h
WITH (timescaledb.continuous) AS
SELECT
    time_bucket('1 day', "time") AS bucket_time,
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


CREATE MATERIALIZED VIEW fieldkit.sensor_data_6h
WITH (timescaledb.continuous) AS
SELECT
    time_bucket('6 hours', "time") AS bucket_time,
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


CREATE MATERIALIZED VIEW fieldkit.sensor_data_1h
WITH (timescaledb.continuous) AS
SELECT
    time_bucket('1 hours', "time") AS bucket_time,
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


CREATE MATERIALIZED VIEW fieldkit.sensor_data_10m
WITH (timescaledb.continuous) AS
SELECT
    time_bucket('10 minute', "time") AS bucket_time,
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
