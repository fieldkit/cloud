ALTER TABLE fieldkit.project ADD COLUMN bounds JSON;
ALTER TABLE fieldkit.project ADD COLUMN show_stations BOOLEAN NOT NULL DEFAULT false;

