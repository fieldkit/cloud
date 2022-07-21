/**
 * Existing data w/o this column isn't that useful so may as well purge
 * and then reprocess to get good data into the tables.
 */
DELETE FROM fieldkit.module_sensor;
DELETE FROM fieldkit.station_module;

ALTER TABLE fieldkit.station_module ADD COLUMN module_index INTEGER NOT NULL;
ALTER TABLE fieldkit.module_sensor RENAME COLUMN position TO sensor_index;
