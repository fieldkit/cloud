ALTER TABLE fieldkit.station ADD COLUMN hidden BOOL;
UPDATE fieldkit.station SET hidden = false;