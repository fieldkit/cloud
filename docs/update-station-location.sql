

SELECT
	id, name, status_json::json->'latitude', status_json::json->'longitude'
FROM fieldkit.station WHERE owner_id = 5;

UPDATE
fieldkit.station
SET status_json = (status_json::jsonb || jsonb_build_object('latitude', 34.064542) || jsonb_build_object('longitude', -118.293137))
WHERE id = 172;


UPDATE
fieldkit.station
SET status_json = (status_json::jsonb || jsonb_build_object('latitude', 34.004468) || jsonb_build_object('longitude', -117.361303))
WHERE id = 165;


UPDATE
fieldkit.station
SET status_json = (status_json::jsonb || jsonb_build_object('latitude', 35.493026) || jsonb_build_object('longitude', -114.683632))
WHERE id = 178;
