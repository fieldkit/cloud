ALTER TABLE fieldkit.module_sensor DROP COLUMN reading;
ALTER TABLE fieldkit.module_sensor ADD COLUMN reading_last FLOAT;
ALTER TABLE fieldkit.module_sensor ADD COLUMN reading_time TIMESTAMP;
