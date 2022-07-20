ALTER TABLE fieldkit.twitter_oauth DROP COLUMN source_id;

DROP VIEW IF EXISTS fieldkit.device_locations_view;
DROP VIEW IF EXISTS fieldkit.record_visible;

DROP TABLE IF EXISTS fieldkit.record_analysis;
DROP TABLE IF EXISTS fieldkit.old_record;

DROP TABLE fieldkit.source_twitter_account;
DROP TABLE fieldkit.device_stream;
DROP TABLE fieldkit.device_firmware;
DROP TABLE fieldkit.device_schema;
DROP TABLE fieldkit.schema;
DROP TABLE fieldkit.device_notes;
DROP TABLE fieldkit.device_location;
DROP TABLE fieldkit.device;
DROP TABLE fieldkit.source;
DROP TABLE fieldkit.source_token;
DROP TABLE fieldkit.team_user;
DROP TABLE fieldkit.team;
DROP TABLE fieldkit.expedition;
