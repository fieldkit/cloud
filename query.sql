/*
SELECT id, name FROM fieldkit.station;
*/

TRUNCATE TABLE fieldkit.data_visibility;

INSERT INTO fieldkit.data_visibility (start_time, end_time, station_id, user_id) VALUES ('2018/01/01', '2021/01/01', 2, 2);
INSERT INTO fieldkit.data_visibility (start_time, end_time, station_id, user_id) VALUES ('2020/01/01', '2020/03/01', 2, 3);
INSERT INTO fieldkit.data_visibility (start_time, end_time, station_id, user_id) VALUES ('2020/04/01', '2020/08/01', 2, 3);

SELECT * FROM fieldkit.data_visibility ORDER BY start_time;


		WITH
		with_timestamp_differences AS (
			SELECT
				*,
										   LAG(time) OVER (ORDER BY time) AS previous_timestamp,
				EXTRACT(epoch FROM (time - LAG(time) OVER (ORDER BY time))) AS time_difference
			FROM fieldkit.aggregated_24h
			WHERE time >= '2019/01/01' AND time <= '2021/01/01' AND station_id IN (2) AND sensor_id IN (1)
			ORDER BY time
		),
		with_temporal_clustering AS (
			SELECT
				*,
				CASE WHEN s.time_difference > 86400
					OR s.time_difference IS NULL THEN true
					ELSE NULL
				END AS new_temporal_cluster
			FROM with_timestamp_differences AS s
		),
		time_grouped AS (
			SELECT
				*,
				COUNT(new_temporal_cluster) OVER (
					ORDER BY s.time
					ROWS UNBOUNDED PRECEDING
				) AS time_group
			FROM with_temporal_clustering s
		)
		SELECT
			id,
			time,
			station_id,
			sensor_id,
			ST_AsBinary(location) AS location,
			value,
			time_group
		FROM time_grouped;
