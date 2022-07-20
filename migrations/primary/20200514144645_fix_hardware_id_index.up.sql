DROP INDEX fieldkit.station_module_hardware_id_idx;

CREATE UNIQUE INDEX ON fieldkit.station_module(meta_record_id, hardware_id);
