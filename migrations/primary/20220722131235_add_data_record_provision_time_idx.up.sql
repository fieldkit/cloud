CREATE INDEX data_record_provision_id_idx ON fieldkit.data_record (provision_id);
CREATE INDEX data_record_provision_id_time_idx ON fieldkit.data_record (provision_id, time)