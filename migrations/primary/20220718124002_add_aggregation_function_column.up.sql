ALTER TABLE fieldkit.sensor_meta ADD COLUMN aggregation_function TEXT;
UPDATE fieldkit.sensor_meta SET aggregation_function = 'avg' WHERE id != 70;
UPDATE fieldkit.sensor_meta SET aggregation_function = 'max' WHERE id = 70;