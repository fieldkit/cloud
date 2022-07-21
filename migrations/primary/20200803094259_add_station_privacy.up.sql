
ALTER TABLE fieldkit.project ADD COLUMN privacy INTEGER;
UPDATE fieldkit.project SET privacy = 0;
UPDATE fieldkit.project SET privacy = 1 WHERE NOT private;
ALTER TABLE fieldkit.project ALTER COLUMN privacy SET NOT NULL;

ALTER TABLE fieldkit.station DROP COLUMN private;
ALTER TABLE fieldkit.project DROP COLUMN private;
