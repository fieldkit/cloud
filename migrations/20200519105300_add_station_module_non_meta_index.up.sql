ALTER TABLE fieldkit.station_module ADD COLUMN source_id INTEGER;
ALTER TABLE fieldkit.station_module ADD COLUMN updated_at TIMESTAMP;
UPDATE fieldkit.station_module SET updated_at = '2020-01-01 00:00:00';
ALTER TABLE fieldkit.station_module ALTER COLUMN updated_at SET NOT NULL;

CREATE        INDEX ON fieldkit.station_module(updated_at);
CREATE UNIQUE INDEX ON fieldkit.station_module(hardware_id, source_id);

/* Give them similar names due to the order :) */
DROP INDEX fieldkit.station_module_meta_record_id_hardware_id_idx;
CREATE UNIQUE INDEX ON fieldkit.station_module(hardware_id, meta_record_id);
