ALTER TABLE fieldkit.station ADD COLUMN status TEXT;

UPDATE fieldkit.station SET status = 'up';