CREATE OR REPLACE VIEW fieldkit.project_and_station_activity AS
    SELECT sa.created_at, ps.project_id, sa.station_id, sa.id AS station_activity_id, NULL AS project_activity_id
	  FROM fieldkit.project_station AS ps
      JOIN fieldkit.station_activity AS sa ON (ps.station_id = sa.station_id)
    UNION
	SELECT pa.created_at, pa.project_id, NULL, NULL, pa.id AS project_activity_id
	  FROM fieldkit.project_activity AS pa;
