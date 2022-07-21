ALTER TABLE fieldkit.station ADD COLUMN location geometry(POINT, 4326);
ALTER TABLE fieldkit.station ADD COLUMN location_name TEXT;

ALTER TABLE fieldkit.station ADD COLUMN updated_at TIMESTAMP;
UPDATE fieldkit.station SET updated_at = created_at;
ALTER TABLE fieldkit.station ALTER COLUMN updated_at SET NOT NULL;
