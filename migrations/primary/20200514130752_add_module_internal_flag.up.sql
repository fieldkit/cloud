ALTER TABLE fieldkit.station_module ADD COLUMN flags INTEGER;

UPDATE fieldkit.station_module SET flags = 0 WHERE position != 255;

UPDATE fieldkit.station_module SET flags = 1 WHERE position  = 255;

ALTER TABLE fieldkit.station_module ALTER COLUMN flags SET NOT NULL;
