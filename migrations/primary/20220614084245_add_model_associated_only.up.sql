ALTER TABLE fieldkit.station_model ADD COLUMN only_visible_via_association BOOLEAN;
UPDATE fieldkit.station_model SET only_visible_via_association = FALSE;
ALTER TABLE fieldkit.station_model ALTER only_visible_via_association SET NOT NULL;