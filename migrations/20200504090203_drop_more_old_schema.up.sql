DROP TABLE fieldkit.device_stream_location;

ALTER SEQUENCE fieldkit.ingestion_id_seq RENAME TO ingestion_old_id_seq;

ALTER SEQUENCE fieldkit.ingestion_id_seq1 RENAME TO ingestion_id_seq;
