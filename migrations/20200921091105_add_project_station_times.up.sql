ALTER TABLE fieldkit.project_station ADD COLUMN created_at TIMESTAMP;
UPDATE fieldkit.project_station SET created_at = '2000-01-01';
ALTER TABLE fieldkit.project_station ALTER COLUMN created_at SET NOT NULL;

ALTER TABLE fieldkit.project_user ADD COLUMN created_at TIMESTAMP;
UPDATE fieldkit.project_user SET created_at = '2000-01-01';
ALTER TABLE fieldkit.project_user ALTER COLUMN created_at SET NOT NULL;
