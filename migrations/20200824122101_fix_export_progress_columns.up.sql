ALTER TABLE fieldkit.data_export ADD COLUMN message TEXT;
ALTER TABLE fieldkit.data_export RENAME COLUMN kind TO format;
